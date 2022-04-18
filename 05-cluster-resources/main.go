package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"cluster-resources/services"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		infra, err := services.Init(ctx)
		if err != nil {
			return err
		}

		if err := infra.RunNginxIngress(); err != nil {
			return err
		}

//		if err := infra.RunPrometheus(); err != nil {
//			return err
//		}

		return infra.RunConsulStack()
	})
}
