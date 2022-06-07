package handler

import (
	"time"

	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/model"
	"github.com/usace/wat-api/wat"
)

func MockDag() model.DirectedAcyclicGraph {
	manifests := make([]model.ModelManifest, 3)
	t := "EC2"
	i := "m2.micro"
	var min int64 = 0
	var desired int64 = 2
	var max int64 = 128
	instance_types := make([]*string, 1)
	instance_types[0] = &i
	/*manifests[0] = model.ModelManifest{
		ModelComputeResources: model.ModelComputeResources{
			MinCpus:       &min,
			DesiredCpus:   &desired,
			MaxCpus:       &max,
			InstanceTypes: instance_types,
			Type:          &t,
			Managed:       true,
		},
		Plugin: model.Plugin{Name: "fragilitycurveplugin", ImageAndTag: "williamlehman/fragilitycurveplugin:v0.0.7"},
	}*/
	manifests[0] = model.ModelManifest{
		ModelComputeResources: model.ModelComputeResources{
			MinCpus:       &min,
			DesiredCpus:   &desired,
			MaxCpus:       &max,
			InstanceTypes: instance_types,
			Type:          &t,
			Managed:       true,
		},
		Plugin: model.Plugin{Name: "hydrograph_scaler", ImageAndTag: "williamlehman/hydrographscaler:v0.0.8"},
	}
	manifests[1] = model.ModelManifest{
		ModelComputeResources: model.ModelComputeResources{
			MinCpus:       &min,
			DesiredCpus:   &desired,
			MaxCpus:       &max,
			InstanceTypes: instance_types,
			Type:          &t,
			Managed:       true,
		},
		Plugin: model.Plugin{Name: "ras-mutator", ImageAndTag: "lawlerseth/ras-mutator:v0.1.0"},
	}
	manifests[2] = model.ModelManifest{
		ModelComputeResources: model.ModelComputeResources{
			MinCpus:       &min,
			DesiredCpus:   &desired,
			MaxCpus:       &max,
			InstanceTypes: instance_types,
			Type:          &t,
			Managed:       true,
		},
		Plugin: model.Plugin{Name: "ras-unsteady", ImageAndTag: "lawlerseth/ras-unsteady:v0.1.0"},
	}
	return model.DirectedAcyclicGraph{
		Nodes: manifests,
	}
}
func MockStochasticJob(config config.WatConfig) wat.StochasticJob {
	tw := model.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	dag := MockDag()
	sj := wat.StochasticJob{
		Dag:                    dag, //yo
		TimeWindow:             tw,
		TotalRealizations:      2,
		EventsPerRealization:   10,
		InitialRealizationSeed: 1234,
		InitialEventSeed:       1234,
		Outputdestination: model.ResourceInfo{
			Scheme:    "s3",
			Authority: "configs",
			Fragment:  "/runs/",
		},
		Inputsource: model.ResourceInfo{
			Scheme:    "s3",
			Authority: "configs",
			Fragment:  "/data/",
		},
		DeleteOutputAfterRealization: false,
	}
	return sj
}
