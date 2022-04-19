package yandex

import (
	"pulumi-crosswalks/utils"

	"github.com/pulumi/pulumi-yandex/sdk/go/yandex"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func InitProvider(ctx *pulumi.Context, creds pulumi.AnyOutput) (*yandex.Provider, error) {
	provider, err := yandex.NewProvider(ctx, "provider", &yandex.ProviderArgs{
		StorageAccessKey: utils.ExtractStringFromPulumiMap(creds, "access_key"),
		StorageSecretKey: utils.ExtractStringFromPulumiMap(creds, "secret_key"),
		Token:            utils.ExtractStringFromPulumiMap(creds, "apikey"),
		FolderId:         utils.ExtractStringFromPulumiMap(creds, "folder_id"),
	})
	if err != nil {
		return nil, err
	}

	return provider, nil
}
