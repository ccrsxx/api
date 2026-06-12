-- +goose Up
CREATE INDEX index_account_on_user_id ON account(user_id);

CREATE INDEX index_guestbook_on_created_at ON guestbook(created_at DESC);

-- +goose Down
DROP INDEX index_account_on_user_id;

DROP INDEX index_guestbook_on_created_at;