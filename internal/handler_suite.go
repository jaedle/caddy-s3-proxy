package internal

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

const headerContentLength = "Content-Length"
const headerContentType = "Content-Type"
const headerIfNoneMatch = "If-None-Match"
const headerEtag = "Etag"
const headerCacheControl = "Cache-Control"

func New(c Config) caddyhttp.MiddlewareHandler {
	return &handler{
		bucket:   c.Bucket,
		S3Client: c.S3Client,
	}
}

type Config struct {
	Bucket   string
	S3Client *s3.Client
}

type handler struct {
	bucket   string
	S3Client *s3.Client
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request, _ caddyhttp.Handler) error {
	if !isAllowedMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	obj, err := h.S3Client.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: aws.String(h.bucket),
		Key:    aws.String(strings.TrimPrefix(r.URL.Path, "/")),
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	defer func() { _ = obj.Body.Close() }()

	if isCacheHit(r, obj) {
		notModified(w, obj)
	} else {
		ok(w, obj)
	}

	return nil
}

func isAllowedMethod(r *http.Request) bool {
	return r.Method == http.MethodGet
}

func isCacheHit(r *http.Request, obj *s3.GetObjectOutput) bool {
	return r.Header.Get(headerIfNoneMatch) == aws.ToString(obj.ETag)
}

func notModified(w http.ResponseWriter, obj *s3.GetObjectOutput) {
	setCommonHeaders(w, obj)
	w.WriteHeader(http.StatusNotModified)
}

func ok(w http.ResponseWriter, obj *s3.GetObjectOutput) {
	setCommonHeaders(w, obj)
	w.Header().Set(headerContentLength, strconv.FormatInt(aws.ToInt64(obj.ContentLength), 10))

	_, _ = io.Copy(w, obj.Body)
}

func setCommonHeaders(w http.ResponseWriter, obj *s3.GetObjectOutput) {
	w.Header().Set(headerContentType, aws.ToString(obj.ContentType))
	w.Header().Set(headerEtag, aws.ToString(obj.ETag))

	if cacheControl := aws.ToString(obj.CacheControl); len(cacheControl) > 0 {
		w.Header().Set(headerCacheControl, cacheControl)
	}
}
