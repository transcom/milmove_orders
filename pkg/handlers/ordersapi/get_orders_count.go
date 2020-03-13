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

// GetOrdersCountHandler returns Orders Count by issuer
type GetOrdersCountHandler struct {
	handlers.HandlerContext
}

// Handle (GetOrdersCountHandler) responds to GET /orders/{issuer}/count
func (h GetOrdersCountHandler) Handle(params ordersoperations.GetOrdersCountParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	logger := h.LoggerFromContext(ctx)

	clientCert := auth.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to access this API"))
	}

	count, err := models.FetchElectronicOrderCountByIssuer(h.DB(), params.Issuer)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	if !verifyOrdersReadAccess(models.Issuer(params.Issuer), clientCert) {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to read Orders from this issuer"))
	}

	count64 := int64(count)
	ordersCountPayload := &ordersmessages.OrdersCount{
		Issuer: ordersmessages.Issuer(params.Issuer),
		Count:  &count64,
	}

	return ordersoperations.NewGetOrdersCountOK().WithPayload(ordersCountPayload)
}
