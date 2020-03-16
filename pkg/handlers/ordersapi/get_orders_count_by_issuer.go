package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/transcom/milmove_orders/pkg/auth"
	"github.com/transcom/milmove_orders/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/milmove_orders/pkg/gen/ordersmessages"
	"github.com/transcom/milmove_orders/pkg/handlers"
	"github.com/transcom/milmove_orders/pkg/models"
)

// GetOrdersCountByIssuerHandler returns Orders Count by Issuer
type GetOrdersCountByIssuerHandler struct {
	handlers.HandlerContext
}

// Handle (GetOrdersCountByIssuerHandler) responds to GET /orders/{issuer}/count
func (h GetOrdersCountByIssuerHandler) Handle(params ordersoperations.GetOrdersCountByIssuerParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	logger := h.LoggerFromContext(ctx)

	clientCert := auth.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to access this API"))
	}

	var startDateTime, endDateTime *string
	if params.StartDateTime != nil {
		dt := params.StartDateTime.String()
		startDateTime = &dt
	}
	if params.EndDateTime != nil {
		dt := params.EndDateTime.String()
		endDateTime = &dt
	}
	ordersCountByIssuer, err := models.FetchElectronicOrderCountByIssuer(h.DB(), params.Issuer, startDateTime, endDateTime)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	if !verifyOrdersReadAccess(models.Issuer(params.Issuer), clientCert) {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to read Orders from this issuer"))
	}

	ordersCountByIssuerPayload := &ordersmessages.OrdersCountByIssuer{
		Issuer:        ordersmessages.Issuer(params.Issuer),
		Count:         ordersCountByIssuer,
		StartDateTime: params.StartDateTime,
		EndDateTime:   params.EndDateTime,
	}

	return ordersoperations.NewGetOrdersCountByIssuerOK().WithPayload(ordersCountByIssuerPayload)
}
