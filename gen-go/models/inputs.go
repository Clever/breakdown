package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = fmt.Sprintf
var _ = url.QueryEscape
var _ = strconv.FormatInt
var _ = strings.Replace
var _ = validate.Maximum
var _ = strfmt.NewFormats

// HealthCheckInput holds the input parameters for a healthCheck operation.
type HealthCheckInput struct {
}

// Validate returns an error if any of the HealthCheckInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i HealthCheckInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i HealthCheckInput) Path() (string, error) {
	path := "/_health"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// GetThingsInput holds the input parameters for a getThings operation.
type GetThingsInput struct {
}

// Validate returns an error if any of the GetThingsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetThingsInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i GetThingsInput) Path() (string, error) {
	path := "/v2/things"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}

// DeleteThingInput holds the input parameters for a deleteThing operation.
type DeleteThingInput struct {
	ID string
}

// ValidateDeleteThingInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateDeleteThingInput(id string) error {

	return nil
}

// DeleteThingInputPath returns the URI path for the input.
func DeleteThingInputPath(id string) (string, error) {
	path := "/v2/things/{id}"
	urlVals := url.Values{}

	pathid := id
	if pathid == "" {
		err := fmt.Errorf("id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{id}", pathid, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetThingInput holds the input parameters for a getThing operation.
type GetThingInput struct {
	ID string
}

// ValidateGetThingInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetThingInput(id string) error {

	return nil
}

// GetThingInputPath returns the URI path for the input.
func GetThingInputPath(id string) (string, error) {
	path := "/v2/things/{id}"
	urlVals := url.Values{}

	pathid := id
	if pathid == "" {
		err := fmt.Errorf("id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{id}", pathid, -1)

	return path + "?" + urlVals.Encode(), nil
}

// CreateOrUpdateThingInput holds the input parameters for a createOrUpdateThing operation.
type CreateOrUpdateThingInput struct {
	Thing *Thing
	ID    string
}

// Validate returns an error if any of the CreateOrUpdateThingInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i CreateOrUpdateThingInput) Validate() error {

	if i.Thing != nil {
		if err := i.Thing.Validate(nil); err != nil {
			return err
		}
	}

	return nil
}

// Path returns the URI path for the input.
func (i CreateOrUpdateThingInput) Path() (string, error) {
	path := "/v2/things/{id}"
	urlVals := url.Values{}

	pathid := i.ID
	if pathid == "" {
		err := fmt.Errorf("id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{id}", pathid, -1)

	return path + "?" + urlVals.Encode(), nil
}
