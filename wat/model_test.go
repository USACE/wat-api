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
	bytes, err := yaml.Marshal(mm)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))

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

func TestHSMModelManifestSeralization(t *testing.T) {

	inputs := make([]Input, 0)
	outputs := make([]Output, 3)
	outputs[0] = Output{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	outputs[1] = Output{
		Name:      "hsm2.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	outputs[2] = Output{
		Name:      "hsm3.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	paths := make([]ResourceInfo, 1)
	paths[0] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/hsm-Test",
		Fragment:  "hsm.json",
	}

	mc := ModelConfiguration{
		Name:                        "hsm",
		ModelConfigurationResources: paths,
	}
	var mincpus int64 = 1
	var maxcpus int64 = 4
	var desiredcpus int64 = 2
	computeType := "EC2"
	instances := make([]*string, 1)
	instance := "m2.micro"
	instances[0] = &instance
	computeResoures := ModelComputeResources{
		MinCpus:       &mincpus,
		MaxCpus:       &maxcpus,
		DesiredCpus:   &desiredcpus,
		Type:          &computeType,
		InstanceTypes: instances,
		Managed:       true,
	}
	mm := ModelManifest{
		Plugin:                Plugin{Name: "hydrographscaler", ImageAndTag: "williamlehman/hydrographscaler:v0.0.1"},
		ModelConfiguration:    mc,
		ModelComputeResources: computeResoures,
		Inputs:                inputs,
		Outputs:               outputs,
	}
	bytes, err := yaml.Marshal(mm)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
func TestHSMModelPayloadSeralization(t *testing.T) {
	tw := TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	event := IndexedSeed{Index: 1, Seed: 5678}
	realization := IndexedSeed{Index: 1, Seed: 1234}
	eventConfiguration := EventConfiguration{
		OutputDestination: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
		},
		Realization:     realization,
		Event:           event,
		EventTimeWindow: tw,
	}
	//someone has to make data somewhere... probably needs to be computed output
	prevModelOutput := make([]ComputedOutput, 0)
	outputs := make([]Output, 1)
	outputs[0] = Output{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	paths := make([]ResourceInfo, 1)
	paths[0] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/hsm-Test",
		Fragment:  "hsm.json",
	}
	mc := ModelConfiguration{
		Name:                        "hsm",
		ModelConfigurationResources: paths,
	}
	ml := ModelLinks{
		LinkedInputs:     prevModelOutput,
		NecessaryOutputs: outputs,
	}
	mPayload := ModelPayload{
		TargetPlugin:       "hydrographscaler",
		PluginImageAndTag:  "williamlehman/hydrographscaler:v0.0.1",
		ModelConfiguration: mc,
		ModelLinks:         ml,
		EventConfiguration: eventConfiguration,
	}
	bytes, err := yaml.Marshal(mPayload)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
func TestRASMutatorModelManifestSeralization(t *testing.T) {

	inputs := make([]Input, 1)
	inputs[0] = Input{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	outputs := make([]Output, 4)
	outputs[0] = Output{
		Name:      "muncie.p04.tmp.hdf",
		Parameter: "ras p hdf file",
		Format:    "hdf",
	}
	outputs[1] = Output{
		Name:      "muncie.b04",
		Parameter: "ras b file",
		Format:    ".b**",
	}
	outputs[2] = Output{
		Name:      "muncie.prj",
		Parameter: "ras project file",
		Format:    ".prj",
	}
	outputs[3] = Output{
		Name:      "muncie.x04",
		Parameter: "ras x file",
		Format:    ".x**",
	}
	paths := make([]ResourceInfo, 4)
	paths[0] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.p04.tmp.hdf",
	}
	paths[1] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.b04",
	}
	paths[2] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.prj",
	}
	paths[3] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.x04",
	}
	mc := ModelConfiguration{
		Name:                        "Muncie",
		ModelConfigurationResources: paths,
	}
	var mincpus int64 = 1
	var maxcpus int64 = 4
	var desiredcpus int64 = 2
	computeType := "EC2"
	instances := make([]*string, 1)
	instance := "m2.micro"
	instances[0] = &instance
	computeResoures := ModelComputeResources{
		MinCpus:       &mincpus,
		MaxCpus:       &maxcpus,
		DesiredCpus:   &desiredcpus,
		Type:          &computeType,
		InstanceTypes: instances,
		Managed:       true,
	}
	mm := ModelManifest{
		Plugin:                Plugin{Name: "ras-mutator", ImageAndTag: "williamlehman/ras-mutator:v0.0.1"},
		ModelConfiguration:    mc,
		ModelComputeResources: computeResoures,
		Inputs:                inputs,
		Outputs:               outputs,
	}
	bytes, err := yaml.Marshal(mm)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
