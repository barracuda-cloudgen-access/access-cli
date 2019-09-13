// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetGroupParams creates a new GetGroupParams object
// with the default values initialized.
func NewGetGroupParams() *GetGroupParams {
	var ()
	return &GetGroupParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetGroupParamsWithTimeout creates a new GetGroupParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetGroupParamsWithTimeout(timeout time.Duration) *GetGroupParams {
	var ()
	return &GetGroupParams{

		timeout: timeout,
	}
}

// NewGetGroupParamsWithContext creates a new GetGroupParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetGroupParamsWithContext(ctx context.Context) *GetGroupParams {
	var ()
	return &GetGroupParams{

		Context: ctx,
	}
}

// NewGetGroupParamsWithHTTPClient creates a new GetGroupParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetGroupParamsWithHTTPClient(client *http.Client) *GetGroupParams {
	var ()
	return &GetGroupParams{
		HTTPClient: client,
	}
}

/*GetGroupParams contains all the parameters to send to the API endpoint
for the get group operation typically these are written to a http.Request
*/
type GetGroupParams struct {

	/*ID
	  The ID of the group to retrieve

	*/
	ID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get group params
func (o *GetGroupParams) WithTimeout(timeout time.Duration) *GetGroupParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get group params
func (o *GetGroupParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get group params
func (o *GetGroupParams) WithContext(ctx context.Context) *GetGroupParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get group params
func (o *GetGroupParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get group params
func (o *GetGroupParams) WithHTTPClient(client *http.Client) *GetGroupParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get group params
func (o *GetGroupParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get group params
func (o *GetGroupParams) WithID(id int64) *GetGroupParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get group params
func (o *GetGroupParams) SetID(id int64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetGroupParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
