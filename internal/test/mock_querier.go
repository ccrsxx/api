package test

import (
	"context"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// MockQuerier implements sqlc.Querier with overridable function fields.
// Any method not explicitly set will panic, catching unmocked calls instantly.
var _ sqlc.Querier = (*MockQuerier)(nil)

type MockQuerier struct {
	CreateAccountFn          func(ctx context.Context, arg sqlc.CreateAccountParams) (sqlc.Account, error)
	CreateGuestbookFn        func(ctx context.Context, arg sqlc.CreateGuestbookParams) (sqlc.CreateGuestbookRow, error)
	CreateUserFn             func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	DeleteGuestbookFn        func(ctx context.Context, id pgtype.UUID) error
	GetAccountByProviderFn   func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error)
	GetContentBySlugFn       func(ctx context.Context, slug string) (sqlc.Content, error)
	GetContentLikeStatusFn   func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error)
	GetContentStatsByTypeFn  func(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error)
	GetGuestbookByIDFn       func(ctx context.Context, id pgtype.UUID) (sqlc.Guestbook, error)
	GetTotalContentMetaFn    func(ctx context.Context, contentID pgtype.UUID) (sqlc.GetTotalContentMetaRow, error)
	GetUserByIDFn            func(ctx context.Context, id pgtype.UUID) (sqlc.User, error)
	GetUserWithAccountByIDFn func(ctx context.Context, id pgtype.UUID) (sqlc.GetUserWithAccountByIDRow, error)
	IncrementContentLikeFn   func(ctx context.Context, arg sqlc.IncrementContentLikeParams) (sqlc.IncrementContentLikeRow, error)
	IncrementContentViewFn   func(ctx context.Context, arg sqlc.IncrementContentViewParams) (sqlc.IncrementContentViewRow, error)
	ListContentByTypeFn      func(ctx context.Context, type_ string) ([]sqlc.ListContentByTypeRow, error)
	ListGuestbookFn          func(ctx context.Context) ([]sqlc.ListGuestbookRow, error)
	UpdateUserFn             func(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	UpsertContentFn          func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error)
	UpsertIPAddressFn        func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error)
}

func (m *MockQuerier) CreateAccount(ctx context.Context, arg sqlc.CreateAccountParams) (sqlc.Account, error) {
	if m.CreateAccountFn == nil {
		panic("MockQuerier.CreateAccount called but not mocked")
	}
	return m.CreateAccountFn(ctx, arg)
}

func (m *MockQuerier) CreateGuestbook(ctx context.Context, arg sqlc.CreateGuestbookParams) (sqlc.CreateGuestbookRow, error) {
	if m.CreateGuestbookFn == nil {
		panic("MockQuerier.CreateGuestbook called but not mocked")
	}
	return m.CreateGuestbookFn(ctx, arg)
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	if m.CreateUserFn == nil {
		panic("MockQuerier.CreateUser called but not mocked")
	}
	return m.CreateUserFn(ctx, arg)
}

func (m *MockQuerier) DeleteGuestbook(ctx context.Context, id pgtype.UUID) error {
	if m.DeleteGuestbookFn == nil {
		panic("MockQuerier.DeleteGuestbook called but not mocked")
	}
	return m.DeleteGuestbookFn(ctx, id)
}

func (m *MockQuerier) GetAccountByProvider(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
	if m.GetAccountByProviderFn == nil {
		panic("MockQuerier.GetAccountByProvider called but not mocked")
	}
	return m.GetAccountByProviderFn(ctx, arg)
}

func (m *MockQuerier) GetContentBySlug(ctx context.Context, slug string) (sqlc.Content, error) {
	if m.GetContentBySlugFn == nil {
		panic("MockQuerier.GetContentBySlug called but not mocked")
	}
	return m.GetContentBySlugFn(ctx, slug)
}

func (m *MockQuerier) GetContentLikeStatus(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
	if m.GetContentLikeStatusFn == nil {
		panic("MockQuerier.GetContentLikeStatus called but not mocked")
	}
	return m.GetContentLikeStatusFn(ctx, arg)
}

func (m *MockQuerier) GetContentStatsByType(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error) {
	if m.GetContentStatsByTypeFn == nil {
		panic("MockQuerier.GetContentStatsByType called but not mocked")
	}
	return m.GetContentStatsByTypeFn(ctx, type_)
}

func (m *MockQuerier) GetGuestbookByID(ctx context.Context, id pgtype.UUID) (sqlc.Guestbook, error) {
	if m.GetGuestbookByIDFn == nil {
		panic("MockQuerier.GetGuestbookByID called but not mocked")
	}
	return m.GetGuestbookByIDFn(ctx, id)
}

func (m *MockQuerier) GetTotalContentMeta(ctx context.Context, contentID pgtype.UUID) (sqlc.GetTotalContentMetaRow, error) {
	if m.GetTotalContentMetaFn == nil {
		panic("MockQuerier.GetTotalContentMeta called but not mocked")
	}
	return m.GetTotalContentMetaFn(ctx, contentID)
}

func (m *MockQuerier) GetUserByID(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
	if m.GetUserByIDFn == nil {
		panic("MockQuerier.GetUserByID called but not mocked")
	}
	return m.GetUserByIDFn(ctx, id)
}

func (m *MockQuerier) GetUserWithAccountByID(ctx context.Context, id pgtype.UUID) (sqlc.GetUserWithAccountByIDRow, error) {
	if m.GetUserWithAccountByIDFn == nil {
		panic("MockQuerier.GetUserWithAccountByID called but not mocked")
	}
	return m.GetUserWithAccountByIDFn(ctx, id)
}

func (m *MockQuerier) IncrementContentLike(ctx context.Context, arg sqlc.IncrementContentLikeParams) (sqlc.IncrementContentLikeRow, error) {
	if m.IncrementContentLikeFn == nil {
		panic("MockQuerier.IncrementContentLike called but not mocked")
	}
	return m.IncrementContentLikeFn(ctx, arg)
}

func (m *MockQuerier) IncrementContentView(ctx context.Context, arg sqlc.IncrementContentViewParams) (sqlc.IncrementContentViewRow, error) {
	if m.IncrementContentViewFn == nil {
		panic("MockQuerier.IncrementContentView called but not mocked")
	}
	return m.IncrementContentViewFn(ctx, arg)
}

func (m *MockQuerier) ListContentByType(ctx context.Context, type_ string) ([]sqlc.ListContentByTypeRow, error) {
	if m.ListContentByTypeFn == nil {
		panic("MockQuerier.ListContentByType called but not mocked")
	}
	return m.ListContentByTypeFn(ctx, type_)
}

func (m *MockQuerier) ListGuestbook(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
	if m.ListGuestbookFn == nil {
		panic("MockQuerier.ListGuestbook called but not mocked")
	}
	return m.ListGuestbookFn(ctx)
}

func (m *MockQuerier) UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	if m.UpdateUserFn == nil {
		panic("MockQuerier.UpdateUser called but not mocked")
	}
	return m.UpdateUserFn(ctx, arg)
}

func (m *MockQuerier) UpsertContent(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
	if m.UpsertContentFn == nil {
		panic("MockQuerier.UpsertContent called but not mocked")
	}
	return m.UpsertContentFn(ctx, arg)
}

func (m *MockQuerier) UpsertIPAddress(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
	if m.UpsertIPAddressFn == nil {
		panic("MockQuerier.UpsertIPAddress called but not mocked")
	}
	return m.UpsertIPAddressFn(ctx, ipAddress)
}
