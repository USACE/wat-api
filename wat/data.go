package wat

//Input define where a model needs input
type Input struct {
	Name      string `json:"name" yaml:"name"`
	Parameter string `json:"parameter" yaml:"parameter"`
	Format    string `json:"format" yaml:"format"`
}
type ResourceInfo struct {
	Type         string
	ResourcePath string
	//https://pkg.go.dev/go.lsp.dev/uri  consider this.
}

//Output defines where a model can produce output the format, parameter and the link information
type Output struct {
	Name           string             `json:"name" yaml:"name"`
	ProducingModel ModelConfiguration `json:"producing_model,omitempty" yaml:"producing_model,omitempty"`
	Parameter      string             `json:"parameter" yaml:"parameter"`
	Format         string             `json:"format" yaml:"format"`
	//ResourceInfo //where do we put the data, and how do we get to it?
}

type LinkedInput struct {
	Input  `json:"input" yaml:"input"`
	Source Output `json:"source" yaml:"source"`
}

//@TODO: think more broadly about types of data, sources of data and how to differentiate them in a DAG
