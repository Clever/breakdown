-- name: HealthCheck :exec
SELECT 1;

-- name: ListRepos :many
SELECT * FROM repo
ORDER BY name;

-- name: CreateRepo :one
WITH ins AS (
    INSERT INTO repo (
        name
    ) VALUES (
        $1
    )
    ON CONFLICT DO NOTHING
    RETURNING id
)
SELECT id
FROM repo r
WHERE r.name = $1
UNION
SELECT id FROM ins;

-- name: GetRepo :one
SELECT *
FROM repo
WHERE name = $1;

-- name: CreateRepoCommit :one
INSERT INTO repo_commit (
    repo_id, commit_sha, commit_date
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: CreatePackageFile :one
INSERT INTO package_file (
    repo_commit_id, path, type, meta
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: CreateDependency :batchmany
WITH ins AS (
    INSERT INTO dependency (
        name, version, type, is_local
    ) VALUES (
        $1, $2, $3, $4
    )
    ON CONFLICT DO NOTHING
    RETURNING *
)
SELECT * FROM dependency d
WHERE d.name = $1
    AND d.version = $2
    AND d.type = $3
UNION
SELECT * FROM ins;

-- name: GetDependencyId :one
SELECT id
FROM dependency
WHERE type = ?
    AND name = ?
    AND version = ?;

-- name: InsertPackageFileDependency :batchexec
INSERT INTO package_file_dependency (
    package_file_id, dependency_id
) VALUES (
    $1, $2
)
ON CONFLICT DO NOTHING;

-- name: InsertDepDependency :batchexec
INSERT INTO dep_dependency (
    parent_id, dependency_id
) VALUES (
    $1, $2
)
ON CONFLICT DO NOTHING;

-- name: InsertDeployment :batchexec
INSERT INTO deployment (
    commit_sha, application, environment, version, run_type
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetDeploys :many
SELECT * FROM deployment ORDER BY commit_sha;

-- name: GetCommit :one
SELECT *
FROM repo_commit
WHERE repo_id = (SELECT id FROM repo WHERE name = $1)
    AND commit_sha = $2
LIMIT 1;

-- name: GetCommits :many
SELECT r.name, rc.commit_sha, rc.meta
FROM repo_commit rc
LEFT JOIN repo r ON r.id = rc.repo_id
ORDER BY 1, 2;
