package wat

import (
	"testing"
	"time"

	"github.com/usace/wat-api/config"
)

func TestStochasticPayloadGeneration(t *testing.T) {
	tw := TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	sj := StochasticJob{

		TimeWindow:                   tw,
		TotalRealizations:            2,
		EventsPerRealization:         10,
		InitialRealizationSeed:       1234,
		InitialEventSeed:             1234,
		Outputdestination:            "testing",
		Inputsource:                  "testSettings.InputDataDir",
		DeleteOutputAfterRealization: false,
	}
	config := config.WatConfig{}
	_, err := sj.GeneratePayloads(nil, nil, nil, config)
	if err != nil {
		t.Fail()
	}
}
