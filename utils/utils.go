package utils

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ExtractValueFromPulumiMapMap(m pulumi.AnyOutput, mapKey, valueKey string) pulumi.StringOutput {
	return m.ApplyT(func(v interface{}) string {
		d, ok := v.(map[string]interface{})
		if !ok {
			panic(fmt.Sprintf("There is no map with key `%s`", mapKey))
		}
		e, ok := d[mapKey].(map[string]interface{})
		if !ok {
			panic(fmt.Sprintf("Can't find values for mapKey `%s`. Is it contains in parent map?", mapKey))
		}
		return e[valueKey].(string)
	}).(pulumi.StringOutput)
}

func ConvertMapSliceToSliceByKey(s []map[string]pulumi.Resource, key string) (r []pulumi.Resource) {
	for _, m := range s {
		r = append(r, m[key])
	}
	return
}

func ExtractStringFromPulumiMap(m pulumi.AnyOutput, valueKey string) pulumi.StringOutput {
	return m.ApplyT(func(v interface{}) string {
		e, ok := v.(map[string]interface{})
		if !ok {
			panic("It is not a map!")
		}
		return e[valueKey].(string)
	}).(pulumi.StringOutput)
}
