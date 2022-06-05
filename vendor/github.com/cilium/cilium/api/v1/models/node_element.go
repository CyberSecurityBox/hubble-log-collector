// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2017-2021 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NodeElement Known node in the cluster
//
// +k8s:deepcopy-gen=true
//
// swagger:model NodeElement
type NodeElement struct {

	// Address used for probing cluster connectivity
	HealthEndpointAddress *NodeAddressing `json:"health-endpoint-address,omitempty"`

	// Name of the node including the cluster association. This is typically
	// <clustername>/<hostname>.
	//
	Name string `json:"name,omitempty"`

	// Primary address used for intra-cluster communication
	PrimaryAddress *NodeAddressing `json:"primary-address,omitempty"`

	// Alternative addresses assigned to the node
	SecondaryAddresses []*NodeAddressingElement `json:"secondary-addresses"`
}

// Validate validates this node element
func (m *NodeElement) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHealthEndpointAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePrimaryAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSecondaryAddresses(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NodeElement) validateHealthEndpointAddress(formats strfmt.Registry) error {

	if swag.IsZero(m.HealthEndpointAddress) { // not required
		return nil
	}

	if m.HealthEndpointAddress != nil {
		if err := m.HealthEndpointAddress.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("health-endpoint-address")
			}
			return err
		}
	}

	return nil
}

func (m *NodeElement) validatePrimaryAddress(formats strfmt.Registry) error {

	if swag.IsZero(m.PrimaryAddress) { // not required
		return nil
	}

	if m.PrimaryAddress != nil {
		if err := m.PrimaryAddress.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("primary-address")
			}
			return err
		}
	}

	return nil
}

func (m *NodeElement) validateSecondaryAddresses(formats strfmt.Registry) error {

	if swag.IsZero(m.SecondaryAddresses) { // not required
		return nil
	}

	for i := 0; i < len(m.SecondaryAddresses); i++ {
		if swag.IsZero(m.SecondaryAddresses[i]) { // not required
			continue
		}

		if m.SecondaryAddresses[i] != nil {
			if err := m.SecondaryAddresses[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("secondary-addresses" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *NodeElement) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NodeElement) UnmarshalBinary(b []byte) error {
	var res NodeElement
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
