package ordersapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/milmove_orders/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/milmove_orders/pkg/handlers"
	"github.com/transcom/milmove_orders/pkg/models"
	"github.com/transcom/milmove_orders/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrdersCountByIssuerSuccess() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/count", string(order.Issuer)), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:          true,
		AllowAirForceOrdersRead: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersCountParams{
		HTTPRequest: req,
		Issuer:      string(order.Issuer),
	}

	handler := GetOrdersCountHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersCountOK{}, response)
	okResponse, ok := response.(*ordersoperations.GetOrdersCountOK)
	if !ok {
		return
	}
	suite.Equal(string(order.Issuer), string(okResponse.Payload.Issuer))
	suite.Equal(int64(1), *okResponse.Payload.Count)
}

func (suite *HandlerSuite) TestGetOrdersCountByIssuerNoApiPerm() {
	req := httptest.NewRequest("GET", "/orders/v1/issuers/air-force/count", nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersCountParams{
		HTTPRequest: req,
		Issuer:      string(models.IssuerAirForce),
	}

	handler := GetOrdersCountHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse, ok := response.(*handlers.ErrResponse)
	if !ok {
		return
	}
	suite.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestGetOrdersCountByIssuerReadPerms() {
	testCases := map[string]struct {
		cert   *models.ClientCert
		issuer models.Issuer
	}{
		"Army": {
			makeAllPowerfulClientCert(),
			models.IssuerArmy,
		},
		"Navy": {
			makeAllPowerfulClientCert(),
			models.IssuerNavy,
		},
		"MarineCorps": {
			makeAllPowerfulClientCert(),
			models.IssuerMarineCorps,
		},
		"CoastGuard": {
			makeAllPowerfulClientCert(),
			models.IssuerCoastGuard,
		},
		"AirForce": {
			makeAllPowerfulClientCert(),
			models.IssuerAirForce,
		},
	}
	testCases["Army"].cert.AllowArmyOrdersRead = false
	testCases["Navy"].cert.AllowNavyOrdersRead = false
	testCases["MarineCorps"].cert.AllowMarineCorpsOrdersRead = false
	testCases["CoastGuard"].cert.AllowCoastGuardOrdersRead = false
	testCases["AirForce"].cert.AllowAirForceOrdersRead = false

	for name, testCase := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			assertions := testdatagen.Assertions{
				ElectronicOrder: models.ElectronicOrder{
					Issuer: testCase.issuer,
				},
			}
			order := testdatagen.MakeElectronicOrder(suite.DB(), assertions)
			req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/count", string(order.Issuer)), nil)
			req = suite.AuthenticateClientCertRequest(req, testCase.cert)

			params := ordersoperations.GetOrdersCountParams{
				HTTPRequest: req,
				Issuer:      string(order.Issuer),
			}

			handler := GetOrdersCountHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
			response := handler.Handle(params)

			suite.IsType(&handlers.ErrResponse{}, response)
			errResponse, ok := response.(*handlers.ErrResponse)
			if !ok {
				return
			}
			suite.Equal(http.StatusForbidden, errResponse.Code)
		})
	}
}
