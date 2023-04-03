// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// RepoCommit A repo commit
//
// swagger:model RepoCommit
type RepoCommit struct {

	// First 8 chars of commit SHA
	// Required: true
	CommitSha *string `json:"commit_sha"`

	// package files
	PackageFiles RepoPackageFiles `json:"package_files,omitempty"`

	// Full repo name "github.com/Clever/<name>"
	// Required: true
	RepoName *string `json:"repo_name"`
}

// Validate validates this repo commit
func (m *RepoCommit) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCommitSha(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePackageFiles(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRepoName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RepoCommit) validateCommitSha(formats strfmt.Registry) error {

	if err := validate.Required("commit_sha", "body", m.CommitSha); err != nil {
		return err
	}

	return nil
}

func (m *RepoCommit) validatePackageFiles(formats strfmt.Registry) error {

	if swag.IsZero(m.PackageFiles) { // not required
		return nil
	}

	if err := m.PackageFiles.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("package_files")
		}
		return err
	}

	return nil
}

func (m *RepoCommit) validateRepoName(formats strfmt.Registry) error {

	if err := validate.Required("repo_name", "body", m.RepoName); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RepoCommit) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RepoCommit) UnmarshalBinary(b []byte) error {
	var res RepoCommit
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
