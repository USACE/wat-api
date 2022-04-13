package wat

//Input define where a model needs input
type Input struct {
	Name      string `json:"name" yaml:"name"`
	Parameter string `json:"parameter" yaml:"parameter"`
	Format    string `json:"format" yaml:"format"`
}

//Output defines where a model can produce output the format, parameter and the link information
type Output struct {
	Name      string `json:"name" yaml:"name"`
	Parameter string `json:"parameter" yaml:"parameter"`
	Format    string `json:"format" yaml:"format"`
}

type LinkedInput struct {
	Input  `json:"input" yaml:"input"`
	Source Output `json:"source" yaml:"source"`
}

//@TODO: think more broadly about types of data, sources of data and how to differentiate them in a DAG
