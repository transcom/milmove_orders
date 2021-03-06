package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/transcom/milmove_orders/pkg/auth"
	"github.com/transcom/milmove_orders/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/milmove_orders/pkg/handlers"
	"github.com/transcom/milmove_orders/pkg/models"
)

// GetOrdersByIssuerAndOrdersNumHandler returns Orders for a specific issuer by ordersNum
type GetOrdersByIssuerAndOrdersNumHandler struct {
	handlers.HandlerContext
}

// Handle (GetOrdersByIssuerAndOrdersNumHandler) responds to GET /issuers/{issuer}/orders/{ordersNum}
func (h GetOrdersByIssuerAndOrdersNumHandler) Handle(params ordersoperations.GetOrdersByIssuerAndOrdersNumParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	logger := h.LoggerFromContext(ctx)

	clientCert := auth.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to access this API"))
	}
	if !verifyOrdersReadAccess(models.Issuer(params.Issuer), clientCert) {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to read orders from this issuer"))
	}

	orders, err := models.FetchElectronicOrderByIssuerAndOrdersNum(h.DB(), params.Issuer, params.OrdersNum)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	return ordersoperations.NewGetOrdersByIssuerAndOrdersNumOK().WithPayload(ordersPayload)
}
