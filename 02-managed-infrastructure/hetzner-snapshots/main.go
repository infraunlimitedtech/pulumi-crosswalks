package main

import (
        "fmt"
        "strconv"
	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		infraStack, err := pulumi.NewStackReference(ctx, "spigell/managed-infrastructure/hetzner-test", nil)
		if err != nil {
			return err
		}

		s := infraStack.GetOutput(pulumi.String("infra:nodes:info"))

                m := s.ApplyT(func (i interface{}) (map[string]interface{}, error) {
			m := i.(map[string]interface{})

                        res := make(map[string]interface{})

			for k, v := range m {
				srv := v.(map[string]interface{})

				id, err := strconv.Atoi(srv["id"].(string))
				if err != nil {
					return res, err
				}

	                        s, err := hcloud.NewSnapshot(ctx, fmt.Sprintf("snap-%s", k), &hcloud.SnapshotArgs{
					ServerId: pulumi.Int(id),
	                                Description: pulumi.String(k),
				})
				if err != nil {
					return res, err
				}

                                res[k] = s.ID()

			}
			return res, nil
		})
	        ctx.Export("nodes", pulumi.Unsecret(m))

		return nil
	})
}
