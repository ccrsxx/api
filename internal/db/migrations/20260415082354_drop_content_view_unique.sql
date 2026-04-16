-- +goose Up
ALTER TABLE content_view DROP CONSTRAINT content_view_content_id_ip_address_id_key;

-- +goose Down
ALTER TABLE content_view ADD CONSTRAINT content_view_content_id_ip_address_id_key UNIQUE (content_id, ip_address_id);
