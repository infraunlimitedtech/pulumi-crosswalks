package yandex

import (
	"errors"
	"fmt"

	"github.com/pulumi/pulumi-yandex/sdk/go/yandex"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type S3Config struct {
	Buckets []*Bucket
}

type Bucket struct {
	ID     string
	ACL    string
	Bucket string
	Prefix string
}

type S3Infra struct {
	storage map[string]map[string]interface{}
}

func ManageS3(ctx *pulumi.Context, cfg *S3Config, creds pulumi.AnyOutput, p *yandex.Provider) (*S3Infra, error) {
	s := make(map[string]map[string]interface{})
	for _, bucket := range cfg.Buckets {
		args := &yandex.StorageBucketArgs{
			Acl: pulumi.String(bucket.ACL),
		}

		if bucket.ID == "" {
			return nil, errors.New("please specify id for bucket")
		}

		switch {
		case bucket.Prefix != "":
			args.BucketPrefix = pulumi.String(bucket.Prefix)
		case bucket.Bucket != "":
			args.Bucket = pulumi.String(bucket.Bucket)
		default:
			return nil, errors.New("please specify prefix or bucket")
		}

		res, err := yandex.NewStorageBucket(ctx, fmt.Sprintf("%s-bucket", bucket.ID), args, pulumi.Provider(p))
		if err != nil {
			return nil, err
		}
		s[bucket.ID] = make(map[string]interface{})
		s[bucket.ID]["url"] = res.BucketDomainName
	}
	return &S3Infra{
		storage: s,
	}, nil
}

func (s *S3Infra) GetStorage() map[string]map[string]interface{} {
	return s.storage
}
