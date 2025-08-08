package block

import (
	"context"
	"github.com/scaleway/scaleway-cli/v2/core"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
)

type customVolumeType struct {
	*block.VolumeType
	KgCo2Equivalent *float32 `json:"kg_co2_equivalent"`
	M3WaterUsage    *float32 `json:"m3_water_usage"`
}

func blockVolumeTypeListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "Type",
				FieldName: "type",
			},
			{
				Label:     "Pricing (GB/hour)",
				FieldName: "pricing",
			},
			{
				Label:     "Snapshot pricing (GB/hour)",
				FieldName: "snapshot_pricing",
			},
			{
				Label:     "KgCo2Equivalent (kg/hour)",
				FieldName: "kg_co2_equivalent",
			},
			{
				Label:     "M3WaterUsage (m3/hour)",
				FieldName: "m3_water_usage",
			},
		},
	}
	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			originalRes, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}

			client := core.ExtractClient(ctx)

			serverTypes := originalRes.(*block.VolumeType)

		},
	)
	return c
}
