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

// NewDeleteGroupParams creates a new DeleteGroupParams object
// with the default values initialized.
func NewDeleteGroupParams() *DeleteGroupParams {
	var ()
	return &DeleteGroupParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteGroupParamsWithTimeout creates a new DeleteGroupParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteGroupParamsWithTimeout(timeout time.Duration) *DeleteGroupParams {
	var ()
	return &DeleteGroupParams{

		timeout: timeout,
	}
}

// NewDeleteGroupParamsWithContext creates a new DeleteGroupParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteGroupParamsWithContext(ctx context.Context) *DeleteGroupParams {
	var ()
	return &DeleteGroupParams{

		Context: ctx,
	}
}

// NewDeleteGroupParamsWithHTTPClient creates a new DeleteGroupParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteGroupParamsWithHTTPClient(client *http.Client) *DeleteGroupParams {
	var ()
	return &DeleteGroupParams{
		HTTPClient: client,
	}
}

/*DeleteGroupParams contains all the parameters to send to the API endpoint
for the delete group operation typically these are written to a http.Request
*/
type DeleteGroupParams struct {

	/*ID
	  The ID of the group to delete

	*/
	ID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete group params
func (o *DeleteGroupParams) WithTimeout(timeout time.Duration) *DeleteGroupParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete group params
func (o *DeleteGroupParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete group params
func (o *DeleteGroupParams) WithContext(ctx context.Context) *DeleteGroupParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete group params
func (o *DeleteGroupParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete group params
func (o *DeleteGroupParams) WithHTTPClient(client *http.Client) *DeleteGroupParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete group params
func (o *DeleteGroupParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the delete group params
func (o *DeleteGroupParams) WithID(id int64) *DeleteGroupParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete group params
func (o *DeleteGroupParams) SetID(id int64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteGroupParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
