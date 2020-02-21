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

// IndexOrdersForMemberHandler returns a list of Orders matching the provided search parameters
type IndexOrdersForMemberHandler struct {
	handlers.HandlerContext
}

// Handle (IndexOrdersForMemberHandler) responds to GET /edipis/{edipi}/orders
func (h IndexOrdersForMemberHandler) Handle(params ordersoperations.IndexOrdersForMemberParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	logger := h.LoggerFromContext(ctx)

	clientCert := auth.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to access this API"))
	}
	allowedIssuers := clientCert.GetAllowedOrdersIssuersRead()
	if len(allowedIssuers) == 0 {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to read any Orders"))
	}

	orders, err := models.FetchElectronicOrdersByEdipiAndIssuers(h.DB(), params.Edipi, allowedIssuers)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewIndexOrdersForMemberOK().WithPayload([]*ordersmessages.Orders{})
	} else if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	ordersPayloads := make([]*ordersmessages.Orders, len(orders))
	for i, o := range orders {
		payload, err := payloadForElectronicOrderModel(o)
		if err != nil {
			return handlers.ResponseForError(logger, err)
		}
		ordersPayloads[i] = payload
	}

	return ordersoperations.NewIndexOrdersForMemberOK().WithPayload(ordersPayloads)
}
