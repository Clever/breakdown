// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CommitInformation commit information
//
// swagger:model CommitInformation
type CommitInformation struct {

	// commit sha
	CommitSha string `json:"commit_sha,omitempty"`

	// meta
	Meta JSONObject `json:"meta,omitempty"`

	// repo name
	RepoName string `json:"repo_name,omitempty"`
}

// Validate validates this commit information
func (m *CommitInformation) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMeta(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CommitInformation) validateMeta(formats strfmt.Registry) error {

	if swag.IsZero(m.Meta) { // not required
		return nil
	}

	if err := m.Meta.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("meta")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CommitInformation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CommitInformation) UnmarshalBinary(b []byte) error {
	var res CommitInformation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
