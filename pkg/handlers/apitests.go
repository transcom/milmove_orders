package handlers

import (
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"runtime/debug"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/milmove_orders/pkg/testingsuite"
)

// BaseHandlerTestSuite abstracts the common methods needed for handler tests
type BaseHandlerTestSuite struct {
	testingsuite.PopTestSuite
	logger       Logger
	filesToClose []*runtime.File
}

// NewBaseHandlerTestSuite returns a new BaseHandlerTestSuite
func NewBaseHandlerTestSuite(logger Logger, packageName testingsuite.PackageName) BaseHandlerTestSuite {
	return BaseHandlerTestSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(packageName),
		logger:       logger,
	}
}

// TestLogger returns the logger to use in the suite
func (suite *BaseHandlerTestSuite) TestLogger() Logger {
	return suite.logger
}

// TestFilesToClose returns the list of files needed to close at the end of tests
func (suite *BaseHandlerTestSuite) TestFilesToClose() []*runtime.File {
	return suite.filesToClose
}

// SetTestFilesToClose sets the list of files needed to close at the end of tests
func (suite *BaseHandlerTestSuite) SetTestFilesToClose(filesToClose []*runtime.File) {
	suite.filesToClose = filesToClose
}

// CloseFile adds a single file to close at the end of tests to the list of files
func (suite *BaseHandlerTestSuite) CloseFile(file *runtime.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

// IsNotErrResponse enforces handler does not return an error response
func (suite *BaseHandlerTestSuite) IsNotErrResponse(response middleware.Responder) {
	r, ok := response.(*ErrResponse)
	if ok {
		suite.logger.Error("Received an unexpected error response from handler: ", zap.Error(r.Err))
		// Formally lodge a complaint
		suite.IsType(&ErrResponse{}, response)
	}
}

// CheckErrorResponse verifies error response is what is expected
func (suite *BaseHandlerTestSuite) CheckErrorResponse(resp middleware.Responder, code int, name string) {
	errResponse, ok := resp.(*ErrResponse)
	if !ok || errResponse.Code != code {
		suite.T().Fatalf("Expected %s, Response: %v, Code: %v", name, resp, code)
		debug.PrintStack()
	}
}

// CheckNotErrorResponse verifies there is no error response
func (suite *BaseHandlerTestSuite) CheckNotErrorResponse(resp middleware.Responder) {
	errResponse, ok := resp.(*ErrResponse)
	if ok {
		suite.NoError(errResponse.Err)
		suite.FailNowf("Received error response", "Code: %v", errResponse.Code)
	}
}

// CheckResponseBadRequest looks at BadRequest errors
func (suite *BaseHandlerTestSuite) CheckResponseBadRequest(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusBadRequest, "BadRequest")
}

// CheckResponseUnauthorized looks at Unauthorized errors
func (suite *BaseHandlerTestSuite) CheckResponseUnauthorized(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusUnauthorized, "Unauthorized")
}

// CheckResponseForbidden looks at Forbidden errors
func (suite *BaseHandlerTestSuite) CheckResponseForbidden(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusForbidden, "Forbidden")
}

// CheckResponseNotFound looks at NotFound errors
func (suite *BaseHandlerTestSuite) CheckResponseNotFound(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusNotFound, "NotFound")
}

// CheckResponseInternalServerError looks at InternalServerError errors
func (suite *BaseHandlerTestSuite) CheckResponseInternalServerError(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusInternalServerError, "InternalServerError")
}

// CheckResponseTeapot enforces that response come from a Teapot
func (suite *BaseHandlerTestSuite) CheckResponseTeapot(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusTeapot, "Teapot")
}

// Fixture allows us to include a fixture like a PDF in the test
func (suite *BaseHandlerTestSuite) Fixture(name string) *runtime.File {
	fixtureDir := "testdatagen/testdata"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Error(err)
	}

	fixturePath := path.Join(cwd, "..", "..", fixtureDir, name)

	// #nosec never comes from user input
	file, err := os.Open(fixturePath)
	if err != nil {
		suite.logger.Fatal("Error opening fixture file", zap.Error(err))
	}

	info, err := file.Stat()
	if err != nil {
		suite.logger.Fatal("Error accessing fixture stats", zap.Error(err))
	}

	header := multipart.FileHeader{
		Filename: info.Name(),
		Size:     info.Size(),
	}

	returnFile := &runtime.File{
		Header: &header,
		Data:   file,
	}
	suite.CloseFile(returnFile)

	return returnFile
}
