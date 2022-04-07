package wat

type Link struct {
	Input  `json:"input"`
	Output `json:"output"`
}

//InputDataLocations define where a model will get input from the format, parameter and the link information
type Input struct {
	Name      string `json:"name"`
	Parameter string `json:"parameter"`
	Format    string `json:"format"`
}

//OutputDataLocations define where a model can produce output the format, parameter and the link information
type Output struct {
	Name                 string   `json:"name"`
	Parameter            string   `json:"parameter"`
	Format               string   `json:"format"`
	LinkInfo             LinkInfo `json:"link_info"`
	GeneratingPluginName string   `json:"generating_plugin_name"`
	GeneratingModelName  string   `json:"generating_model_name"`
}

type LinkInfo interface {
	Type() LinkType
}

type LinkType byte

const (
	CloudStore LinkType = 0
	LocalStore LinkType = 1
	Virtual    LinkType = 2
)

type LocalCSVLink struct {
	Path string `json:"Path"`
}

func (lcsv LocalCSVLink) Type() LinkType {
	return LocalStore
}

type LocalDSSTimeSeriesLink struct {
	DSSAPart    string `json:"DSSAPart"`
	DSSBPart    string `json:"DSSBPart"`
	DSSCPart    string `json:"DSSCPart"`
	DSSFilePath string `json:"DSSFilePath"`
}

func (ldssTS LocalDSSTimeSeriesLink) Type() LinkType {
	return LocalStore
}
