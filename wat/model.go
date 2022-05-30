package wat

import (
	"fmt"
	"io/ioutil"

	"github.com/USACE/filestore"
	"gopkg.in/yaml.v2"
)

//ModelConfiguration is a model name and an optional model alternative
type ModelConfiguration struct {
	Name        string `json:"model_name" yaml:"model_name"`
	Alternative string `json:"model_alternative,omitempty" yaml:"model_alternative,omitempty"` //model library guid?
	//ModelConfigurationResources []ResourceInfo `json:"model_configuration_paths" yaml:"model_configuration_paths"` //probably a uri?
}
type ModelComputeResources struct {
	MinCpus       *int64    `json:"min_cpus" yaml:"min_cpus"`
	DesiredCpus   *int64    `json:"desired_cpus" yaml:"desired_cpus"`
	MaxCpus       *int64    `json:"max_cpus" yaml:"max_cpus"`
	InstanceTypes []*string `json:"instance_types" yaml:"instance_types"`
	Type          *string   `json:"compute_environment_type" yaml:"compute_environment_type"`
	Managed       bool      `json:"compute_environment_management_state" yaml:"compute_environment_management_state"`
}

//ModelManifest is defined by a set of files, provides inputs and ouptuts, is recognizable by a Model Library MCAT
type ModelManifest struct {
	//Batch or Lambda
	//TaskType              string `json:"task_type" yaml:"task_type"`
	Plugin                `json:"plugin" yaml:"plugin"`
	ModelConfiguration    `json:"model_configuration" yaml:"model_configuration"`
	ModelComputeResources `json:"model_compute_resources" yaml:"model_compute_resources"`
	Inputs                []DataDescription `json:"inputs" yaml:"inputs"`
	Outputs               []DataDescription `json:"outputs" yaml:"outputs"`
}
type ModelLinks struct {
	LinkedInputs     []LinkedDataDescription `json:"linked_inputs" yaml:"linked_inputs"`
	NecessaryOutputs []LinkedDataDescription `json:"required_outputs" yaml:"required_outputs"`
}
type ModelPayload struct {
	//Plugin       Plugin `json:"target_plugin" yaml:"target_plugin"`
	ModelConfiguration `json:"model_configuration" yaml:"model_configuration"`
	ModelLinks         `json:"model_links" yaml:"model_links"`
	//EventConfiguration `json:"event_config" yaml:"event_config"`
}

func (mp ModelPayload) EventConfiguration() EventConfiguration {
	return EventConfiguration{} //look through model links to find an input that is an event configuration...
}
func (mp *ModelPayload) SetEventConfiguration(ec EventConfiguration) {
	//look through model links to find an input that is an event configuration... and set it!
}

// LoadModelPayload
func LoadModelPayloadFromS3(payloadFile string, fs filestore.FileStore) (ModelPayload, error) {
	var p ModelPayload
	fmt.Println("reading payload:", payloadFile)
	data, err := fs.GetObject(payloadFile)
	if err != nil {
		return p, err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return p, err
	}
	//fmt.Println(string(body))
	err = yaml.Unmarshal(body, &p)
	if err != nil {
		return p, err
	}
	//fmt.Println(p)
	return p, nil
}
