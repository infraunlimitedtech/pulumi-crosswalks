package yandex

import (
	"github.com/pulumi/pulumi-yandex/sdk/go/yandex"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ManageS3Helper(ctx *pulumi.Context) (map[string]interface{}, error) {
	s3Helper, err := yandex.NewIamServiceAccount(ctx, "s3Helper", &yandex.IamServiceAccountArgs{
		Name:        pulumi.String("pulumi-s3-helper"),
		Description: pulumi.String("service account to manage storages. Managed by Pulumi"),
	})
	if err != nil {
		return nil, err
	}

	apiKey, err := yandex.NewIamServiceAccountApiKey(ctx, "s3Helper", &yandex.IamServiceAccountApiKeyArgs{
		Description:      pulumi.String("key for s3 helper. Managed by Pulumi"),
		ServiceAccountId: s3Helper.ID(),
	})
	if err != nil {
		return nil, err
	}

	_, err = yandex.NewResourcemanagerFolderIamMember(ctx, "admin", &yandex.ResourcemanagerFolderIamMemberArgs{
		FolderId: s3Helper.FolderId,
		Member:   pulumi.Sprintf("serviceAccount:%s", s3Helper.ID()),
		Role:     pulumi.String("storage.admin"),
	})
	if err != nil {
		return nil, err
	}

	s3Key, err := yandex.NewIamServiceAccountStaticAccessKey(ctx, "s3Helper", &yandex.IamServiceAccountStaticAccessKeyArgs{
		Description:      pulumi.String("static access key for object storage. Managed by Pulumi"),
		ServiceAccountId: s3Helper.ID(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"folder_id":  s3Helper.FolderId,
		"account_id": s3Helper.ID(),
		"access_key": s3Key.AccessKey,
		"secret_key": s3Key.SecretKey,
		"apikey":     apiKey.SecretKey,
	}, nil
}
