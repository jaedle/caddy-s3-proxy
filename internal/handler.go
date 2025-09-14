package internal

import (
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

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
	_, _ = io.Copy(w, obj.Body)
	return nil
}
