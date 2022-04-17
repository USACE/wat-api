package wat

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

func TestModelManifestSeralization(t *testing.T) {
	inputs := make([]Input, 0)
	outputs := make([]Output, 1)
	outputs[0] = Output{
		Name:      "hydrograph1",
		Parameter: "flow",
		Format:    "csv",
	}
	paths := make([]string, 1)
	paths[0] = "/hsm.json"
	mc := ModelConfiguration{
		Name:                    "TestModel",
		ModelConfigurationPaths: paths,
	}
	mm := ModelManifest{
		ModelConfiguration: mc,
		Inputs:             inputs,
		Outputs:            outputs,
	}
	bytes, err := json.Marshal(mm)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	t.Log(string(bytes))

}
func TestModelPayloadSeralization(t *testing.T) {
	tw := TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	event := IndexedSeed{Index: 1, Seed: 5678}
	realization := IndexedSeed{Index: 1, Seed: 1234}
	eventConfiguration := EventConfiguration{
		OutputDestination: "/testing/",
		Realization:       realization,
		Event:             event,
		EventTimeWindow:   tw,
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
	paths := make([]string, 1)
	paths[0] = "/hsm.json"
	mc := ModelConfiguration{
		Name:                    "TestModel",
		ModelConfigurationPaths: paths,
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
	t.Log(string(bytes))
	t.Log("\n")
	ybytes, err := yaml.Marshal(mmanifest)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	t.Log(string(ybytes))
}
