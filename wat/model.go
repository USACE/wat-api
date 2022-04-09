package wat

//Model is defined by a set of files, provides inputs and ouptuts, is recognizable by a Model Library MCAT
type ModelManifest struct {
	Name                   string   `json:"model_name"` //model library guid?
	ModelConfigurationPath string   `json:"model_configuration_path"`
	Inputs                 []Input  `json:"inputs"`
	Outputs                []Output `json:"outputs"`
	EventConfiguration     `json:"event_config"`
}
