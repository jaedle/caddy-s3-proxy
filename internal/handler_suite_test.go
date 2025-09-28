package internal_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/jaedle/caddy-s3-proxy/internal"
	s3test "github.com/jaedle/caddy-s3-proxy/test/s3"
	"github.com/stretchr/testify/suite"
)

const anHtmlFilename = "some.html"
const anotherHtmlFilename = "another.html"

const anHtmlFileContent = "<html></html>"

const contentTypeTextHtml = "text/html"

const headerContentType = "Content-Type"
const headerContentLength = "Content-Length"
const headerIfNoneMatch = "If-None-Match"
const headerEtag = "Etag"
const headerCacheControl = "Cache-Control"
const headerLastModified = "Last-Modified"

const aCacheDirective = "public, max-age=3600"

type handlerTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int

	handler      caddyhttp.MiddlewareHandler
	testS3Client s3test.S3Test
	bucket       string
}

func (s *handlerTestSuite) SetupTest() {
	s.testS3Client = s3test.New()
	s.bucket = s.testS3Client.ABucket(s.T())

	s.handler = internal.New(internal.Config{
		S3Client: s.testS3Client.S3Client,
		Bucket:   s.bucket,
	})
}

func (s *handlerTestSuite) TearDownTest() {
	s.testS3Client.Clean(s.T())
}

func (s *handlerTestSuite) givenAnObject(obj s3test.Object) {
	s.testS3Client.Put(s.T(), obj)
}

func (s *handlerTestSuite) obj(key string) s3test.Object {
	return s3test.Obj(s.bucket, key)
}

func (s *handlerTestSuite) body(res *http.Response) []byte {
	all, err := io.ReadAll(res.Body)
	s.Require().NoError(err)
	return all
}

func (s *handlerTestSuite) do(req *http.Request) *http.Response {
	res := httptest.NewRecorder()
	s.Require().NoError(s.handler.ServeHTTP(res, req, nil))
	return res.Result()
}

func (s *handlerTestSuite) expectHeaderEqual(res *http.Response, expected string, name string) bool {
	return s.Equal(expected, res.Header.Get(name))
}

func (s *handlerTestSuite) expectHeaderPresent(res *http.Response, name string) {
	s.NotEmpty(res.Header.Get(name))
}

func (s *handlerTestSuite) expectHeaderNotPresent(res *http.Response, name string) {
	s.Empty(res.Header.Get(name))
}

func (s *handlerTestSuite) noBody(res *http.Response) {
	s.Empty(s.body(res))
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
