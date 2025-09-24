package s3test

import "github.com/aws/aws-sdk-go-v2/aws"

const anyContentType = "application/octet-stream"

type Object struct {
	bucket      string
	key         string
	content     *string
	contentType *string
}

func Obj(bucket string, key string) Object {
	return Object{
		bucket:      bucket,
		key:         key,
		contentType: aws.String(anyContentType),
	}
}

func (o Object) Content(s string) Object {
	o.content = &s
	return o
}

func (o Object) ContentType(ct string) Object {
	o.contentType = &ct
	return o
}
