package model

//Plugin can compute a model, can be used to get model manifests
type Plugin struct {
	Name            string   `json:"plugin_name" yaml:"plugin_name"`
	ImageAndTag     string   `json:"plugin_image_and_tag" yaml:"plugin_image_and_tag"`
	CommandLineArgs []string `json:"commandline_args" yaml:"commandline_args"`
}
