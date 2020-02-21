package auth

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/milmove_orders/pkg/testingsuite"
)

const (
	// OrdersTestHost
	OrdersTestHost string = "orders.example.com"
)

// ApplicationTestServername is a collection of the test servernames
func ApplicationTestServername() ApplicationServername {
	appnames := ApplicationServername{
		OrdersServername: OrdersTestHost,
	}
	return appnames
}

type authSuite struct {
	testingsuite.BaseTestSuite
	logger Logger
}

func TestAuthSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &authSuite{logger: logger}
	suite.Run(t, hs)
}
