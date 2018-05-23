// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// BlobAttributes Blob attributes set by the user
// swagger:model BlobAttributes
type BlobAttributes struct {

	// description
	Description string `json:"description,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// Storage path
	// Max Length: 200
	Path string `json:"path,omitempty"`

	// properties
	Properties interface{} `json:"properties,omitempty"`
}

// Validate validates this blob attributes
func (m *BlobAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePath(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BlobAttributes) validatePath(formats strfmt.Registry) error {

	if swag.IsZero(m.Path) { // not required
		return nil
	}

	if err := validate.MaxLength("path", "body", string(m.Path), 200); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *BlobAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BlobAttributes) UnmarshalBinary(b []byte) error {
	var res BlobAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}