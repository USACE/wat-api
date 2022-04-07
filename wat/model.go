package wat

//Model is defined by a set of files, provides inputs and ouptuts, is recognizable by a Model Library MCAT
type Model struct {
	Name    string   `json:"model_name"` //model library guid?
	Inputs  []string `json:"inputs"`
	Outputs []string `json:"outputs"`
}
