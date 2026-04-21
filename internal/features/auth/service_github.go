package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/clients/github"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
)

func (s *Service) CreateOauthTokenForGithubUser(ctx context.Context, accessToken string) (string, error) {
	githubUser, err := s.getGithubUser(ctx, accessToken)

	if err != nil {
		return "", fmt.Errorf("create session get github user error: %w", err)
	}

	user, err := s.createGithubUser(ctx, githubUser)

	if err != nil {
		return "", fmt.Errorf("create session upsert github user error: %w", err)
	}

	sessionToken, err := s.GenerateOauthToken(user.ID.String())

	if err != nil {
		return "", fmt.Errorf("create session generate session token error: %w", err)
	}

	return sessionToken, nil
}

func (s *Service) ValidateOauthState(ctx context.Context, r *http.Request) (oauth2.Token, error) {
	githubCookie, err := r.Cookie("oauth-state")

	if err != nil {
		return oauth2.Token{}, fmt.Errorf("validate oauth request cookie error: %w", err)
	}

	stateFromGithub := r.URL.Query().Get("state")

	if githubCookie.Value != stateFromGithub {
		return oauth2.Token{}, fmt.Errorf("validate oauth request invalid state")
	}

	oauthCode := r.URL.Query().Get("code")

	token, err := s.githubOauthConfig.Exchange(ctx, oauthCode)

	if err != nil {
		return oauth2.Token{}, fmt.Errorf("validate oauth request exchange error: %w", err)
	}

	if !token.Valid() {
		return oauth2.Token{}, fmt.Errorf("validate oauth request invalid token")
	}

	return *token, nil
}

func (s *Service) createGithubUser(ctx context.Context, githubUser github.User) (sqlc.User, error) {
	providerAccountID := strconv.FormatInt(githubUser.ID, 10)

	account, err := s.db.GetAccountByProvider(ctx, sqlc.GetAccountByProviderParams{
		Provider:          "github",
		ProviderAccountID: providerAccountID,
	})

	if err == nil {
		return s.updateExistingGithubUser(ctx, account.UserID, githubUser)
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return sqlc.User{}, fmt.Errorf("create github user get account error: %w", err)
	}

	return s.createNewGithubUser(ctx, githubUser, providerAccountID)
}

func (s *Service) createNewGithubUser(ctx context.Context, githubUser github.User, providerAccountID string) (sqlc.User, error) {
	githubName := githubUser.Login

	if githubUser.Name != nil {
		githubName = *githubUser.Name
	}

	githubEmail := ""

	if githubUser.Email != nil {
		githubEmail = *githubUser.Email
	}

	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return sqlc.User{}, fmt.Errorf("create new github user begin tx error: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Warn("create new github user rollback tx error", "error", err)
		}
	}()

	qtx := s.db.WithTx(tx)

	user, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  githubName,
		Image: pgtype.Text{String: githubUser.AvatarURL, Valid: githubUser.AvatarURL != ""},
		Email: pgtype.Text{String: githubEmail, Valid: githubEmail != ""},
	})

	if err != nil {
		return sqlc.User{}, fmt.Errorf("create new github user insert error: %w", err)
	}

	_, err = qtx.CreateAccount(ctx, sqlc.CreateAccountParams{
		UserID:            user.ID,
		Username:          pgtype.Text{String: githubUser.Login, Valid: true},
		Provider:          "github",
		ProviderAccountID: providerAccountID,
	})

	if err != nil {
		return sqlc.User{}, fmt.Errorf("create new github account insert error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.User{}, fmt.Errorf("create new github user commit tx error: %w", err)
	}

	return user, nil
}

func (s *Service) updateExistingGithubUser(ctx context.Context, userID pgtype.UUID, githubUser github.User) (sqlc.User, error) {
	user, err := s.db.GetUserByID(ctx, userID)

	if err != nil {
		return sqlc.User{}, fmt.Errorf("update github user get user error: %w", err)
	}

	parsedGithubName := user.Name

	if githubUser.Name != nil {
		parsedGithubName = *githubUser.Name
	}

	nameNeedsUpdate := user.Name != parsedGithubName
	avatarNeedsUpdate := user.Image.String != githubUser.AvatarURL

	if nameNeedsUpdate || avatarNeedsUpdate {
		updatedUser, err := s.db.UpdateUser(ctx, sqlc.UpdateUserParams{
			ID:    user.ID,
			Name:  parsedGithubName,
			Image: pgtype.Text{String: githubUser.AvatarURL, Valid: githubUser.AvatarURL != ""},
		})

		if err != nil {
			return sqlc.User{}, fmt.Errorf("update github user update db error: %w", err)
		}

		return updatedUser, nil
	}

	return user, nil
}
