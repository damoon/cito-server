package cito

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go"
)

func EnsureBucket(mc *minio.Client, bucket, location string) error {

	exists, err := mc.BucketExists(bucket)
	if err != nil {
		return fmt.Errorf("failed to access bucket %s: %s", bucket, err)
	}
	if exists {
		return nil
	}

	err = mc.MakeBucket(bucket, location)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %s", bucket, err)
	}
	log.Printf("bucket %s created\n", bucket)

	return nil
}

func bucketExists(ctx context.Context, mc *minio.Client, bucket string) (bool, error) {

	type bucketExists struct {
		found bool
		err   error
	}
	res := make(chan bucketExists)

	go func() {
		exists, err := mc.BucketExists(bucket)
		res <- bucketExists{found: exists, err: err}
	}()

	select {
	case res := <-res:
		return res.found, res.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func objectExists(ctx context.Context, mc *minio.Client, bucket, object string) (bool, error) {

	type exists struct {
		found bool
		err   error
	}
	res := make(chan exists)

	go func() {
		_, err := mc.StatObject(bucket, object, minio.StatObjectOptions{})
		if err != nil {
			errResponse := minio.ToErrorResponse(err)
			if errResponse.Code == "NoSuchKey" {
				res <- exists{found: false}
			}
			res <- exists{err: err}
		}
		res <- exists{found: true}
	}()

	select {
	case res := <-res:
		return res.found, res.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}
