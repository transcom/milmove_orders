package handlers

import (
	"context"
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/milmove_orders/pkg/auth"
	"github.com/transcom/milmove_orders/pkg/logging"
	"github.com/transcom/mymove/pkg/iws"
)

// HandlerContext provides access to all the contextual references needed by individual handlers
//go:generate mockery -name HandlerContext
type HandlerContext interface {
	DB() *pop.Connection
	LoggerFromContext(ctx context.Context) Logger
	LoggerFromRequest(r *http.Request) Logger
	IWSPersonLookup() iws.PersonLookup
	SetIWSPersonLookup(rbs iws.PersonLookup)
	SetAppNames(appNames auth.ApplicationServername)
	AppNames() auth.ApplicationServername
	SetTraceID(traceID uuid.UUID)
	GetTraceID() uuid.UUID
}

// FeatureFlag struct for feature flags
type FeatureFlag struct {
	Name   string
	Active bool
}

// A single handlerContext is passed to each handler
type handlerContext struct {
	db              *pop.Connection
	logger          Logger
	iwsPersonLookup iws.PersonLookup
	appNames        auth.ApplicationServername
	traceID         uuid.UUID
}

// NewHandlerContext returns a new handlerContext with its required private fields set.
func NewHandlerContext(db *pop.Connection, logger Logger) HandlerContext {
	return &handlerContext{
		db:     db,
		logger: logger,
	}
}

func (hctx *handlerContext) LoggerFromRequest(r *http.Request) Logger {
	return hctx.LoggerFromContext(r.Context())
}

func (hctx *handlerContext) LoggerFromContext(ctx context.Context) Logger {
	if logger, ok := logging.FromContext(ctx).(Logger); ok {
		return logger
	}
	return hctx.logger
}

// DB returns a POP db connection for the context
func (hctx *handlerContext) DB() *pop.Connection {
	return hctx.db
}

// AppNames returns a struct of all the app names for the current environment
func (hctx *handlerContext) AppNames() auth.ApplicationServername {
	return hctx.appNames
}

// SetAppNames is a simple setter for private field
func (hctx *handlerContext) SetAppNames(appNames auth.ApplicationServername) {
	hctx.appNames = appNames
}

func (hctx *handlerContext) IWSPersonLookup() iws.PersonLookup {
	return hctx.iwsPersonLookup
}

func (hctx *handlerContext) SetIWSPersonLookup(rbs iws.PersonLookup) {
	hctx.iwsPersonLookup = rbs
}

func (hctx *handlerContext) SetTraceID(traceID uuid.UUID) {
	hctx.traceID = traceID
}

func (hctx *handlerContext) GetTraceID() uuid.UUID {
	return hctx.traceID
}
