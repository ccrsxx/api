package auth

import (
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/google/uuid"
)

type Controller struct {
	service *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{
		service: svc,
	}
}

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Image    string `json:"image"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (c *Controller) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r.Context())

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	err = api.NewSuccessResponse(w, http.StatusOK, UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Role:     user.Role,
		Email:    user.Email.String,
		Image:    user.Image.String,
		Username: user.Username.String,
	})

	if err != nil {
		slog.Warn("get me response error", "error", err)
		return
	}
}

func (c *Controller) LoginGithub(w http.ResponseWriter, r *http.Request) {
	_, err := c.service.ValidateOauthToken(r.Context(), r)

	// If user is already logged in, redirect to back guestbook
	if err == nil {
		frontendGuestbookURL := c.service.frontendPublicURL + "/guestbook"

		http.Redirect(w, r, frontendGuestbookURL, http.StatusTemporaryRedirect)

		return
	}

	randomID := uuid.New().String()

	url := c.service.githubOauthConfig.AuthCodeURL(randomID)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth-state",
		Path:     "/",
		Value:    randomID,
		MaxAge:   int(oauthStateExpiry.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Domain is intentionally omitted so it stays securely locked to your API domain
	})

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (c *Controller) LogoutGithub(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth-token",
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		Domain:   c.service.frontendBaseURL,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth-state",
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Domain is intentionally omitted so it stays securely locked to your API domain
	})

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) LoginGithubCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	oauthToken, err := c.service.ValidateOauthState(ctx, r)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	sessionToken, err := c.service.CreateOauthTokenForGithubUser(ctx, oauthToken.AccessToken)

	if err != nil {
		api.HandleHTTPError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth-token",
		Path:     "/",
		Value:    sessionToken,
		MaxAge:   int(oauthTokenExpiry.Seconds()),
		Secure:   true,
		HttpOnly: true,
		Domain:   c.service.frontendBaseURL,
		SameSite: http.SameSiteLaxMode,
	})

	frontendGuestbookURL := c.service.frontendPublicURL + "/guestbook"

	http.Redirect(w, r, frontendGuestbookURL, http.StatusTemporaryRedirect)
}
