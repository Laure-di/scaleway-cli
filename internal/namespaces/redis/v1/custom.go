package redis

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(clusterWaitCommand()))
	cmds.MustFind("redis", "cluster", "create").Override(clusterCreateBuilder)
	cmds.MustFind("redis", "cluster", "delete").Override(clusterDeleteBuilder)
	cmds.MustFind("redis", "cluster", "metrics").Override(clusterMetricsBuilder)

	//human.RegisterMarshalerFunc(redis.ClusterMetricsResponse{}, clusterMetricsMarshalerFunc)

	return cmds
}
