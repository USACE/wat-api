package handler

import (
	"time"

	"github.com/usace/wat-api/wat"
)

func MockStochasticJob() wat.StochasticJob {
	tw := wat.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	plugins := make([]wat.Plugin, 1)
	plugins[0] = wat.Plugin{Name: "fragilitycurveplugin", ImageAndTag: "williamlehman/fragilitycurveplugin:v0.0.2"}
	sj := wat.StochasticJob{
		SelectedPlugins:              plugins,
		TimeWindow:                   tw,
		TotalRealizations:            2,
		EventsPerRealization:         10,
		InitialRealizationSeed:       1234,
		InitialEventSeed:             1234,
		Outputdestination:            "/data/",
		Inputsource:                  "/data/",
		DeleteOutputAfterRealization: false,
	}
	return sj
}
