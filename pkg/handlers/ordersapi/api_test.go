package ordersapi

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/milmove_orders/pkg/auth"
	"github.com/transcom/milmove_orders/pkg/handlers"
	"github.com/transcom/milmove_orders/pkg/models"
	"github.com/transcom/milmove_orders/pkg/testingsuite"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	handlers.BaseHandlerTestSuite
}

// AuthenticateClientCertRequest authenticates mutual TLS auth API users with the provided ClientCert object
func (suite *HandlerSuite) AuthenticateClientCertRequest(req *http.Request, cert *models.ClientCert) *http.Request {
	ctx := auth.SetClientCertInRequestContext(req, cert)
	return req.WithContext(ctx)
}

// SetupTest sets up the test suite by preparing the DB
func (suite *HandlerSuite) SetupTest() {
	errTruncateAll := suite.DB().TruncateAll()
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database: %#v", errTruncateAll)
	}
}

// AfterTest completes tests by trying to close open files
func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.TestFilesToClose() {
		err := file.Data.Close()
		suite.NoError(err)
	}
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &HandlerSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(logger, testingsuite.CurrentPackage()),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func makeAllPowerfulClientCert() *models.ClientCert {
	return &models.ClientCert{
		AllowAirForceOrdersRead:     true,
		AllowAirForceOrdersWrite:    true,
		AllowArmyOrdersRead:         true,
		AllowArmyOrdersWrite:        true,
		AllowCoastGuardOrdersRead:   true,
		AllowCoastGuardOrdersWrite:  true,
		AllowMarineCorpsOrdersRead:  true,
		AllowMarineCorpsOrdersWrite: true,
		AllowNavyOrdersRead:         true,
		AllowNavyOrdersWrite:        true,
		AllowOrdersAPI:              true,
	}
}
