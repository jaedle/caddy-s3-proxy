package plugin

import (
	"fmt"
	"strconv"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

type Configuration struct {
	AwsEndpoint  string
	Bucket       string
	UsePathStyle bool
}

func Parse(d *caddyfile.Dispenser) (*Configuration, error) {
	var c Configuration
	for d.Next() {
		switch d.Val() {
		case "aws_endpoint":
			strArg(d, &c.AwsEndpoint)
		case "bucket":
			strArg(d, &c.Bucket)
		case "use_path_style":
			boolArg(d, &c.UsePathStyle)
		default:
			return nil, fmt.Errorf("found unknown directive '%s'", d.Val())
		}
	}

	if err := c.verify(); err != nil {
		return nil, err
	} else {
		return &c, nil
	}
}

func strArg(d *caddyfile.Dispenser, s *string) {
	if !d.NextArg() {
		return
	}

	*s = d.Val()
}

func boolArg(d *caddyfile.Dispenser, b *bool) {
	var s string
	strArg(d, &s)
	actual, _ := strconv.ParseBool(s)
	*b = actual
}

func (c *Configuration) verify() error {
	if len(c.Bucket) == 0 {
		return fmt.Errorf("missing required 'bucket' directive")
	}
	return nil
}