func TestRASMutatorModelPayloadSeralization(t *testing.T) {
	tw := TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	event := IndexedSeed{Index: 1, Seed: 5678}
	realization := IndexedSeed{Index: 1, Seed: 1234}
	eventConfiguration := EventConfiguration{
		OutputDestination: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
		},
		Realization:     realization,
		Event:           event,
		EventTimeWindow: tw,
	}
	//someone has to make data somewhere... probably needs to be computed output
	prevModelOutput := make([]ComputedOutput, 1)
	prevModelOutput[0] = ComputedOutput{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    "csv",
		ResourceInfo: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
			Fragment:  "hsm1.csv",
		},
	}
	inputs := make([]Input, 1)
	inputs[0] = Input{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
	}

	outputs := make([]Output, 4)
	outputs[0] = Output{
		Name:      "muncie.p04.tmp.hdf",
		Parameter: "ras p hdf file",
		Format:    "hdf",
	}
	outputs[1] = Output{
		Name:      "muncie.b04",
		Parameter: "ras b file",
		Format:    ".b**",
	}
	outputs[2] = Output{
		Name:      "muncie.prj",
		Parameter: "ras project file",
		Format:    ".prj",
	}
	outputs[3] = Output{
		Name:      "muncie.x04",
		Parameter: "ras x file",
		Format:    ".x**",
	}

	paths := make([]ResourceInfo, 4)
	paths[0] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.p04.tmp.hdf",
	}
	paths[1] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.b04",
	}
	paths[2] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.prj",
	}
	paths[3] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.x04",
	}
	mc := ModelConfiguration{
		Name:                        "Muncie",
		ModelConfigurationResources: paths,
	}
	ml := ModelLinks{
		LinkedInputs:     prevModelOutput,
		NecessaryOutputs: outputs,
	}
	mPayload := ModelPayload{
		TargetPlugin:       "ras-mutator",
		PluginImageAndTag:  "williamlehman/ras-mutator:v0.0.1",
		ModelConfiguration: mc,
		ModelLinks:         ml,
		EventConfiguration: eventConfiguration,
	}
	bytes, err := yaml.Marshal(mPayload)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
func TestRASRunnerModelManifestSeralization(t *testing.T) {
	inputs := make([]Input, 0)
	outputs := make([]Output, 2)
	outputs[0] = Output{
		Name:      "muncie.p04.hdf",
		Parameter: "ras results hdf file",
		Format:    ".hdf",
	}
	outputs[1] = Output{
		Name:      "muncie.log",
		Parameter: "ras log file",
		Format:    ".log",
	}
	paths := make([]ResourceInfo, 4)
	paths[0] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.p04.tmp.hdf",
	}
	paths[1] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.b04",
	}
	paths[2] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.prj",
	}
	paths[3] = ResourceInfo{
		Scheme:    "https",
		Authority: "/model-library/Muncie-Test",
		Fragment:  "muncie.x04",
	}
	mc := ModelConfiguration{
		Name:                        "Muncie",
		ModelConfigurationResources: paths,
	}
	var mincpus int64 = 1
	var maxcpus int64 = 4
	var desiredcpus int64 = 2
	computeType := "EC2"
	instances := make([]*string, 1)
	instance := "m2.micro"
	instances[0] = &instance
	computeResoures := ModelComputeResources{
		MinCpus:       &mincpus,
		MaxCpus:       &maxcpus,
		DesiredCpus:   &desiredcpus,
		Type:          &computeType,
		InstanceTypes: instances,
		Managed:       true,
	}
	mm := ModelManifest{
		Plugin:                Plugin{Name: "ras-runner", ImageAndTag: "williamlehman/ras-runner:v0.0.1"},
		ModelConfiguration:    mc,
		ModelComputeResources: computeResoures,
		Inputs:                inputs,
		Outputs:               outputs,
	}
	bytes, err := yaml.Marshal(mm)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
func TestRASRunnerModelPayloadSeralization(t *testing.T) {
	tw := TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	event := IndexedSeed{Index: 1, Seed: 5678}
	realization := IndexedSeed{Index: 1, Seed: 1234}
	eventConfiguration := EventConfiguration{
		OutputDestination: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
		},
		Realization:     realization,
		Event:           event,
		EventTimeWindow: tw,
	}
	//someone has to make data somewhere... probably needs to be computed output
	prevModelOutput := make([]ComputedOutput, 0)
	outputs := make([]Output, 2)
	outputs[0] = Output{
		Name:      "muncie.p04.hdf",
		Parameter: "ras results hdf file",
		Format:    ".hdf",
	}
	outputs[1] = Output{
		Name:      "muncie.log",
		Parameter: "ras log file",
		Format:    ".log",
	}
	paths := make([]ResourceInfo, 4)
	paths[0] = ResourceInfo{
		Scheme:    "http",
		Authority: "/minio/runs/realization_1/event_1",
		Fragment:  "muncie.p04.tmp.hdf",
	}
	paths[1] = ResourceInfo{
		Scheme:    "http",
		Authority: "/minio/runs/realization_1/event_1",
		Fragment:  "muncie.b04",
	}
	paths[2] = ResourceInfo{
		Scheme:    "http",
		Authority: "/minio/runs/realization_1/event_1",
		Fragment:  "muncie.prj",
	}
	paths[3] = ResourceInfo{
		Scheme:    "http",
		Authority: "/minio/runs/realization_1/event_1",
		Fragment:  "muncie.x04",
	}
	mc := ModelConfiguration{
		Name:                        "Muncie",
		ModelConfigurationResources: paths,
	}
	ml := ModelLinks{
		LinkedInputs:     prevModelOutput,
		NecessaryOutputs: outputs,
	}
	mPayload := ModelPayload{
		TargetPlugin:       "ras-runner",
		PluginImageAndTag:  "williamlehman/ras-runner:v0.0.1",
		ModelConfiguration: mc,
		ModelLinks:         ml,
		EventConfiguration: eventConfiguration,
	}
	bytes, err := yaml.Marshal(mPayload)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
