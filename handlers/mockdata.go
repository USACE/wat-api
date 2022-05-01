package handler

import (
	"time"

	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/wat"
)

func MockPlugins() []wat.Plugin {
	plugins := make([]wat.Plugin, 3)
	plugins[0] = wat.Plugin{Name: "fragilitycurveplugin", ImageAndTag: "williamlehman/fragilitycurveplugin:v0.0.4"}
	plugins[1] = wat.Plugin{Name: "hydrograph_scaler", ImageAndTag: "williamlehman/hydrographscaler:v0.0.5"}
	plugins[2] = wat.Plugin{Name: "hydrograph_stats", ImageAndTag: "tbd/hydrographstats:v0.0.2"}
	return plugins
}
func MockStochasticJob(config config.WatConfig) wat.StochasticJob {
	tw := wat.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	plugins := MockPlugins()
	sj := wat.StochasticJob{
		SelectedPlugins:        plugins,
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
