// Code generated by go-swagger; DO NOT EDIT.

package access_proxies

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

// NewEditProxyParams creates a new EditProxyParams object
// with the default values initialized.
func NewEditProxyParams() *EditProxyParams {
	var ()
	return &EditProxyParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewEditProxyParamsWithTimeout creates a new EditProxyParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewEditProxyParamsWithTimeout(timeout time.Duration) *EditProxyParams {
	var ()
	return &EditProxyParams{

		timeout: timeout,
	}
}

// NewEditProxyParamsWithContext creates a new EditProxyParams object
// with the default values initialized, and the ability to set a context for a request
func NewEditProxyParamsWithContext(ctx context.Context) *EditProxyParams {
	var ()
	return &EditProxyParams{

		Context: ctx,
	}
}

// NewEditProxyParamsWithHTTPClient creates a new EditProxyParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewEditProxyParamsWithHTTPClient(client *http.Client) *EditProxyParams {
	var ()
	return &EditProxyParams{
		HTTPClient: client,
	}
}

/*EditProxyParams contains all the parameters to send to the API endpoint
for the edit proxy operation typically these are written to a http.Request
*/
type EditProxyParams struct {

	/*ID
	  The ID of the proxy to edit

	*/
	ID strfmt.UUID
	/*Proxy
	  Proxy information to modify

	*/
	Proxy EditProxyBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the edit proxy params
func (o *EditProxyParams) WithTimeout(timeout time.Duration) *EditProxyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the edit proxy params
func (o *EditProxyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the edit proxy params
func (o *EditProxyParams) WithContext(ctx context.Context) *EditProxyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the edit proxy params
func (o *EditProxyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the edit proxy params
func (o *EditProxyParams) WithHTTPClient(client *http.Client) *EditProxyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the edit proxy params
func (o *EditProxyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the edit proxy params
func (o *EditProxyParams) WithID(id strfmt.UUID) *EditProxyParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the edit proxy params
func (o *EditProxyParams) SetID(id strfmt.UUID) {
	o.ID = id
}

// WithProxy adds the proxy to the edit proxy params
func (o *EditProxyParams) WithProxy(proxy EditProxyBody) *EditProxyParams {
	o.SetProxy(proxy)
	return o
}

// SetProxy adds the proxy to the edit proxy params
func (o *EditProxyParams) SetProxy(proxy EditProxyBody) {
	o.Proxy = proxy
}

// WriteToRequest writes these params to a swagger request
func (o *EditProxyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID.String()); err != nil {
		return err
	}

	if err := r.SetBodyParam(o.Proxy); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}