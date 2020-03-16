// Code generated by go-swagger; DO NOT EDIT.

package ordersoperations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetOrdersCountByIssuerHandlerFunc turns a function with the right signature into a get orders count by issuer handler
type GetOrdersCountByIssuerHandlerFunc func(GetOrdersCountByIssuerParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetOrdersCountByIssuerHandlerFunc) Handle(params GetOrdersCountByIssuerParams) middleware.Responder {
	return fn(params)
}

// GetOrdersCountByIssuerHandler interface for that can handle valid get orders count by issuer params
type GetOrdersCountByIssuerHandler interface {
	Handle(GetOrdersCountByIssuerParams) middleware.Responder
}

// NewGetOrdersCountByIssuer creates a new http.Handler for the get orders count by issuer operation
func NewGetOrdersCountByIssuer(ctx *middleware.Context, handler GetOrdersCountByIssuerHandler) *GetOrdersCountByIssuer {
	return &GetOrdersCountByIssuer{Context: ctx, Handler: handler}
}

/*GetOrdersCountByIssuer swagger:route GET /issuers/{issuer}/count getOrdersCountByIssuer

Retrieve a count of Orders by issuer

Gets a Count of Orders by issuer.
## Errors
Users of this endpoint must have permission to read Orders for the `issuer` associated with the Orders. If not, this endpoint will return `403 Forbidden`.

*/
type GetOrdersCountByIssuer struct {
	Context *middleware.Context
	Handler GetOrdersCountByIssuerHandler
}

func (o *GetOrdersCountByIssuer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetOrdersCountByIssuerParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
