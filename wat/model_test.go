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
	paths := make([]ResourceInfo, 1)
	paths[0] = ResourceInfo{Fragment: "/hsm.json"}
	mc := ModelConfiguration{
		Name:                        "TestModel",
		ModelConfigurationResources: paths,
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
		OutputDestination: ResourceInfo{
			Authority: "/testing/",
		},
		Realization:     realization,
		Event:           event,
		EventTimeWindow: tw,
	}
	//someone has to make data somewhere... probably needs to be computed output
	prevModelOutput := make([]ComputedOutput, 2)
	prevModelOutput[0] = ComputedOutput{
		Name:      "OutputFromAnotherModel1",
		Parameter: "speed",
		Format:    "mph",
	}
	prevModelOutput[1] = ComputedOutput{
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
	paths := make([]ResourceInfo, 1)
	paths[0] = ResourceInfo{Fragment: "/hsm.json"}
	mc := ModelConfiguration{
		Name:                        "TestModel",
		ModelConfigurationResources: paths,
	}
	/*m := Model{
		ModelConfiguration: mc,
		Inputs:                 inputs,
		Outputs:                outputs,
	}*/
	linkedInputs := make([]ComputedOutput, 2)
	linkedInputs[0] = ComputedOutput{
		Name:      inputs[0].Name,
		Format:    inputs[0].Format,
		Parameter: inputs[0].Parameter,
		ResourceInfo: ResourceInfo{
			Scheme:    "s3://",
			Authority: "testing/",
			Fragment:  inputs[0].Name,
		},
	}
	linkedInputs[1] = ComputedOutput{
		Name:      inputs[1].Name,
		Format:    inputs[1].Format,
		Parameter: inputs[1].Parameter,
		ResourceInfo: ResourceInfo{
			Scheme:    "s3://",
			Authority: "testing/",
			Fragment:  inputs[1].Name,
		},
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
