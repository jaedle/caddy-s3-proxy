package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/jaedle/caddy-s3-proxy/internal"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int

	handler caddyhttp.MiddlewareHandler
}

func (s *handlerTestSuite) SetupTest() {
	s.handler = internal.New()
}

func (s *handlerTestSuite) TestHandlesRequest() {
	res := s.do(httptest.NewRequest("GET", "/", nil))

	s.Equal(http.StatusNotFound, res.StatusCode)

}

func (s *handlerTestSuite) do(req *http.Request) *http.Response {
	res := httptest.NewRecorder()
	s.Require().NoError(s.handler.ServeHTTP(res, req, nil))
	return res.Result()
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
