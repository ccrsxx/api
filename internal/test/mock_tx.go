package test

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// MockTx implements pgx.Tx with overridable Commit and Rollback.
// All other methods panic to catch unmocked calls instantly.
var _ pgx.Tx = (*MockTx)(nil)

type MockTx struct {
	CommitFn   func(ctx context.Context) error
	RollbackFn func(ctx context.Context) error
}

func (m *MockTx) Commit(ctx context.Context) error {
	if m.CommitFn == nil {
		panic("MockTx.Commit called but not mocked")
	}
	return m.CommitFn(ctx)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	if m.RollbackFn == nil {
		panic("MockTx.Rollback called but not mocked")
	}
	return m.RollbackFn(ctx)
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	panic("MockTx.Begin called but not mocked")
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	panic("MockTx.CopyFrom called but not mocked")
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	panic("MockTx.SendBatch called but not mocked")
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	panic("MockTx.LargeObjects called but not mocked")
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	panic("MockTx.Prepare called but not mocked")
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	panic("MockTx.Exec called but not mocked")
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	panic("MockTx.Query called but not mocked")
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	panic("MockTx.QueryRow called but not mocked")
}

func (m *MockTx) Conn() *pgx.Conn {
	panic("MockTx.Conn called but not mocked")
}

// MockBeginner satisfies any interface with Begin(ctx) (pgx.Tx, error).
type MockBeginner struct {
	BeginFn func(ctx context.Context) (pgx.Tx, error)
}

func (m *MockBeginner) Begin(ctx context.Context) (pgx.Tx, error) {
	if m.BeginFn == nil {
		panic("MockBeginner.Begin called but not mocked")
	}
	return m.BeginFn(ctx)
}
