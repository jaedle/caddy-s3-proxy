//go:build integration

package integration_test

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3test "github.com/jaedle/caddy-s3-proxy/test/s3"
	"github.com/stretchr/testify/suite"
)

type integrationSuite struct {
	suite.Suite
	s3testClient s3test.S3Test
	caddyCmd     *exec.Cmd
}

func (s *integrationSuite) SetupSuite() {
	s.s3testClient = s3test.New()
	s.s3testClient.Clean(s.T())

	_, err := s.s3testClient.S3Client.CreateBucket(s.T().Context(), &s3.CreateBucketInput{
		Bucket: aws.String("example-bucket"),
	})
	s.Require().NoError(err)

	err = s3.NewBucketExistsWaiter(s.s3testClient.S3Client).Wait(s.T().Context(), &s3.HeadBucketInput{
		Bucket: aws.String("example-bucket"),
	}, 5*time.Second)
	s.Require().NoError(err)
}

func (s *integrationSuite) TearDown() {
	s.s3testClient.Clean(s.T())
	s.stopCaddy()
}

func (s *integrationSuite) TestWorks() {
	s.s3testClient.Put(s.T(), s3test.Obj("example-bucket", "some.html"))

	s.startCaddy()

	resp, err := http.Get("http://localhost:2015/some.html")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

const startupTimeout = 3 * time.Second

func (s *integrationSuite) startCaddy() {
	out := new(bytes.Buffer)
	command := exec.Command("build/caddy", "run")
	command.Dir = "../../example"
	command.Env = append(os.Environ(), []string{
		"AWS_ACCESS_KEY_ID=test",
		"AWS_SECRET_ACCESS_KEY=test",
		"AWS_REGION=us-east-1",
		"AWS_DEFAULT_REGION=us-east-1",
	}...)
	command.Stdout = out
	command.Stderr = out

	err := command.Start()
	s.Require().NoError(err)

	s.caddyCmd = command
	s.waitForCaddy()
}

func (s *integrationSuite) waitForCaddy() {
	ready := make(chan interface{}, 1)
	go func() {
		for {
			if resp, err := http.Get("http://localhost:2015/health"); err == nil && resp.StatusCode == http.StatusOK {
				ready <- nil
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
	select {
	case <-ready:
		return
	case <-time.After(startupTimeout):
		s.Fail("timed out waiting for caddy")
	}
}

func (s *integrationSuite) stopCaddy() {
	if s.caddyCmd == nil {
		return
	}

	s.Require().NoError(s.caddyCmd.Process.Kill())
	s.Require().NoError(s.caddyCmd.Wait())

	s.caddyCmd = nil
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(integrationSuite))
}
