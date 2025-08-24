package plugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"

	"github.com/jaedle/caddy-s3-proxy/internal"
)

var (
	_ caddy.Provisioner           = (*S3Proxy)(nil)
	_ caddyfile.Unmarshaler       = (*S3Proxy)(nil)
	_ caddyhttp.MiddlewareHandler = (*S3Proxy)(nil)
)

func init() {
	caddy.RegisterModule(&S3Proxy{})
	httpcaddyfile.RegisterHandlerDirective("s3_proxy", parseCaddyfile)
}

func (*S3Proxy) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.s3_proxy",
		New: func() caddy.Module { return new(S3Proxy) },
	}
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	proxy := new(S3Proxy)
	if err := proxy.UnmarshalCaddyfile(h.Dispenser); err != nil {
		return nil, err
	} else {
		return proxy, nil
	}
}

func (m *S3Proxy) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if cfg, err := Parse(d); err != nil {
		return err
	} else {
		m.Config = cfg
		return nil
	}
}

type S3Proxy struct {
	Config *Configuration `json:"config,omitempty"`
}

func (m *S3Proxy) Provision(ctx caddy.Context) error {
	return nil
}

func (m *S3Proxy) Validate() error {
	if m.Config == nil || m.Config.Bucket == "" {
		return fmt.Errorf("bucket is required")
	}
	return nil
}

func (m *S3Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), m.endpointOptions())
	if err != nil {
		return err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if m.Config != nil {
			o.UsePathStyle = m.Config.UsePathStyle
		}
	})

	return internal.New(internal.Config{
		Bucket:   m.Config.Bucket,
		S3Client: s3Client,
	}).ServeHTTP(w, r, next)
}

func (m *S3Proxy) endpointOptions() config.LoadOptionsFunc {
	if len(m.Config.AwsEndpoint) == 0 {
		return nil
	}

	return config.WithBaseEndpoint(m.Config.AwsEndpoint)
}
