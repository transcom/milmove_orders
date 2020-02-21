// Code generated by go-swagger; DO NOT EDIT.

package ordersoperations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	ordersmessages "github.com/transcom/milmove_orders/pkg/gen/ordersmessages"
)

// IndexOrdersForMemberOKCode is the HTTP code returned for type IndexOrdersForMemberOK
const IndexOrdersForMemberOKCode int = 200

/*IndexOrdersForMemberOK Successful

swagger:response indexOrdersForMemberOK
*/
type IndexOrdersForMemberOK struct {

	/*
	  In: Body
	*/
	Payload []*ordersmessages.Orders `json:"body,omitempty"`
}

// NewIndexOrdersForMemberOK creates IndexOrdersForMemberOK with default headers values
func NewIndexOrdersForMemberOK() *IndexOrdersForMemberOK {

	return &IndexOrdersForMemberOK{}
}

// WithPayload adds the payload to the index orders for member o k response
func (o *IndexOrdersForMemberOK) WithPayload(payload []*ordersmessages.Orders) *IndexOrdersForMemberOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the index orders for member o k response
func (o *IndexOrdersForMemberOK) SetPayload(payload []*ordersmessages.Orders) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *IndexOrdersForMemberOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*ordersmessages.Orders, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// IndexOrdersForMemberBadRequestCode is the HTTP code returned for type IndexOrdersForMemberBadRequest
const IndexOrdersForMemberBadRequestCode int = 400

/*IndexOrdersForMemberBadRequest Bad request

swagger:response indexOrdersForMemberBadRequest
*/
type IndexOrdersForMemberBadRequest struct {
}

// NewIndexOrdersForMemberBadRequest creates IndexOrdersForMemberBadRequest with default headers values
func NewIndexOrdersForMemberBadRequest() *IndexOrdersForMemberBadRequest {

	return &IndexOrdersForMemberBadRequest{}
}

// WriteResponse to the client
func (o *IndexOrdersForMemberBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// IndexOrdersForMemberUnauthorizedCode is the HTTP code returned for type IndexOrdersForMemberUnauthorized
const IndexOrdersForMemberUnauthorizedCode int = 401

/*IndexOrdersForMemberUnauthorized must be authenticated to use this endpoint

swagger:response indexOrdersForMemberUnauthorized
*/
type IndexOrdersForMemberUnauthorized struct {
}

// NewIndexOrdersForMemberUnauthorized creates IndexOrdersForMemberUnauthorized with default headers values
func NewIndexOrdersForMemberUnauthorized() *IndexOrdersForMemberUnauthorized {

	return &IndexOrdersForMemberUnauthorized{}
}

// WriteResponse to the client
func (o *IndexOrdersForMemberUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// IndexOrdersForMemberForbiddenCode is the HTTP code returned for type IndexOrdersForMemberForbidden
const IndexOrdersForMemberForbiddenCode int = 403

/*IndexOrdersForMemberForbidden Forbidden

swagger:response indexOrdersForMemberForbidden
*/
type IndexOrdersForMemberForbidden struct {
}

// NewIndexOrdersForMemberForbidden creates IndexOrdersForMemberForbidden with default headers values
func NewIndexOrdersForMemberForbidden() *IndexOrdersForMemberForbidden {

	return &IndexOrdersForMemberForbidden{}
}

// WriteResponse to the client
func (o *IndexOrdersForMemberForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// IndexOrdersForMemberNotFoundCode is the HTTP code returned for type IndexOrdersForMemberNotFound
const IndexOrdersForMemberNotFoundCode int = 404

/*IndexOrdersForMemberNotFound No orders found

swagger:response indexOrdersForMemberNotFound
*/
type IndexOrdersForMemberNotFound struct {
}

// NewIndexOrdersForMemberNotFound creates IndexOrdersForMemberNotFound with default headers values
func NewIndexOrdersForMemberNotFound() *IndexOrdersForMemberNotFound {

	return &IndexOrdersForMemberNotFound{}
}

// WriteResponse to the client
func (o *IndexOrdersForMemberNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// IndexOrdersForMemberInternalServerErrorCode is the HTTP code returned for type IndexOrdersForMemberInternalServerError
const IndexOrdersForMemberInternalServerErrorCode int = 500

/*IndexOrdersForMemberInternalServerError Server error

swagger:response indexOrdersForMemberInternalServerError
*/
type IndexOrdersForMemberInternalServerError struct {
}

// NewIndexOrdersForMemberInternalServerError creates IndexOrdersForMemberInternalServerError with default headers values
func NewIndexOrdersForMemberInternalServerError() *IndexOrdersForMemberInternalServerError {

	return &IndexOrdersForMemberInternalServerError{}
}

// WriteResponse to the client
func (o *IndexOrdersForMemberInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
