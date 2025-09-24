package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (s *handlerTestSuite) TestNotFound() {
	s.givenAnObject(s.obj(anHtmlFilename))

	res := s.do(httptest.NewRequest("GET", "/"+anotherHtmlFilename, nil))

	s.Equal(http.StatusNotFound, res.StatusCode)
	s.noBody(res)
}

func (s *handlerTestSuite) TestDeniesNotAllowedMethods() {
	notAllowed := []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}

	for _, m := range notAllowed {
		s.T().Run(m, func(t *testing.T) {
			res := s.do(httptest.NewRequest(m, "/", nil))
			assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
		})
	}
}
