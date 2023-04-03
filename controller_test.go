package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/Clever/breakdown/db"
	"github.com/Clever/breakdown/gen-go/models"
	"github.com/Clever/kayvee-go/v7/logger"
	"github.com/go-openapi/swag"
	"github.com/jackc/pgx/v4"
)

func getQueries() (*db.Queries, *pgx.Conn, error) {
	testDb, err := db.TestDB()
	if err != nil {
		return nil, nil, err
	}

	return db.New(testDb), testDb, nil
}

type deployTest struct {
	name        string
	input       func(MyController) error
	expectError bool
	expected    func() []string
}

func TestDeploy(t *testing.T) {
	tests := []deployTest{
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
