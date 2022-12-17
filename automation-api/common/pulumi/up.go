package pulumi

import (
	"context"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

func (p *Pulumi) Up(ctx context.Context) error {
	opts := []optup.Option{}
	opts = append(opts, optup.ProgressStreams(os.Stdout))

	_, err := p.stack.Up(ctx, opts...)

	return err
}
