// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/oNaiPs/fyde-cli/models"
)

// EditUserReader is a Reader for the EditUser structure.
type EditUserReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *EditUserReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewEditUserOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewEditUserUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewEditUserNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewEditUserOK creates a EditUserOK with default headers values
func NewEditUserOK() *EditUserOK {
	return &EditUserOK{}
}

/*EditUserOK handles this case with default header values.

User edited
*/
type EditUserOK struct {
	Payload *EditUserOKBody
}

func (o *EditUserOK) Error() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] editUserOK  %+v", 200, o.Payload)
}

func (o *EditUserOK) GetPayload() *EditUserOKBody {
	return o.Payload
}

func (o *EditUserOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(EditUserOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewEditUserUnauthorized creates a EditUserUnauthorized with default headers values
func NewEditUserUnauthorized() *EditUserUnauthorized {
	return &EditUserUnauthorized{}
}

/*EditUserUnauthorized handles this case with default header values.

unauthorized: invalid credentials or missing authentication headers
*/
type EditUserUnauthorized struct {
	Payload *models.GenericUnauthorizedResponse
}

func (o *EditUserUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] editUserUnauthorized  %+v", 401, o.Payload)
}

func (o *EditUserUnauthorized) GetPayload() *models.GenericUnauthorizedResponse {
	return o.Payload
}

func (o *EditUserUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericUnauthorizedResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewEditUserNotFound creates a EditUserNotFound with default headers values
func NewEditUserNotFound() *EditUserNotFound {
	return &EditUserNotFound{}
}

/*EditUserNotFound handles this case with default header values.

user not found
*/
type EditUserNotFound struct {
}

func (o *EditUserNotFound) Error() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] editUserNotFound ", 404)
}

func (o *EditUserNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*EditUserBody edit user body
swagger:model EditUserBody
*/
type EditUserBody struct {

	// user
	User *EditUserParamsBodyUser `json:"user,omitempty"`
}

// Validate validates this edit user body
func (o *EditUserBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateUser(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *EditUserBody) validateUser(formats strfmt.Registry) error {

	if swag.IsZero(o.User) { // not required
		return nil
	}

	if o.User != nil {
		if err := o.User.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("user" + "." + "user")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *EditUserBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *EditUserBody) UnmarshalBinary(b []byte) error {
	var res EditUserBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*EditUserOKBody edit user o k body
swagger:model EditUserOKBody
*/
type EditUserOKBody struct {
	models.User

	// devices
	Devices []interface{} `json:"devices"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *EditUserOKBody) UnmarshalJSON(raw []byte) error {
	// EditUserOKBodyAO0
	var editUserOKBodyAO0 models.User
	if err := swag.ReadJSON(raw, &editUserOKBodyAO0); err != nil {
		return err
	}
	o.User = editUserOKBodyAO0

	// EditUserOKBodyAO1
	var dataEditUserOKBodyAO1 struct {
		Devices []interface{} `json:"devices"`
	}
	if err := swag.ReadJSON(raw, &dataEditUserOKBodyAO1); err != nil {
		return err
	}

	o.Devices = dataEditUserOKBodyAO1.Devices

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o EditUserOKBody) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	editUserOKBodyAO0, err := swag.WriteJSON(o.User)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, editUserOKBodyAO0)

	var dataEditUserOKBodyAO1 struct {
		Devices []interface{} `json:"devices"`
	}

	dataEditUserOKBodyAO1.Devices = o.Devices

	jsonDataEditUserOKBodyAO1, errEditUserOKBodyAO1 := swag.WriteJSON(dataEditUserOKBodyAO1)
	if errEditUserOKBodyAO1 != nil {
		return nil, errEditUserOKBodyAO1
	}
	_parts = append(_parts, jsonDataEditUserOKBodyAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this edit user o k body
func (o *EditUserOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with models.User
	if err := o.User.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (o *EditUserOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *EditUserOKBody) UnmarshalBinary(b []byte) error {
	var res EditUserOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*EditUserParamsBodyUser edit user params body user
swagger:model EditUserParamsBodyUser
*/
type EditUserParamsBodyUser struct {

	// email
	// Format: email
	Email strfmt.Email `json:"email,omitempty"`

	// enabled
	Enabled bool `json:"enabled,omitempty"`

	// group ids
	GroupIds []int64 `json:"group_ids"`

	// name
	Name string `json:"name,omitempty"`

	// phone number
	PhoneNumber string `json:"phone_number,omitempty"`
}

// Validate validates this edit user params body user
func (o *EditUserParamsBodyUser) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *EditUserParamsBodyUser) validateEmail(formats strfmt.Registry) error {

	if swag.IsZero(o.Email) { // not required
		return nil
	}

	if err := validate.FormatOf("user"+"."+"user"+"."+"email", "body", "email", o.Email.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *EditUserParamsBodyUser) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *EditUserParamsBodyUser) UnmarshalBinary(b []byte) error {
	var res EditUserParamsBodyUser
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
