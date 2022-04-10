package wat

type ModelConfiguration struct {
	Name                   string `json:"model_name"` //model library guid?
	ModelConfigurationPath string `json:"model_configuration_path"`
}

//Model is defined by a set of files, provides inputs and ouptuts, is recognizable by a Model Library MCAT
type Model struct {
	ModelConfiguration `json:"model_configuration"`
	Inputs             []Input  `json:"inputs"`
	Outputs            []Output `json:"outputs"`
}
type ModelLinks struct {
	LinkedInputs     []LinkedInput `json:"linked_inputs"`
	NecessaryOutputs []Output      `json:"required_outputs"`
}
type ModelManifest struct {
	TargetPlugin       string `json:"target_plugin"`
	ModelConfiguration `json:"model_configuration"`
	ModelLinks         `json:"model_links"`
	EventConfiguration `json:"event_config"`
}
