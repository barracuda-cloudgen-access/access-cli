// Code generated by go-swagger; DO NOT EDIT.

package access_resources

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewCreateResourceParams creates a new CreateResourceParams object
// with the default values initialized.
func NewCreateResourceParams() *CreateResourceParams {
	var ()
	return &CreateResourceParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateResourceParamsWithTimeout creates a new CreateResourceParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateResourceParamsWithTimeout(timeout time.Duration) *CreateResourceParams {
	var ()
	return &CreateResourceParams{

		timeout: timeout,
	}
}

// NewCreateResourceParamsWithContext creates a new CreateResourceParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateResourceParamsWithContext(ctx context.Context) *CreateResourceParams {
	var ()
	return &CreateResourceParams{

		Context: ctx,
	}
}

// NewCreateResourceParamsWithHTTPClient creates a new CreateResourceParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateResourceParamsWithHTTPClient(client *http.Client) *CreateResourceParams {
	var ()
	return &CreateResourceParams{
		HTTPClient: client,
	}
}

/*CreateResourceParams contains all the parameters to send to the API endpoint
for the create resource operation typically these are written to a http.Request
*/
type CreateResourceParams struct {

	/*Resource
	  Information about the resource to create

	*/
	Resource CreateResourceBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create resource params
func (o *CreateResourceParams) WithTimeout(timeout time.Duration) *CreateResourceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create resource params
func (o *CreateResourceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create resource params
func (o *CreateResourceParams) WithContext(ctx context.Context) *CreateResourceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create resource params
func (o *CreateResourceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create resource params
func (o *CreateResourceParams) WithHTTPClient(client *http.Client) *CreateResourceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create resource params
func (o *CreateResourceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithResource adds the resource to the create resource params
func (o *CreateResourceParams) WithResource(resource CreateResourceBody) *CreateResourceParams {
	o.SetResource(resource)
	return o
}

// SetResource adds the resource to the create resource params
func (o *CreateResourceParams) SetResource(resource CreateResourceBody) {
	o.Resource = resource
}

// WriteToRequest writes these params to a swagger request
func (o *CreateResourceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.Resource); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}