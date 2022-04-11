package wat

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestModelManifestSeralization(t *testing.T) {
	tw := TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	eventConfiguration := EventConfiguration{
		OutputDestination: "/testing/",
		RealizationNumber: 1,
		EventNumber:       1,
		EventTimeWindow:   tw,
		RealizationSeed:   1234,
		EventSeed:         5678,
	}
	//someone has to make data somewhere...
	prevModelOutput := make([]Output, 2)
	prevModelOutput[0] = Output{
		Name:      "OutputFromAnotherModel1",
		Parameter: "speed",
		Format:    "mph",
	}
	prevModelOutput[1] = Output{
		Name:      "OutputFromAnotherModel2",
		Parameter: "distance",
		Format:    "mi",
	}
	inputs := make([]Input, 2)
	inputs[0] = Input{
		Name:      "input1",
		Parameter: "speed",
		Format:    "mph",
	}
	inputs[1] = Input{
		Name:      "input2",
		Parameter: "distance",
		Format:    "mi",
	}
	outputs := make([]Output, 1)
	outputs[0] = Output{
		Name:      "output1",
		Parameter: "time",
		Format:    "hours",
	}
	mc := ModelConfiguration{
		Name:                   "TestModel",
		ModelConfigurationPath: "/hsm.json",
	}
	/*m := Model{
		ModelConfiguration: mc,
		Inputs:                 inputs,
		Outputs:                outputs,
	}*/
	linkedInputs := make([]LinkedInput, 2)
	linkedInputs[0] = LinkedInput{
		Input:  inputs[0],
		Source: prevModelOutput[0],
	}
	linkedInputs[1] = LinkedInput{
		Input:  inputs[1],
		Source: prevModelOutput[1],
	}
	ml := ModelLinks{
		LinkedInputs:     linkedInputs,
		NecessaryOutputs: outputs,
	}
	mmanifest := ModelPayload{
		TargetPlugin:       "SpeedAndDistanceToTimePlugin",
		ModelConfiguration: mc,
		ModelLinks:         ml,
		EventConfiguration: eventConfiguration,
	}
	bytes, err := json.Marshal(mmanifest)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	log.Fatal(string(bytes))
}
