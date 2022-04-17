package wat

//ModelConfiguration is a name and a path to a configuration
type ModelConfiguration struct {
	Name                    string   `json:"model_name" yaml:"model_name"`                               //model library guid?
	ModelConfigurationPaths []string `json:"model_configuration_paths" yaml:"model_configuration_paths"` //probably a uri?
}

//ModelManifest is defined by a set of files, provides inputs and ouptuts, is recognizable by a Model Library MCAT
type ModelManifest struct {
	ModelConfiguration `json:"model_configuration" yaml:"model_configuration"`
	Inputs             []Input  `json:"inputs" yaml:"inputs"`
	Outputs            []Output `json:"outputs" yaml:"outputs"`
}
type ModelLinks struct {
	LinkedInputs     []LinkedInput `json:"linked_inputs" yaml:"linked_inputs"`
	NecessaryOutputs []Output      `json:"required_outputs" yaml:"required_outputs"`
}
type ModelPayload struct {
	TargetPlugin       string `json:"target_plugin" yaml:"target_plugin"`
	ModelConfiguration `json:"model_configuration" yaml:"model_configuration"`
	ModelLinks         `json:"model_links" yaml:"model_links"`
	EventConfiguration `json:"event_config" yaml:"event_config"`
}
