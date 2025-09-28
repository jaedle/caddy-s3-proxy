package internal_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func (s *handlerTestSuite) TestServesContent() {
	s.givenAnObject(
		s.obj(anHtmlFilename).Content(anHtmlFileContent).ContentType(contentTypeTextHtml),
	)

	res := s.do(httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil))

	s.Equal(http.StatusOK, res.StatusCode)
	s.expectHeaderEqual(res, contentTypeTextHtml, headerContentType)
	s.expectHeaderEqual(res, fmt.Sprintf("%d", len(anHtmlFileContent)), headerContentLength)
	s.expectHeaderPresent(res, headerEtag)
	s.expectHeaderPresent(res, headerLastModified)
	s.Equal([]byte(anHtmlFileContent), s.body(res))
}

func (s *handlerTestSuite) TestServesLastModifiedAsTimestamp() {
	s.givenAnObject(
		s.obj(anHtmlFilename),
	)

	res := s.do(httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil))

	s.Equal(http.StatusOK, res.StatusCode)
	lastModified := res.Header.Get(headerLastModified)
	_, err := time.Parse(http.TimeFormat, lastModified)
	s.NoError(err, "must return a valid HTTP timestamp")
}

func (s *handlerTestSuite) TestCachesByEtag() {
	s.givenAnObject(
		s.obj(anHtmlFilename).ContentType(contentTypeTextHtml).Content(anHtmlFileContent),
	)

	res := s.do(httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil))
	etag := res.Header.Get(headerEtag)

	req := httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil)
	req.Header.Set(headerIfNoneMatch, etag)
	res = s.do(req)

	s.Equal(http.StatusNotModified, res.StatusCode)
	s.Equal(etag, res.Header.Get(headerEtag))
	s.expectHeaderEqual(res, contentTypeTextHtml, headerContentType)
	s.expectHeaderNotPresent(res, headerContentLength)
	s.expectHeaderPresent(res, headerLastModified)
	s.noBody(res)
}

func (s *handlerTestSuite) TestServesFileOnEtagMismatch() {
	s.givenAnObject(
		s.obj(anHtmlFilename).Content(anHtmlFileContent),
	)

	req := httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil)
	req.Header.Set(headerIfNoneMatch, "anything-that-does-not-match")
	res := s.do(req)

	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal([]byte(anHtmlFileContent), s.body(res))
}

func (s *handlerTestSuite) TestServesCacheControlWhenPresent() {
	s.givenAnObject(
		s.obj(anHtmlFilename).CacheControl(aCacheDirective),
	)

	res := s.do(httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil))

	s.Equal(http.StatusOK, res.StatusCode)
	s.expectHeaderEqual(res, aCacheDirective, headerCacheControl)
}

func (s *handlerTestSuite) TestServesCacheControlWhenPresentOnNonModified() {
	s.givenAnObject(
		s.obj(anHtmlFilename).CacheControl(aCacheDirective),
	)
	res := s.do(httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil))
	etag := res.Header.Get(headerEtag)

	req := httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil)
	req.Header.Set(headerIfNoneMatch, etag)
	res = s.do(req)

	s.Equal(http.StatusNotModified, res.StatusCode)
	s.expectHeaderEqual(res, aCacheDirective, headerCacheControl)
}

func (s *handlerTestSuite) TestDoesNotServeCacheControlWhenNotPresent() {
	s.givenAnObject(
		s.obj(anHtmlFilename).ContentType(contentTypeTextHtml),
	)

	res := s.do(httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil))

	s.Equal(http.StatusOK, res.StatusCode)
	s.expectHeaderNotPresent(res, headerCacheControl)
}

func (s *handlerTestSuite) TestCachesByLastModified() {
	s.givenAnObject(
		s.obj(anHtmlFilename).ContentType(contentTypeTextHtml).Content(anHtmlFileContent),
	)

	res := s.do(httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil))
	lastModified := res.Header.Get(headerLastModified)

	req := httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil)
	req.Header.Set(headerModifiedSince, lastModified)
	res = s.do(req)

	s.Equal(http.StatusNotModified, res.StatusCode)
	s.expectHeaderEqual(res, contentTypeTextHtml, headerContentType)
	s.expectHeaderNotPresent(res, headerContentLength)
	s.expectHeaderPresent(res, headerEtag)
	s.expectHeaderEqual(res, lastModified, headerLastModified)
	s.noBody(res)
}

func (s *handlerTestSuite) TestServesFileOnLastModifiedMismatch() {
	s.givenAnObject(
		s.obj(anHtmlFilename).ContentType(contentTypeTextHtml).Content(anHtmlFileContent),
	)

	req := httptest.NewRequest(http.MethodGet, "/"+anHtmlFilename, nil)
	req.Header.Set(headerModifiedSince, time.Now().Add(-1*time.Hour).Format(http.TimeFormat))

	res := s.do(req)

	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal([]byte(anHtmlFileContent), s.body(res))
}
