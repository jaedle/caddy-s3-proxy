package plugin_test

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	plugin "github.com/jaedle/caddy-s3-proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const defaultValueAwsEndpoint = ""
const defaultValueUsePathStyle = false

func TestParsesConfiguration(t *testing.T) {
	actual := mustParse(t, `
aws_endpoint http://localhost:9000
bucket a-bucket-name
use_path_style true
`)

	assert.Equal(t, "http://localhost:9000", actual.AwsEndpoint)
	assert.Equal(t, "a-bucket-name", actual.Bucket)
	assert.Equal(t, true, actual.UsePathStyle)
}

func TestParsesMinimalConfiguration(t *testing.T) {
	actual := mustParse(t, `
bucket a-bucket-name
`)
	assert.Equal(t, defaultValueAwsEndpoint, actual.AwsEndpoint)
	assert.Equal(t, defaultValueUsePathStyle, actual.UsePathStyle)
}

func TestFailsOnMissingBucket(t *testing.T) {
	mustFailParsing(t, "", "missing required 'bucket' directive")
}

func TestFailsOnUnknownDirective(t *testing.T) {
	mustFailParsing(t, `
bucket a-bucket-name
unknown_directive value
`, "found unknown directive 'unknown_directive'")
}

func mustFailParsing(t *testing.T, in string, msg string) {
	cfg, err := plugin.Parse(caddyfile.NewTestDispenser(in))
	assert.Error(t, err)
	assert.Nil(t, cfg)

	assert.EqualError(t, err, msg)
}

func mustParse(t *testing.T, input string) *plugin.Configuration {
	cfg, err := plugin.Parse(caddyfile.NewTestDispenser(input))
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	return cfg
}
