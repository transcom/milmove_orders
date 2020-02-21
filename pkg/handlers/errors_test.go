package handlers

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/go-openapi/runtime/middleware"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/milmove_orders/pkg/models"
	"github.com/transcom/milmove_orders/pkg/testingsuite"
)

type fakeModel struct {
	ID   uuid.UUID
	Name string
}

type ErrorsSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestErrorsSuite(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	hs := &ErrorsSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *ErrorsSuite) TestResponseForErrorWhenASQLErrorIsEncountered() {
	var actual middleware.Responder
	var electronicOrder []*models.ElectronicOrder
	var noTableModel []*fakeModel

	// invalid column
	errInvalidColumn := suite.DB().Where("move_iid = $1", "123").All(&electronicOrder)
	// invalid arguments
	errInvalidArguments := suite.DB().Where("id in (?) and foo = ?", 1, 2, 3, "bar").All(electronicOrder)
	// invalid table
	errNoTable := suite.DB().Where("1=1").First(noTableModel)
	// invalid sql
	errInvalidQuery := suite.DB().Where("this should not compile").All(&electronicOrder)

	// slice to hold all errors and assert against
	errs := []error{errInvalidColumn, errNoTable, errInvalidArguments, errInvalidQuery}

	for _, err := range errs {
		actual = ResponseForError(suite.logger, err)
		res, ok := actual.(*ErrResponse)
		suite.True(ok)
		suite.Equal(500, res.Code)
		suite.Equal(SQLErrMessage, res.Err.Error())
	}

}

func (suite *ErrorsSuite) TestResponseForErrorNil() {

	var err error
	actual := ResponseForError(suite.logger, err)
	res, ok := actual.(*ErrResponse)
	suite.True(ok)
	suite.Equal(res.Code, 500)
	suite.Equal(res.Err.Error(), NilErrMessage)

}
