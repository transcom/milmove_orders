package ordersapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/milmove_orders/pkg/auth"
	"github.com/transcom/milmove_orders/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/milmove_orders/pkg/handlers"
	"github.com/transcom/milmove_orders/pkg/models"
)

// PostRevisionToOrdersHandler adds a Revision to Orders by uuid
type PostRevisionToOrdersHandler struct {
	handlers.HandlerContext
}

// Handle (params ordersoperations.PostRevisionToOrdersParams) responds to POST /orders/{uuid}
func (h PostRevisionToOrdersHandler) Handle(params ordersoperations.PostRevisionToOrdersParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	logger := h.LoggerFromContext(ctx)

	clientCert := auth.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrWriteForbidden, "Not permitted to access this API"))
	}

	id, err := uuid.FromString(params.UUID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	orders, err := models.FetchElectronicOrderByID(h.DB(), id)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	if !verifyOrdersWriteAccess(orders.Issuer, clientCert) {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrWriteForbidden, "Not permitted to write Orders from this issuer"))
	}

	for _, r := range orders.Revisions {
		// SeqNum collision
		if r.SeqNum == int(*params.Revision.SeqNum) {
			return handlers.ResponseForError(
				logger,
				errors.WithMessage(
					models.ErrWriteConflict,
					fmt.Sprintf("Cannot POST Revision with SeqNum %d to Orders %s: a Revision with that SeqNum already exists in those Orders", r.SeqNum, params.UUID)))
		}
	}

	newRevision := toElectronicOrdersRevision(orders, params.Revision)
	verrs, err := models.CreateElectronicOrdersRevision(ctx, h.DB(), newRevision)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	orders.Revisions = append(orders.Revisions, *newRevision)

	orderPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return ordersoperations.NewPostRevisionToOrdersCreated().WithPayload(orderPayload)
}
