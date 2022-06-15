package handler

import (
	"time"

	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/model"
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
		Plugin: model.Plugin{
			Name:            "hydrograph_scaler",
			ImageAndTag:     "williamlehman/hydrographscaler:v0.0.11",
			CommandLineArgs: []string{"./main", "-payload"},
		},
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
		Plugin: model.Plugin{
			Name:            "ras-mutator",
			ImageAndTag:     "lawlerseth/ras-mutator:v0.1.1",
			CommandLineArgs: []string{"./h5rasedit", "wat", "-m", "host.docker.internal:9000", "-f"},
		},
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
		Plugin: model.Plugin{
			Name:            "ras-unsteady",
			ImageAndTag:     "lawlerseth/ras-unsteady:v0.0.2",
			CommandLineArgs: []string{"./watrun", "-m", "host.docker.internal:9000", "-f"},
		},
	}
	return model.DirectedAcyclicGraph{
		Nodes: manifests,
	}
}
func MockStochasticJob(config config.WatConfig) model.StochasticJob {
	tw := model.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	dag := MockDag()
	sj := model.StochasticJob{
		Dag:                    dag, //yo
		TimeWindow:             tw,
		TotalRealizations:      1,
		EventsPerRealization:   1,
		InitialRealizationSeed: 1234,
		InitialEventSeed:       1234,
		Outputdestination: model.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Fragment:  "/runs/",
		},
		Inputsource: model.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Fragment:  "/data/",
		},
		DeleteOutputAfterRealization: false,
	}
	return sj
}
func MockStochastic2dJob(config config.WatConfig) model.StochasticJob {
	tw := model.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	dag := MockDag()
	sj := model.StochasticJob{
		Dag:                    dag, //yo
		TimeWindow:             tw,
		TotalRealizations:      1,
		EventsPerRealization:   1,
		InitialRealizationSeed: 1234,
		InitialEventSeed:       1234,
		Outputdestination: model.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Fragment:  "/runs2d/",
		},
		Inputsource: model.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Fragment:  "/data/",
		},
		DeleteOutputAfterRealization: false,
	}
	return sj
}
