package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/Clever/breakdown/db"
	"github.com/Clever/breakdown/gen-go/models"
	"github.com/Clever/kayvee-go/v7/logger"
	"github.com/go-openapi/swag"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func getQueries() (*db.Queries, *pgx.Conn, error) {
	testDb, err := db.TestDB()
	if err != nil {
		return nil, nil, err
	}

	return db.New(testDb), testDb, nil
}

type controllerTest struct {
	name        string
	input       func(MyController) error
	expectError bool
	expected    func() []string
}

func TestDeploy(t *testing.T) {
	tests := []controllerTest{
		{
			name: "inserts one deploy successfully",
			input: func(mc MyController) error {
				return mc.PostDeploy(context.Background(), &models.Deploys{
					&models.Deploy{
						Application: swag.String("breakdown"),
						CommitSha:   swag.String("12345678"),
						Environment: swag.String("production"),
						RunType:     swag.String("docker"),
						Version:     swag.String("1"),
					},
				})
			},
			expectError: false,
			expected: func() []string {
				return []string{"breakdown 12345678 production docker 1"}
			},
		},
		{
			name: "errors on short commit sha",
			input: func(mc MyController) error {
				return mc.PostDeploy(context.Background(), &models.Deploys{
					&models.Deploy{
						Application: swag.String("breakdown"),
						CommitSha:   swag.String("1234"),
						Environment: swag.String("production"),
						RunType:     swag.String("docker"),
						Version:     swag.String("1"),
					},
				})
			},
			expectError: true,
			expected: func() []string {
				return []string{}
			},
		},
	}

	qs, db, err := getQueries()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	defer db.Close(ctx)

	testMC := MyController{
		db:      db,
		queries: qs,
		l:       logger.NewMockCountLogger("test"),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input(testMC)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("got error: %s", err)
			}

			deps, _ := qs.GetDeploys(ctx)
			deploys := []string{}
			for _, dep := range deps {
				deploys = append(deploys, fmt.Sprintf("%s %s %s %s %s", dep.Application, dep.CommitSha, dep.Environment, dep.RunType, dep.Version))
			}

			expectedRes := tt.expected()
			if len(expectedRes) != len(deploys) {
				t.Errorf("differing lengths of expected results, got=%d, want=%d", len(deploys), len(expectedRes))
			}

			for i, deploy := range deploys {
				if deploy != expectedRes[i] {
					t.Errorf("differing deploy at %d, got=%q, want=%q", i, deploy, expectedRes[i])
				}
			}

		})
	}
}

type testCommitInformation struct {
	name        string
	input       func(*pgx.Conn, MyController) error
	expectError bool
}

func TestGetCommitInformation(t *testing.T) {
	tests := []testCommitInformation{
		{
			name: "input commit sha not long enough",
			input: func(conn *pgx.Conn, mc MyController) error {
				_, err := mc.GetCommit(context.Background(), &models.GetCommitInformation{
					CommitSha: swag.String("123"),
					RepoName:  swag.String("breakdown"),
				})
				return err
			},
			expectError: true,
		},
		{
			name: "commit sha not found",
			input: func(conn *pgx.Conn, mc MyController) error {
				_, err := mc.GetCommit(context.Background(), &models.GetCommitInformation{
					CommitSha: swag.String("12312121212"),
					RepoName:  swag.String("breakdown"),
				})
				return err
			},
			expectError: true,
		},
		{
			name:        "finds commit_sha",
			expectError: false,
			input: func(c *pgx.Conn, mc MyController) error {
				var meta pgtype.JSONB
				meta.Set(`{"test":"some_test"}`)
				_, err := c.Exec(context.Background(), `
					WITH new_repo AS (
						INSERT INTO repo (name) VALUES ($1)
						RETURNING id
					)
					INSERT INTO repo_commit (repo_id, commit_sha, commit_date, meta)
					SELECT
						nr.id, $2, CURRENT_TIMESTAMP, $3
					FROM new_repo nr
				`, "breakdown", "11111111", meta)
				if err != nil {
					return fmt.Errorf("insert error: %s", err)
				}

				commit, err := mc.GetCommit(context.Background(), &models.GetCommitInformation{
					CommitSha: swag.String("11111111"),
					RepoName:  swag.String("breakdown"),
				})

				fmt.Printf("%+v", commit.Meta)

				if err != nil {
					return fmt.Errorf("error getting commit: %s", err)
				}

				metaExpected := make(map[string]interface{})
				metaExpected["test"] = "some_test"

				if commit.CommitSha != "11111111" ||
					commit.RepoName != "breakdown" ||
					commit.Meta["test"] != "some_test" {
					return fmt.Errorf("expected commit to be same (%+v)", commit)
				}
				return nil
			},
		},
	}

	qs, db, err := getQueries()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	defer db.Close(ctx)

	testMC := MyController{
		db:      db,
		queries: qs,
		l:       logger.NewMockCountLogger("test"),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input(db, testMC)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("got error: %s", err)
			}

		})
	}
}
