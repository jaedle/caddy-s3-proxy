package internal_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func (s *handlerTestSuite) TestServesContent() {
	s.givenAnObject(
		s.obj(anHtmlFilename).Content(anHtmlFileContent).ContentType("text/html"),
	)

	res := s.do(httptest.NewRequest("GET", "/"+anHtmlFilename, nil))

	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal("text/html", res.Header.Get("Content-Type"))
	s.Equal(fmt.Sprintf("%d", len(anHtmlFileContent)), res.Header.Get("Content-Length"))
	s.Equal([]byte(anHtmlFileContent), body(res, s))
}
