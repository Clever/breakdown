-- +goose Up
-- +goose StatementBegin
ALTER TABLE deployment DROP CONSTRAINT deployment_repo_commit_id_fkey;
ALTER TABLE deployment DROP COLUMN repo_commit_id;
ALTER TABLE deployment ADD COLUMN commit_sha CHAR(8) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
