package wat

//Input define where a model needs input
type RequiredInput struct {
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
type PossibleOutput struct {
	Name           string             `json:"name" yaml:"name"`
	ProducingModel ModelConfiguration `json:"producing_model,omitempty" yaml:"producing_model,omitempty"`
	Parameter      string             `json:"parameter" yaml:"parameter"`
	Format         string             `json:"format" yaml:"format"`
}
type ComputedOutput struct {
	Name           string             `json:"name" yaml:"name"`
	ProducingModel ModelConfiguration `json:"producing_model,omitempty" yaml:"producing_model,omitempty"`
	Parameter      string             `json:"parameter" yaml:"parameter"`
	Format         string             `json:"format" yaml:"format"`
	ResourceInfo                      //uri?
}
type SatisfiedLink struct {
	RequiredInput `json:"input" yaml:"input"`
	Source        ComputedOutput `json:"source" yaml:"source"`
}

//@TODO: think more broadly about types of data, sources of data and how to differentiate them in a DAG
