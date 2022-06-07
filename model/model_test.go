package model

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
	inputs := make([]DataDescription, 1)
	inputs[0] = DataDescription{
		Name:      "Project File",
		Parameter: "Project Specification",
		Format:    ".json",
	}
	outputs := make([]DataDescription, 1)
	outputs[0] = DataDescription{
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

func TestHSMModelManifestSeralization(t *testing.T) {

	inputs := make([]DataDescription, 1)
	inputs[0] = DataDescription{
		Name:      "Project File",
		Parameter: "Project Specification",
		Format:    ".json",
	}
	outputs := make([]DataDescription, 3)
	outputs[0] = DataDescription{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	outputs[1] = DataDescription{
		Name:      "hsm2.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	outputs[2] = DataDescription{
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
	prevModelOutput := make([]LinkedDataDescription, 2)
	prevModelOutput[0] = LinkedDataDescription{
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
	outputs := make([]LinkedDataDescription, 1)
	outputs[0] = LinkedDataDescription{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
		ResourceInfo: ResourceInfo{
			Scheme:    "https",
			Authority: "/minio/runs/realization_1/event_1",
			Fragment:  "hsm.json",
		},
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

	inputs := make([]DataDescription, 5)
	inputs[0] = DataDescription{
		Name:      "hsm1.csv",
		Parameter: "flow",
		Format:    ".csv",
	}
	inputs[1] = DataDescription{
		Name:      "muncie.p04.tmp.hdf",
		Parameter: "ras p hdf file",
		Format:    "hdf",
	}
	inputs[2] = DataDescription{
		Name:      "muncie.b04",
		Parameter: "ras b file",
		Format:    ".b**",
	}
	inputs[3] = DataDescription{
		Name:      "muncie.prj",
		Parameter: "ras project file",
		Format:    ".prj",
	}
	inputs[4] = DataDescription{
		Name:      "muncie.x04",
		Parameter: "ras x file",
		Format:    ".x**",
	}
	outputs := make([]DataDescription, 4)
	outputs[0] = DataDescription{
		Name:      "muncie.p04.tmp.hdf",
		Parameter: "ras p hdf file",
		Format:    "hdf",
	}
	outputs[1] = DataDescription{
		Name:      "muncie.b04",
		Parameter: "ras b file",
		Format:    ".b**",
	}
	outputs[2] = DataDescription{
		Name:      "muncie.prj",
		Parameter: "ras project file",
		Format:    ".prj",
	}
	outputs[3] = DataDescription{
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
	prevModelOutput := make([]LinkedDataDescription, 5)
	prevModelOutput[0] = LinkedDataDescription{
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
	prevModelOutput[1] = LinkedDataDescription{
		Name:         "Temp Project HDF File",
		Parameter:    "Project HDF File",
		Format:       ".hdf",
		ResourceInfo: paths[0],
	}
	prevModelOutput[2] = LinkedDataDescription{
		Name:         "RAS B file",
		Parameter:    "B file stuff",
		Format:       ".b**",
		ResourceInfo: paths[1],
	}
	prevModelOutput[3] = LinkedDataDescription{
		Name:         "RAS Project File",
		Parameter:    "Project Specification",
		Format:       ".prj",
		ResourceInfo: paths[2],
	}
	prevModelOutput[4] = LinkedDataDescription{
		Name:         "RAS X File",
		Parameter:    "X File stuff",
		Format:       ".x**",
		ResourceInfo: paths[3],
	}
	outputs := make([]LinkedDataDescription, 4)
	outputs[0] = LinkedDataDescription{
		Name:      "muncie.p04.tmp.hdf",
		Parameter: "ras p hdf file",
		Format:    "hdf",
		ResourceInfo: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
			Fragment:  "muncie.p04.tmp.hdf",
		},
	}
	outputs[1] = LinkedDataDescription{
		Name:      "muncie.b04",
		Parameter: "ras b file",
		Format:    ".b**",
		ResourceInfo: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
			Fragment:  "muncie.b04",
		},
	}
	outputs[2] = LinkedDataDescription{
		Name:      "muncie.prj",
		Parameter: "ras project file",
		Format:    ".prj",
		ResourceInfo: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
			Fragment:  "muncie.prj",
		},
	}
	outputs[3] = LinkedDataDescription{
		Name:      "muncie.x04",
		Parameter: "ras x file",
		Format:    ".x**",
		ResourceInfo: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
			Fragment:  "muncie.x04",
		},
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
	inputs := make([]DataDescription, 4)
	outputs := make([]DataDescription, 2)
	outputs[0] = DataDescription{
		Name:      "muncie.p04.hdf",
		Parameter: "ras results hdf file",
		Format:    ".hdf",
	}
	outputs[1] = DataDescription{
		Name:      "muncie.log",
		Parameter: "ras log file",
		Format:    ".log",
	}
	inputs[0] = DataDescription{
		Name:      "Temp Project HDF File",
		Parameter: "Project HDF File",
		Format:    ".hdf",
	}
	inputs[1] = DataDescription{
		Name:      "RAS B file",
		Parameter: "B file stuff",
		Format:    ".b**",
	}
	inputs[2] = DataDescription{
		Name:      "RAS Project File",
		Parameter: "Project Specification",
		Format:    ".prj",
	}
	inputs[3] = DataDescription{
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
	prevModelOutput := make([]LinkedDataDescription, 4)
	outputs := make([]LinkedDataDescription, 2)
	outputs[0] = LinkedDataDescription{
		Name:      "muncie.p04.hdf",
		Parameter: "ras results hdf file",
		Format:    ".hdf",
		ResourceInfo: ResourceInfo{
			Scheme:    "https",
			Authority: "/minio/runs/realization_1/event_1/",
			Fragment:  "muncie.p04.hdf",
		},
	}
	outputs[1] = LinkedDataDescription{
		Name:      "muncie.log",
		Parameter: "ras log file",
		Format:    ".log",
		ResourceInfo: ResourceInfo{
			Scheme:    "https",
			Authority: "/minio/runs/realization_1/event_1",
			Fragment:  "muncie.log",
		},
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
	prevModelOutput[0] = LinkedDataDescription{
		Name:         "Temp Project HDF File",
		Parameter:    "Project HDF File",
		Format:       ".hdf",
		ResourceInfo: paths[0],
	}
	prevModelOutput[1] = LinkedDataDescription{
		Name:         "RAS B file",
		Parameter:    "B file stuff",
		Format:       ".b**",
		ResourceInfo: paths[1],
	}
	prevModelOutput[2] = LinkedDataDescription{
		Name:         "RAS Project File",
		Parameter:    "Project Specification",
		Format:       ".prj",
		ResourceInfo: paths[2],
	}
	prevModelOutput[3] = LinkedDataDescription{
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
