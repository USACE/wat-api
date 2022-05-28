package wat

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"gopkg.in/yaml.v2"
)

func TestEventConfiguration(t *testing.T) {
	eventConfiguration := MockEventConfiguration()
	bytes, err := json.Marshal(eventConfiguration)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
func TestModelManifestSeralization(t *testing.T) {
	inputs := make([]Input, 1)
	inputs[0] = Input{
		Name:      "Project File",
		Parameter: "Project Specification",
		Format:    ".json",
	}
	outputs := make([]Output, 1)
	outputs[0] = Output{
		Name:      "hydrograph1",
		Parameter: "flow",
		Format:    ".csv",
	}
	mc := ModelConfiguration{
		Name: "TestModel",
		//Alternative: "",
	}

	mm := ModelManifest{
		Plugin: Plugin{
			Name:        "hydrographscaler",
			ImageAndTag: "williamlehman/hydrographscaler:v0.0.1",
		},
		ModelConfiguration: mc,
		ModelComputeResources: ModelComputeResources{
			MinCpus:       aws.Int64(0),
			DesiredCpus:   aws.Int64(1),
			MaxCpus:       aws.Int64(1),
			InstanceTypes: []*string{aws.String("m2.micro")},
			Type:          aws.String("EC2"),
			Managed:       false,
		},
		Inputs:  inputs,
		Outputs: outputs,
	}
	bytes, err := yaml.Marshal(mm)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))

}
func TestModelPayloadSeralization(t *testing.T) {
	eventConfiguration := MockEventConfiguration()
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
	mc := ModelConfiguration{
		Name: "TestModel",
	}
	linkedInputs := make([]ComputedOutput, 3)
	linkedInputs[0] = ComputedOutput{
		Name:      "Project File",
		Format:    "Project Specification",
		Parameter: ".json",
		ResourceInfo: ResourceInfo{
			Scheme:    "s3://",
			Authority: "testing/",
			Fragment:  "hsm.json",
		},
	}
	linkedInputs[1] = ComputedOutput{
		Name:      inputs[0].Name,
		Format:    inputs[0].Format,
		Parameter: inputs[0].Parameter,
		ResourceInfo: ResourceInfo{
			Scheme:    "s3://",
			Authority: "testing/",
			Fragment:  inputs[0].Name,
		},
	}
	linkedInputs[2] = ComputedOutput{
		Name:      inputs[1].Name,
		Format:    inputs[1].Format,
		Parameter: inputs[1].Parameter,
		ResourceInfo: ResourceInfo{
			Scheme:    "s3://",
			Authority: "testing/",
			Fragment:  inputs[1].Name,
		},
	}
	linkedInputs = append(linkedInputs, eventConfiguration.ToInput())
	ml := ModelLinks{
		LinkedInputs:     linkedInputs,
		NecessaryOutputs: outputs,
	}
	mmanifest := ModelPayload{
		ModelConfiguration: mc,
		ModelLinks:         ml,
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

	inputs := make([]Input, 1)
	inputs[0] = Input{
		Name:      "Project File",
		Parameter: "Project Specification",
		Format:    ".json",
	}
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
	mc := ModelConfiguration{
		Name: "hsm",
		//ModelConfigurationResources: paths,
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
	eventConfiguration := MockEventConfiguration()
	//someone has to make data somewhere... probably needs to be computed output
	prevModelOutput := make([]ComputedOutput, 2)
	prevModelOutput[0] = ComputedOutput{
		Name:      "Project File",
		Parameter: "Project Specification",
		Format:    ".json",
		ResourceInfo: ResourceInfo{
			Scheme:    "https",
			Authority: "/model-library/hsm-Test",
			Fragment:  "hsm.json",
		},
	}
	prevModelOutput[1] = eventConfiguration.ToInput()
	outputs := make([]Output, 1)
	outputs[0] = Output{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	mc := ModelConfiguration{
		Name: "hsm",
	}
	ml := ModelLinks{
		LinkedInputs:     prevModelOutput,
		NecessaryOutputs: outputs,
	}
	mPayload := ModelPayload{
		ModelConfiguration: mc,
		ModelLinks:         ml,
	}
	bytes, err := yaml.Marshal(mPayload)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
func TestRASMutatorModelManifestSeralization(t *testing.T) {

	inputs := make([]Input, 5)
	inputs[0] = Input{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	inputs[1] = Input{
		Name:      "muncie.p04.tmp.hdf",
		Parameter: "ras p hdf file",
		Format:    "hdf",
	}
	inputs[2] = Input{
		Name:      "muncie.b04",
		Parameter: "ras b file",
		Format:    ".b**",
	}
	inputs[3] = Input{
		Name:      "muncie.prj",
		Parameter: "ras project file",
		Format:    ".prj",
	}
	inputs[4] = Input{
		Name:      "muncie.x04",
		Parameter: "ras x file",
		Format:    ".x**",
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
	mc := ModelConfiguration{
		Name: "Muncie",
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
	eventConfiguration := MockEventConfiguration()
	//someone has to make data somewhere... probably needs to be computed output
	prevModelOutput := make([]ComputedOutput, 5)
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
	prevModelOutput[1] = ComputedOutput{
		Name:         "Temp Project HDF File",
		Parameter:    "Project HDF File",
		Format:       ".hdf",
		ResourceInfo: paths[0],
	}
	prevModelOutput[2] = ComputedOutput{
		Name:         "RAS B file",
		Parameter:    "B file stuff",
		Format:       ".b**",
		ResourceInfo: paths[1],
	}
	prevModelOutput[3] = ComputedOutput{
		Name:         "RAS Project File",
		Parameter:    "Project Specification",
		Format:       ".prj",
		ResourceInfo: paths[2],
	}
	prevModelOutput[4] = ComputedOutput{
		Name:         "RAS X File",
		Parameter:    "X File stuff",
		Format:       ".x**",
		ResourceInfo: paths[3],
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
	mc := ModelConfiguration{
		Name: "Muncie",
	}
	prevModelOutput = append(prevModelOutput, eventConfiguration.ToInput())
	ml := ModelLinks{
		LinkedInputs:     prevModelOutput,
		NecessaryOutputs: outputs,
	}
	mPayload := ModelPayload{
		ModelConfiguration: mc,
		ModelLinks:         ml,
	}
	bytes, err := yaml.Marshal(mPayload)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
func TestRASRunnerModelManifestSeralization(t *testing.T) {
	inputs := make([]Input, 4)
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
	inputs[0] = Input{
		Name:      "Temp Project HDF File",
		Parameter: "Project HDF File",
		Format:    ".hdf",
	}
	inputs[1] = Input{
		Name:      "RAS B file",
		Parameter: "B file stuff",
		Format:    ".b**",
	}
	inputs[2] = Input{
		Name:      "RAS Project File",
		Parameter: "Project Specification",
		Format:    ".prj",
	}
	inputs[3] = Input{
		Name:      "RAS X File",
		Parameter: "X File stuff",
		Format:    ".x**",
	}
	mc := ModelConfiguration{
		Name: "Muncie",
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
	eventConfiguration := MockEventConfiguration()
	//someone has to make data somewhere... probably needs to be computed output
	prevModelOutput := make([]ComputedOutput, 4)
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
	prevModelOutput[0] = ComputedOutput{
		Name:         "Temp Project HDF File",
		Parameter:    "Project HDF File",
		Format:       ".hdf",
		ResourceInfo: paths[0],
	}
	prevModelOutput[1] = ComputedOutput{
		Name:         "RAS B file",
		Parameter:    "B file stuff",
		Format:       ".b**",
		ResourceInfo: paths[1],
	}
	prevModelOutput[2] = ComputedOutput{
		Name:         "RAS Project File",
		Parameter:    "Project Specification",
		Format:       ".prj",
		ResourceInfo: paths[2],
	}
	prevModelOutput[3] = ComputedOutput{
		Name:         "RAS X File",
		Parameter:    "X File stuff",
		Format:       ".x**",
		ResourceInfo: paths[3],
	}
	mc := ModelConfiguration{
		Name: "Muncie",
	}
	prevModelOutput = append(prevModelOutput, eventConfiguration.ToInput())
	ml := ModelLinks{
		LinkedInputs:     prevModelOutput,
		NecessaryOutputs: outputs,
	}
	mPayload := ModelPayload{
		ModelConfiguration: mc,
		ModelLinks:         ml,
	}
	bytes, err := yaml.Marshal(mPayload)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bytes))
}
