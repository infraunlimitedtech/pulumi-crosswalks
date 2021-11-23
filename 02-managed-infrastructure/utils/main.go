package utils

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ExtractFromExportedMap(p pulumi.Output, key string) pulumi.StringOutput {
	return p.ApplyT(func(v interface{}) string {
		m, ok := v.(map[string]interface{})
		if !ok {
			panic("Exported Map invalid!")
		}
		return m[key].(string)
	}).(pulumi.StringOutput)
}
