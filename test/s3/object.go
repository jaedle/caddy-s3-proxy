package s3test

type Object struct {
	bucket  string
	key     string
	content *string
}

func Obj(bucket string, key string) Object {
	return Object{
		bucket: bucket,
		key:    key,
	}
}

func (o Object) Content(s string) Object {
	o.content = &s
	return o
}
