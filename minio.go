package cito

import (
	"fmt"

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

	return nil
}
