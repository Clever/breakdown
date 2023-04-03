// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

type CiSource string

const (
	CiSourceCircleCi      CiSource = "circle-ci"
	CiSourceGithubActions CiSource = "github-actions"
)

func (e *CiSource) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CiSource(s)
	case string:
		*e = CiSource(s)
	default:
		return fmt.Errorf("unsupported scan type for CiSource: %T", src)
	}
	return nil
}

type NullCiSource struct {
	CiSource CiSource
	Valid    bool // Valid is true if CiSource is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCiSource) Scan(value interface{}) error {
	if value == nil {
		ns.CiSource, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CiSource.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCiSource) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CiSource), nil
}

type PackageType string

const (
	PackageTypeGomod PackageType = "gomod"
	PackageTypeNpm   PackageType = "npm"
)

func (e *PackageType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PackageType(s)
	case string:
		*e = PackageType(s)
	default:
		return fmt.Errorf("unsupported scan type for PackageType: %T", src)
	}
	return nil
}

type NullPackageType struct {
	PackageType PackageType
	Valid       bool // Valid is true if PackageType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPackageType) Scan(value interface{}) error {
	if value == nil {
		ns.PackageType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PackageType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPackageType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PackageType), nil
}

type CiWorkflow struct {
	ID           int64
	Source       CiSource
	SourceID     string
	RepoID       int64
	CommitSha    string
	WorkflowTime time.Time
	DurationS    int32
	Repo         sql.NullString
}

type DepDependency struct {
	ParentID     int64
	DependencyID int64
}

type Dependency struct {
	ID      int64
	Name    string
	Version string
	Type    PackageType
	IsLocal bool
}

type Deployment struct {
	ID          int64
	Application string
	Version     string
	Date        time.Time
	RunType     string
	Environment string
	CommitSha   string
}

type PackageFile struct {
	ID           int64
	RepoCommitID int64
	Path         string
	Type         PackageType
	Meta         pgtype.JSONB
}

type PackageFileDependency struct {
	PackageFileID int64
	DependencyID  int64
}

type Repo struct {
	ID   int64
	Name string
}

type RepoCommit struct {
	ID         int64
	RepoID     int64
	CommitSha  string
	CommitDate time.Time
	Meta       pgtype.JSONB
}