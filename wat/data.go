package wat

//Input define where a model needs input
type Input struct {
	Name      string `json:"name"`
	Parameter string `json:"parameter"`
	Format    string `json:"format"`
}

//Output defines where a model can produce output the format, parameter and the link information
type Output struct {
	Name      string `json:"name"`
	Parameter string `json:"parameter"`
	Format    string `json:"format"`
}

type LinkedInput struct {
	Input  `json:"input"`
	Source Output `json:"source"`
}

//@TODO: think more broadly about types of data, sources of data and how to differentiate them in a DAG
