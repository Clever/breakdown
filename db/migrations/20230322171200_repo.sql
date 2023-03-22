-- +goose Up
CREATE TABLE IF NOT EXISTS repo (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS repo_commit (
    id BIGSERIAL PRIMARY KEY,
    repo_id BIGINT NOT NULL,
    commit_sha CHAR(8) NOT NULL,
    commit_date TIMESTAMP WITH TIME ZONE NOT NULL,
    meta JSONB,
    UNIQUE(commit_sha),
    FOREIGN KEY(repo_id) REFERENCES repo(id)
);

CREATE INDEX repo_commit__repo_id_commit_date ON repo_commit (repo_id, commit_date);

CREATE TYPE ci_source AS ENUM('circle-ci', 'github-actions');

CREATE TABLE IF NOT EXISTS ci_workflow (
    id BIGSERIAL PRIMARY KEY,
    source ci_source NOT NULL,
    source_id TEXT NOT NULL,
    repo_id BIGINT  NOT NULL,
    commit_sha CHAR(8) NOT NULL,
    workflow_time TIMESTAMP WITH TIME ZONE NOT NULL,
    duration_s INT NOT NULL,
    FOREIGN KEY(repo_id) REFERENCES repo(id)
);

CREATE INDEX ci_workflow__repo_id__workflow_date ON ci_workflow(repo_id, workflow_time);

CREATE TYPE package_type AS ENUM('gomod', 'npm');

CREATE TABLE IF NOT EXISTS package_file (
    id BIGSERIAL PRIMARY KEY,
    repo_commit_id BIGINT  NOT NULL,
    path TEXT NOT NULL,
    type package_type NOT NULL,
    meta JSONB,
    UNIQUE(repo_commit_id, path),
    FOREIGN KEY (repo_commit_id) REFERENCES repo_commit(id)
);

CREATE TABLE IF NOT EXISTS dependency (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    type package_type NOT NULL,
    is_local boolean NOT NULL,
    UNIQUE(type, name, version)
);

CREATE TABLE IF NOT EXISTS dep_dependency (
    parent_id BIGINT  NOT NULL,
    dependency_id BIGINT  NOT NULL,
    FOREIGN KEY(parent_id) REFERENCES dependency(id),
    FOREIGN KEY(dependency_id) REFERENCES dependency(id)
);

CREATE TABLE IF NOT EXISTS package_file_dependency (
    package_file_id BIGINT  NOT NULL,
    dependency_id BIGINT  NOT NULL,
    UNIQUE(package_file_id, dependency_id),
    FOREIGN KEY(package_file_id) REFERENCES package_file(id),
    FOREIGN KEY(dependency_id) REFERENCES dependency(id)
);

CREATE TABLE IF NOT EXISTS deployment (
    id BIGSERIAL PRIMARY KEY,
    repo_commit_id BIGINT NOT NULL,
    application TEXT,
    version TEXT,
    date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    run_type TEXT,
    FOREIGN KEY(repo_commit_id) REFERENCES repo_commit(id)
);

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
