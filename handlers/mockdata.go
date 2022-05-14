package handler

import (
	"time"

	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/wat"
)

func MockDag() wat.DirectedAcyclicGraph {
	manifests := make([]wat.ModelManifest, 2)
	t := "EC2"
	i := "m2.micro"
	var min int64 = 0
	var desired int64 = 2
	var max int64 = 128
	instance_types := make([]*string, 1)
	instance_types[0] = &i
	manifests[0] = wat.ModelManifest{
		ModelComputeResources: wat.ModelComputeResources{
			MinCpus:       &min,
			DesiredCpus:   &desired,
			MaxCpus:       &max,
			InstanceTypes: instance_types,
			Type:          &t,
			Managed:       true,
		},
		Plugin: wat.Plugin{Name: "fragilitycurveplugin", ImageAndTag: "williamlehman/fragilitycurveplugin:v0.0.7"},
	}
	manifests[1] = wat.ModelManifest{
		ModelComputeResources: wat.ModelComputeResources{
			MinCpus:       &min,
			DesiredCpus:   &desired,
			MaxCpus:       &max,
			InstanceTypes: instance_types,
			Type:          &t,
			Managed:       true,
		},
		Plugin: wat.Plugin{Name: "hydrograph_scaler", ImageAndTag: "williamlehman/hydrographscaler:v0.0.7"},
	}
	return wat.DirectedAcyclicGraph{
		Nodes: manifests,
	}
}
func MockStochasticJob(config config.WatConfig) wat.StochasticJob {
	tw := wat.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	dag := MockDag()
	sj := wat.StochasticJob{
		Dag:                    dag, //yo
		TimeWindow:             tw,
		TotalRealizations:      2,
		EventsPerRealization:   10,
		InitialRealizationSeed: 1234,
		InitialEventSeed:       1234,
		Outputdestination: wat.ResourceInfo{
			Scheme:    config.S3_ENDPOINT + "/" + config.S3_BUCKET,
			Authority: "/runs/",
		},
		Inputsource: wat.ResourceInfo{
			Scheme:    config.S3_ENDPOINT + "/" + config.S3_BUCKET,
			Authority: "/data/",
		},
		DeleteOutputAfterRealization: false,
	}
	return sj
}
