// Code generated by go-swagger; DO NOT EDIT.

package ordersoperations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetOrdersCountByIssuerParams creates a new GetOrdersCountByIssuerParams object
// no default values defined in spec.
func NewGetOrdersCountByIssuerParams() GetOrdersCountByIssuerParams {

	return GetOrdersCountByIssuerParams{}
}

// GetOrdersCountByIssuerParams contains all the bound params for the get orders count by issuer operation
// typically these are obtained from a http.Request
//
// swagger:parameters getOrdersCountByIssuer
type GetOrdersCountByIssuerParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Search date-time end
	  In: query
	*/
	EndDateTime *strfmt.DateTime
	/*Organization that issued the Orders.
	  Required: true
	  In: path
	*/
	Issuer string
	/*Search date-time start
	  In: query
	*/
	StartDateTime *strfmt.DateTime
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetOrdersCountByIssuerParams() beforehand.
func (o *GetOrdersCountByIssuerParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qEndDateTime, qhkEndDateTime, _ := qs.GetOK("endDateTime")
	if err := o.bindEndDateTime(qEndDateTime, qhkEndDateTime, route.Formats); err != nil {
		res = append(res, err)
	}

	rIssuer, rhkIssuer, _ := route.Params.GetOK("issuer")
	if err := o.bindIssuer(rIssuer, rhkIssuer, route.Formats); err != nil {
		res = append(res, err)
	}

	qStartDateTime, qhkStartDateTime, _ := qs.GetOK("startDateTime")
	if err := o.bindStartDateTime(qStartDateTime, qhkStartDateTime, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindEndDateTime binds and validates parameter EndDateTime from query.
func (o *GetOrdersCountByIssuerParams) bindEndDateTime(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: date-time
	value, err := formats.Parse("date-time", raw)
	if err != nil {
		return errors.InvalidType("endDateTime", "query", "strfmt.DateTime", raw)
	}
	o.EndDateTime = (value.(*strfmt.DateTime))

	if err := o.validateEndDateTime(formats); err != nil {
		return err
	}

	return nil
}

// validateEndDateTime carries on validations for parameter EndDateTime
func (o *GetOrdersCountByIssuerParams) validateEndDateTime(formats strfmt.Registry) error {

	if err := validate.FormatOf("endDateTime", "query", "date-time", o.EndDateTime.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindIssuer binds and validates parameter Issuer from path.
func (o *GetOrdersCountByIssuerParams) bindIssuer(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Issuer = raw

	if err := o.validateIssuer(formats); err != nil {
		return err
	}

	return nil
}

// validateIssuer carries on validations for parameter Issuer
func (o *GetOrdersCountByIssuerParams) validateIssuer(formats strfmt.Registry) error {

	if err := validate.Enum("issuer", "path", o.Issuer, []interface{}{"army", "navy", "air-force", "marine-corps", "coast-guard"}); err != nil {
		return err
	}

	return nil
}

// bindStartDateTime binds and validates parameter StartDateTime from query.
func (o *GetOrdersCountByIssuerParams) bindStartDateTime(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: date-time
	value, err := formats.Parse("date-time", raw)
	if err != nil {
		return errors.InvalidType("startDateTime", "query", "strfmt.DateTime", raw)
	}
	o.StartDateTime = (value.(*strfmt.DateTime))

	if err := o.validateStartDateTime(formats); err != nil {
		return err
	}

	return nil
}

// validateStartDateTime carries on validations for parameter StartDateTime
func (o *GetOrdersCountByIssuerParams) validateStartDateTime(formats strfmt.Registry) error {

	if err := validate.FormatOf("startDateTime", "query", "date-time", o.StartDateTime.String(), formats); err != nil {
		return err
	}
	return nil
}
