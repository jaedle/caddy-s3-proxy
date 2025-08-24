package s3test

type Object struct {
	bucket string
	key    string
}

func Obj(bucket string, key string) Object {
	return Object{
		bucket: bucket,
		key:    key,
	}
}
