-- +goose Up
-- +goose StatementBegin
ALTER TABLE deployment ADD COLUMN environment TEXT NOT NULL;
ALTER TABLE deployment ALTER COLUMN application SET NOT NULL;
ALTER TABLE deployment ALTER COLUMN version SET NOT NULL;
ALTER TABLE deployment ALTER COLUMN run_type SET NOT NULL;


ALTER TABLE ci_workflow ADD COLUMN repo TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
