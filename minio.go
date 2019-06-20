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
	log.Printf("created bucket %s\n", bucket)

	return nil
}

func objectExists(ctx context.Context, mc *minio.Client, bucket, object string) (bool, error) {
	_, err := mc.StatObject(bucket, object, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
