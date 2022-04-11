package wat

//ModelConfiguration is a name and a path to a configuration
type ModelConfiguration struct {
	Name                   string `json:"model_name"`               //model library guid?
	ModelConfigurationPath string `json:"model_configuration_path"` //probably a uri?
}

//ModelManifest is defined by a set of files, provides inputs and ouptuts, is recognizable by a Model Library MCAT
type ModelManifest struct {
	ModelConfiguration `json:"model_configuration"`
	Inputs             []Input  `json:"inputs"`
	Outputs            []Output `json:"outputs"`
}
type ModelLinks struct {
	LinkedInputs     []LinkedInput `json:"linked_inputs"`
	NecessaryOutputs []Output      `json:"required_outputs"`
}
type ModelPayload struct {
	TargetPlugin       string `json:"target_plugin"`
	ModelConfiguration `json:"model_configuration"`
	ModelLinks         `json:"model_links"`
	EventConfiguration `json:"event_config"`
}
