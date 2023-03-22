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

// RepoPackageFile format of packages
//
// swagger:model RepoPackageFile
type RepoPackageFile struct {

	// error when parsing package-file, if any
	Error string `json:"error,omitempty"`

	// version of go, if any
	GoVersion string `json:"go_version,omitempty"`

	// Name of go module or npm package
	// Required: true
	Name *string `json:"name"`

	// packages
	// Required: true
	Packages *RepoPackages `json:"packages"`

	// path to package file eg "go.mod"
	// Required: true
	Path *string `json:"path"`
}

// Validate validates this repo package file
func (m *RepoPackageFile) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePackages(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePath(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RepoPackageFile) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *RepoPackageFile) validatePackages(formats strfmt.Registry) error {

	if err := validate.Required("packages", "body", m.Packages); err != nil {
		return err
	}

	if m.Packages != nil {
		if err := m.Packages.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("packages")
			}
			return err
		}
	}

	return nil
}

func (m *RepoPackageFile) validatePath(formats strfmt.Registry) error {

	if err := validate.Required("path", "body", m.Path); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RepoPackageFile) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RepoPackageFile) UnmarshalBinary(b []byte) error {
	var res RepoPackageFile
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
