package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Clever/breakdown/db"
	"github.com/Clever/breakdown/gen-go/models"
	"github.com/Clever/breakdown/gen-go/server"
	"github.com/Clever/kayvee-go/v7/logger"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// MyController implements server.Controller
type MyController struct {
	launchConfig LaunchConfig
	db           *pgx.Conn
	queries      *db.Queries
	l            logger.KayveeLogger
}

var _ server.Controller = MyController{}

func (mc MyController) beginTX(ctx context.Context) (pgx.Tx, *db.Queries, error) {
	tx, err := mc.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	qtx := mc.queries.WithTx(tx)
	return tx, qtx, nil
}

// HealthCheck handles GET requests to /_health
func (mc MyController) HealthCheck(ctx context.Context) error {
	tx, qtx, err := mc.beginTX(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	if err = qtx.HealthCheck(ctx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// PostUpload handles POSTs to /v1/upload
func (mc MyController) PostUpload(ctx context.Context, i *models.RepoCommit) error {
	tx, qtx, err := mc.beginTX(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %s", err)
	}

	repoID, err := qtx.CreateRepo(ctx, *i.RepoName)
	if err != nil {
		return fmt.Errorf("creating repo: %s", err)
	}

	repoCommitID, err := qtx.CreateRepoCommit(ctx, db.CreateRepoCommitParams{
		RepoID:     repoID,
		CommitSha:  *i.CommitSha,
		CommitDate: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("creating repo commit: %s", err)
	}

	for _, packageFile := range i.PackageFiles {

		if len(packageFile.Error) > 0 {
			mc.l.ErrorD("parse-error", logger.M{
				"repo":    *i.RepoName,
				"sha":     *i.CommitSha,
				"path":    *packageFile.Path,
				"message": packageFile.Error,
			})
			continue
		}

		// meta := pgtype.JSONB{}
		var meta pgtype.JSONB
		meta.Set(nil)
		var packageType db.PackageType
		switch *packageFile.Type {
		case "gomod":
			packageType = db.PackageTypeGomod
		case "npm":
			packageType = db.PackageTypeNpm
		default:
			return fmt.Errorf("unknown package type %q for %q", *packageFile.Type, *packageFile.Path)
		}

		if len(packageFile.GoVersion) > 0 {
			metaBytes, err := json.Marshal(struct {
				Goversion string `json:"go_version"`
			}{
				Goversion: packageFile.GoVersion,
			})
			if err != nil {
				return err
			}
			meta.Set(metaBytes)
		}

		fileID, err := qtx.CreatePackageFile(ctx, db.CreatePackageFileParams{
			RepoCommitID: repoCommitID,
			Path:         *packageFile.Path,
			Type:         packageType,
			Meta:         meta,
		})

		if err != nil {
			return fmt.Errorf("creating package file: %s", err)
		}

		// First pass over Packages to insert new dependencies and store id to reference later
		depNameToID := map[string]int64{}

		packageFileDepName := fmt.Sprintf("%s@%s", packageFile.Name, packageFile.GoVersion)

		createDepParams := make([]db.CreateDependencyParams, 0)

		for depNameVer, depInfo := range packageFile.Packages {
			// Skip over top level module/package
			if depNameVer == packageFileDepName {
				continue
			}
			createDepParams = append(createDepParams, db.CreateDependencyParams{
				Name:    depInfo.Name,
				Version: depInfo.Version,
				Type:    packageType,
				IsLocal: depInfo.IsLocal,
			})
		}

		err = nil
		depRes := qtx.CreateDependency(ctx, createDepParams)
		depRes.Query(func(i int, deps []db.Dependency, batchErr error) {
			if batchErr != nil {
				err = fmt.Errorf("batching dependencies: %s", batchErr.Error())
				return
			}
			for _, dep := range deps {
				name := fmt.Sprintf("%s@%s", dep.Name, dep.Version)
				depNameToID[name] = dep.ID
			}
		})
		if err != nil {
			return err
		}

		// Insert direct dependencies
		packageDepInfo, ok := packageFile.Packages[packageFileDepName]
		if !ok {
			return fmt.Errorf("top level module/package %q not found in packages", packageFileDepName)
		}

		packageFileDepParams := make([]db.InsertPackageFileDependencyParams, 0)
		for _, directDep := range packageDepInfo.Dependencies {
			depID, ok := depNameToID[directDep]
			if !ok {
				return fmt.Errorf("dependency ID not found for %q (file %q)", directDep, packageFileDepName)
			}

			packageFileDepParams = append(packageFileDepParams, db.InsertPackageFileDependencyParams{
				PackageFileID: fileID,
				DependencyID:  depID,
			})

		}

		err = nil
		batchRes := qtx.InsertPackageFileDependency(ctx, packageFileDepParams)
		batchRes.Exec(func(i int, batchErr error) {
			if batchErr != nil {
				err = fmt.Errorf("batching package file deps: %s", batchErr.Error())
			}
		})
		if err != nil {
			return err
		}

		// Second pass over packages to insert dep <-> dep

		insertDepDepParams := make([]db.InsertDepDependencyParams, 0)

		for depNameVer, depInfo := range packageFile.Packages {
			if depNameVer == packageFileDepName {
				continue
			}
			parentID, ok := depNameToID[depNameVer]
			if !ok {
				return fmt.Errorf("parent ID not found for %q", depNameVer)
			}

			for _, depDep := range depInfo.Dependencies {
				depID, ok := depNameToID[depDep]
				if !ok {
					return fmt.Errorf("%q -> %q dep id not found", depNameVer, depDep)
				}
				insertDepDepParams = append(insertDepDepParams, db.InsertDepDependencyParams{
					ParentID:     parentID,
					DependencyID: depID,
				})
			}
		}

		err = nil
		depBatchRes := qtx.InsertDepDependency(ctx, insertDepDepParams)
		depBatchRes.Exec(func(i int, execErr error) {
			if execErr != nil {
				err = fmt.Errorf("batching dep deps: %s", execErr.Error())
			}
		})
		if err != nil {
			return err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing: %s", err.Error())
	}
	return nil
}

// GetCommit handles GETs to /v1/commit
func (mc MyController) GetCommit(ctx context.Context, i *models.GetCommitInformation) (*models.CommitInformation, error) {
	qs := mc.queries

	if len(*i.CommitSha) < 8 {
		return nil, models.BadRequest{Message: "commit_sha not long enough"}
	}

	commit, err := qs.GetCommit(ctx, db.GetCommitParams{
		Name:      *i.RepoName,
		CommitSha: *i.CommitSha,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NotFound{Message: "repo_commit not found"}
		}
		return nil, models.BadRequest{Message: err.Error()}
	}

	var meta models.JSONObject
	if commit.Meta.Status&pgtype.Present > 0 {
		err = commit.Meta.AssignTo(&meta)
		if err != nil {
			return nil, err
		}
	}

	return &models.CommitInformation{
		CommitSha: commit.CommitSha,
		RepoName:  *i.RepoName,
		Meta:      meta,
	}, nil

}

// PostCustom handles POSTs to /v1/custom
func (mc MyController) PostCustom(ctx context.Context, i *models.CustomData) error {
	return nil
}

// PostDeploy handles POSTs to /v1/deploy
func (mc MyController) PostDeploy(ctx context.Context, deploys *models.Deploys) error {
	tx, qtx, err := mc.beginTX(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	deploymentParams := make([]db.InsertDeploymentParams, 0)
	for _, deploy := range *deploys {
		if len(*deploy.CommitSha) < 8 {
			return fmt.Errorf("%q not long enough, at least 8 chars", *deploy.CommitSha)
		}
		deploymentParams = append(deploymentParams, db.InsertDeploymentParams{
			CommitSha:   (*deploy.CommitSha)[0:8],
			Application: *deploy.Application,
			Environment: *deploy.Environment,
			RunType:     *deploy.RunType,
			Version:     *deploy.Version,
		})
	}

	batchRes := qtx.InsertDeployment(ctx, deploymentParams)
	batchRes.Exec(func(i int, batchErr error) {
		if batchErr != nil {
			err = fmt.Errorf("inserting deployments: %s", batchErr)
		}
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
