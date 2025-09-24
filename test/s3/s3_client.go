package s3test

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const region = "us-east-1"
const localEndpoint = "http://localhost:4577"

type S3Test struct {
	S3Client *s3.Client
}

const testBucketPrefix = "test-bucket-"

func New() S3Test {
	s3Client := s3.NewFromConfig(aws.Config{
		Region:       region,
		Credentials:  credentials.NewStaticCredentialsProvider("test", "test", ""),
		BaseEndpoint: aws.String(localEndpoint),
	}, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return S3Test{
		S3Client: s3Client,
	}
}

func (s *S3Test) Put(t *testing.T, obj Object) {
	_, err := s.S3Client.PutObject(t.Context(), &s3.PutObjectInput{
		Bucket:       aws.String(obj.bucket),
		Key:          aws.String(obj.key),
		Body:         toReader(obj.content),
		ContentType:  obj.contentType,
		CacheControl: obj.cacheControl,
	})
	require.NoError(t, err)
}

func toReader(s *string) *strings.Reader {
	if s == nil {
		return strings.NewReader("")
	}

	return strings.NewReader(*s)
}

func (s *S3Test) ABucket(t *testing.T) string {
	var name = testBucketPrefix + uuid.NewString()
	_, err := s.S3Client.CreateBucket(t.Context(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})

	require.NoError(t, err)
	return name
}

func (s *S3Test) Clean(t *testing.T) {
	paginator := s3.NewListBucketsPaginator(s.S3Client, &s3.ListBucketsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(t.Context())
		require.NoError(t, err)
		for _, bucket := range page.Buckets {
			s.deleteBucket(t, aws.ToString(bucket.Name))
		}
	}
}

func (s *S3Test) deleteBucket(t *testing.T, name string) {
	s.emptyBucket(t, name)

	_, err := s.S3Client.DeleteBucket(t.Context(), &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	require.NoError(t, err)
}

func (s *S3Test) emptyBucket(t *testing.T, name string) {
	paginator := s3.NewListObjectVersionsPaginator(s.S3Client, &s3.ListObjectVersionsInput{
		Bucket: aws.String(name),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(t.Context())
		require.NoError(t, err)

		if isEmpty(page) {
			return
		}

		_, err = s.S3Client.DeleteObjects(t.Context(), &s3.DeleteObjectsInput{
			Bucket: aws.String(name),
			Delete: &types.Delete{Objects: toObjectIdentifiers(page)},
		})
		require.NoError(t, err)
	}
}

func toObjectIdentifiers(page *s3.ListObjectVersionsOutput) []types.ObjectIdentifier {
	var result []types.ObjectIdentifier
	for _, obj := range page.Versions {
		result = append(result, types.ObjectIdentifier{
			Key:       obj.Key,
			VersionId: obj.VersionId,
		})
	}
	return result
}

func isEmpty(page *s3.ListObjectVersionsOutput) bool {
	return len(page.Versions) == 0
}
