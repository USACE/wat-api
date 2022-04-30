package wat

//ModelConfiguration is a name and a path to a configuration
type ModelConfiguration struct {
	Name                        string         `json:"model_name" yaml:"model_name"`                               //model library guid?
	ModelConfigurationResources []ResourceInfo `json:"model_configuration_paths" yaml:"model_configuration_paths"` //probably a uri?
}

//ModelManifest is defined by a set of files, provides inputs and ouptuts, is recognizable by a Model Library MCAT
type ModelManifest struct {
	ModelConfiguration `json:"model_configuration" yaml:"model_configuration"`
	Inputs             []RequiredInput  `json:"inputs" yaml:"inputs"`
	Outputs            []PossibleOutput `json:"outputs" yaml:"outputs"`
}
type ModelLinks struct {
	LinkedInputs     []ComputedOutput `json:"linked_inputs" yaml:"linked_inputs"`
	NecessaryOutputs []PossibleOutput `json:"required_outputs" yaml:"required_outputs"`
}
type ModelPayload struct {
	TargetPlugin       string `json:"target_plugin" yaml:"target_plugin"`
	PluginImageAndTag  string `json:"plugin_image_and_tag" yaml:"plugin_image_and_tag"`
	ModelConfiguration `json:"model_configuration" yaml:"model_configuration"`
	ModelLinks         `json:"model_links" yaml:"model_links"`
	EventConfiguration `json:"event_config" yaml:"event_config"`
}
