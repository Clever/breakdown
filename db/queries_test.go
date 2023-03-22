package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4"
)

type queryTest[K any] struct {
	name       string
	queries    func(*Queries, context.Context)
	getResults func(*Queries, context.Context) []K
	expected   []K
}

func getQueries() (*Queries, *pgx.Conn, error) {
	testDb, err := TestDB()
	if err != nil {
		return nil, nil, err
	}

	return New(testDb), testDb, nil
}

func TestHealthCheck(t *testing.T) {
	queries, db, err := getQueries()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	defer db.Close(ctx)

	if err = queries.HealthCheck(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestRepoQueries(t *testing.T) {
	tests := []queryTest[string]{
		{
			name: "inserts 2 new repos with one duplicate",
			queries: func(q *Queries, ctx context.Context) {
				q.CreateRepo(ctx, "github.com/Clever/breakdown")
				q.CreateRepo(ctx, "github.com/Clever/test")
				q.CreateRepo(ctx, "github.com/Clever/breakdown")
			},
			getResults: func(q *Queries, ctx context.Context) []string {
				repos, _ := q.ListRepos(ctx)
				res := []string{}
				for _, repo := range repos {
					res = append(res, repo.Name)
				}
				return res
			},
			expected: []string{
				"github.com/Clever/breakdown",
				"github.com/Clever/test",
			},
		},
		{
			name: "ListRepos returns repos ordered by name",
			queries: func(q *Queries, ctx context.Context) {
				q.CreateRepo(ctx, "github.com/Clever/test")
				q.CreateRepo(ctx, "github.com/Clever/breakdown")
			},
			getResults: func(q *Queries, ctx context.Context) []string {
				repos, _ := q.ListRepos(ctx)
				res := []string{}
				for _, repo := range repos {
					res = append(res, repo.Name)
				}
				return res
			},
			expected: []string{
				"github.com/Clever/breakdown",
				"github.com/Clever/test",
			},
		},
	}

	runQueryTest(t, tests)
}

func runQueryTest[K any](t *testing.T, tests []queryTest[K]) {
	t.Helper()

	queries, db, err := getQueries()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	defer db.Close(ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := db.Begin(ctx)
			if err != nil {
				t.Fatalf("getting initial tx: %s", err)
			}
			qtx := queries.WithTx(tx)
			tt.queries(qtx, ctx)

			results := tt.getResults(qtx, ctx)

			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("want=%+v\ngot= %+v", tt.expected, results)
			}

			if err = tx.Rollback(ctx); err != nil {
				t.Fatalf("rolling back transaction: %s", err)
			}
		})
	}
}
