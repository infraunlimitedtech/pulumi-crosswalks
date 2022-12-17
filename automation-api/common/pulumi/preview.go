package pulumi

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
)

func (p *Pulumi) Preview(ctx context.Context, realtimeSteamToStdout bool) ([]events.EngineEvent, error) {
	var previewEvents []events.EngineEvent
	prevCh := make(chan events.EngineEvent)

	opts := []optpreview.Option{}
	opts = append(opts, optpreview.EventStreams(prevCh))

	if realtimeSteamToStdout {
		opts = append(opts, optpreview.ProgressStreams(os.Stdout))
	}

	wg := collectEvents(prevCh, &previewEvents)

	_, err := p.stack.Preview(ctx, opts...)
	if err != nil {
		return previewEvents, fmt.Errorf("get preview: %w", err)
	}

	wg.Wait()
	return previewEvents, nil
}
