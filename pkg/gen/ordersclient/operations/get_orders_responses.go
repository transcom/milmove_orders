// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/transcom/milmove_orders/pkg/gen/ordersmessages"
)

// GetOrdersReader is a Reader for the GetOrders structure.
type GetOrdersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetOrdersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetOrdersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetOrdersBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewGetOrdersUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetOrdersForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetOrdersNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetOrdersInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetOrdersOK creates a GetOrdersOK with default headers values
func NewGetOrdersOK() *GetOrdersOK {
	return &GetOrdersOK{}
}

/*GetOrdersOK handles this case with default header values.

Successful
*/
type GetOrdersOK struct {
	Payload *ordersmessages.Orders
}

func (o *GetOrdersOK) Error() string {
	return fmt.Sprintf("[GET /orders/{uuid}][%d] getOrdersOK  %+v", 200, o.Payload)
}

func (o *GetOrdersOK) GetPayload() *ordersmessages.Orders {
	return o.Payload
}

func (o *GetOrdersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(ordersmessages.Orders)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetOrdersBadRequest creates a GetOrdersBadRequest with default headers values
func NewGetOrdersBadRequest() *GetOrdersBadRequest {
	return &GetOrdersBadRequest{}
}

/*GetOrdersBadRequest handles this case with default header values.

Invalid
*/
type GetOrdersBadRequest struct {
}

func (o *GetOrdersBadRequest) Error() string {
	return fmt.Sprintf("[GET /orders/{uuid}][%d] getOrdersBadRequest ", 400)
}

func (o *GetOrdersBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetOrdersUnauthorized creates a GetOrdersUnauthorized with default headers values
func NewGetOrdersUnauthorized() *GetOrdersUnauthorized {
	return &GetOrdersUnauthorized{}
}

/*GetOrdersUnauthorized handles this case with default header values.

must be authenticated to use this endpoint
*/
type GetOrdersUnauthorized struct {
}

func (o *GetOrdersUnauthorized) Error() string {
	return fmt.Sprintf("[GET /orders/{uuid}][%d] getOrdersUnauthorized ", 401)
}

func (o *GetOrdersUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetOrdersForbidden creates a GetOrdersForbidden with default headers values
func NewGetOrdersForbidden() *GetOrdersForbidden {
	return &GetOrdersForbidden{}
}

/*GetOrdersForbidden handles this case with default header values.

Forbidden
*/
type GetOrdersForbidden struct {
}

func (o *GetOrdersForbidden) Error() string {
	return fmt.Sprintf("[GET /orders/{uuid}][%d] getOrdersForbidden ", 403)
}

func (o *GetOrdersForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetOrdersNotFound creates a GetOrdersNotFound with default headers values
func NewGetOrdersNotFound() *GetOrdersNotFound {
	return &GetOrdersNotFound{}
}

/*GetOrdersNotFound handles this case with default header values.

Orders not found
*/
type GetOrdersNotFound struct {
}

func (o *GetOrdersNotFound) Error() string {
	return fmt.Sprintf("[GET /orders/{uuid}][%d] getOrdersNotFound ", 404)
}

func (o *GetOrdersNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetOrdersInternalServerError creates a GetOrdersInternalServerError with default headers values
func NewGetOrdersInternalServerError() *GetOrdersInternalServerError {
	return &GetOrdersInternalServerError{}
}

/*GetOrdersInternalServerError handles this case with default header values.

Server error
*/
type GetOrdersInternalServerError struct {
}

func (o *GetOrdersInternalServerError) Error() string {
	return fmt.Sprintf("[GET /orders/{uuid}][%d] getOrdersInternalServerError ", 500)
}

func (o *GetOrdersInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}