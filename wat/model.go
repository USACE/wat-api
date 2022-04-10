package wat

//Model is defined by a set of files, provides inputs and ouptuts, is recognizable by a Model Library MCAT
type Model struct {
	Name                   string   `json:"model_name"` //model library guid?
	ModelConfigurationPath string   `json:"model_configuration_path"`
	Inputs                 []Input  `json:"inputs"`
	Outputs                []Output `json:"outputs"`
}
type LinkedModel struct {
	Model            `json:"model"`
	LinkedInputs     []LinkedInput `json:"linked_inputs"`
	NecessaryOutputs []Output      `json:"required_outputs"`
}
type ModelManifest struct {
	TargetPlugin       string      `json:"target_plugin"`
	TargetModel        LinkedModel `json:"linked_model"`
	EventConfiguration `json:"event_config"`
}
