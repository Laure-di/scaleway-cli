package k8s

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type k8sVersionGetRequest struct {
	Version string
	Region  scw.Region
}

func k8sVersionGetCommand() *core.Command {
	return &core.Command{
		Short:     `Get a Kubernetes version`,
		Long:      `Get all the details about a specific Kubernetes version.`,
		Namespace: "k8s",
		Verb:      "get",
		Resource:  "version",
		ArgsType:  reflect.TypeOf(k8sVersionGetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "version",
				Short:      "Version from which to get details",
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
		},
		Run: k8sVersionGetRun,
	}
}

func k8sVersionGetRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	request := argsI.(*k8sVersionGetRequest)

	client := core.ExtractClient(ctx)
	apiK8s := k8s.NewAPI(client)

	versions, err := apiK8s.ListVersions(&k8s.ListVersionsRequest{
		Region: request.Region,
	})

	if err != nil {
		return nil, err
	}

	for _, version := range versions.Versions {
		if version.Name == request.Version {
			return version, nil
		}
	}
	return nil, fmt.Errorf("version '%s' not found", request.Version)
}