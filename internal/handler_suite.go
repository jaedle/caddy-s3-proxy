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
	if !h.isAllowedMethod(r) {
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
	if obj == nil {
		return nil
	}
	defer func() { _ = obj.Body.Close() }()

	if cacheHit(r, obj) {
		h.notModified(w, obj)
	} else {
		h.Ok(w, obj)
	}

	return nil
}

func cacheHit(r *http.Request, obj *s3.GetObjectOutput) bool {
	return r.Header.Get(headerIfNoneMatch) == aws.ToString(obj.ETag)
}

func (h *handler) notModified(w http.ResponseWriter, obj *s3.GetObjectOutput) {
	w.Header().Set(headerContentType, aws.ToString(obj.ContentType))
	w.Header().Set(headerEtag, aws.ToString(obj.ETag))

	w.WriteHeader(http.StatusNotModified)
}

func (h *handler) Ok(w http.ResponseWriter, obj *s3.GetObjectOutput) {
	w.Header().Set(headerContentLength, strconv.FormatInt(aws.ToInt64(obj.ContentLength), 10))
	w.Header().Set(headerContentType, aws.ToString(obj.ContentType))
	w.Header().Set(headerEtag, aws.ToString(obj.ETag))
	_, _ = io.Copy(w, obj.Body)
}

func (h *handler) isAllowedMethod(r *http.Request) bool {
	return r.Method == http.MethodGet
}
