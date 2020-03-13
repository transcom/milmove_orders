// Code generated by go-swagger; DO NOT EDIT.

package operations

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

// NewIndexOrdersForMemberParams creates a new IndexOrdersForMemberParams object
// with the default values initialized.
func NewIndexOrdersForMemberParams() *IndexOrdersForMemberParams {
	var ()
	return &IndexOrdersForMemberParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewIndexOrdersForMemberParamsWithTimeout creates a new IndexOrdersForMemberParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewIndexOrdersForMemberParamsWithTimeout(timeout time.Duration) *IndexOrdersForMemberParams {
	var ()
	return &IndexOrdersForMemberParams{

		timeout: timeout,
	}
}

// NewIndexOrdersForMemberParamsWithContext creates a new IndexOrdersForMemberParams object
// with the default values initialized, and the ability to set a context for a request
func NewIndexOrdersForMemberParamsWithContext(ctx context.Context) *IndexOrdersForMemberParams {
	var ()
	return &IndexOrdersForMemberParams{

		Context: ctx,
	}
}

// NewIndexOrdersForMemberParamsWithHTTPClient creates a new IndexOrdersForMemberParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewIndexOrdersForMemberParamsWithHTTPClient(client *http.Client) *IndexOrdersForMemberParams {
	var ()
	return &IndexOrdersForMemberParams{
		HTTPClient: client,
	}
}

/*IndexOrdersForMemberParams contains all the parameters to send to the API endpoint
for the index orders for member operation typically these are written to a http.Request
*/
type IndexOrdersForMemberParams struct {

	/*Edipi
	  EDIPI of the member to retrieve Orders

	*/
	Edipi string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the index orders for member params
func (o *IndexOrdersForMemberParams) WithTimeout(timeout time.Duration) *IndexOrdersForMemberParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the index orders for member params
func (o *IndexOrdersForMemberParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the index orders for member params
func (o *IndexOrdersForMemberParams) WithContext(ctx context.Context) *IndexOrdersForMemberParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the index orders for member params
func (o *IndexOrdersForMemberParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the index orders for member params
func (o *IndexOrdersForMemberParams) WithHTTPClient(client *http.Client) *IndexOrdersForMemberParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the index orders for member params
func (o *IndexOrdersForMemberParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithEdipi adds the edipi to the index orders for member params
func (o *IndexOrdersForMemberParams) WithEdipi(edipi string) *IndexOrdersForMemberParams {
	o.SetEdipi(edipi)
	return o
}

// SetEdipi adds the edipi to the index orders for member params
func (o *IndexOrdersForMemberParams) SetEdipi(edipi string) {
	o.Edipi = edipi
}

// WriteToRequest writes these params to a swagger request
func (o *IndexOrdersForMemberParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param edipi
	if err := r.SetPathParam("edipi", o.Edipi); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
