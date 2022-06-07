package model

//Input define where a model needs input
type DataDescription struct {
	//maybe file name? maybe location name?
	//this struct still needs discussion and consideration.
	Name      string `json:"name" yaml:"name"`
	Parameter string `json:"parameter" yaml:"parameter"`
	Format    string `json:"format" yaml:"format"`
}

type LinkedDataDescription struct {
	Name         string `json:"name" yaml:"name"`
	Parameter    string `json:"parameter,omitempty" yaml:"parameter,omitempty"`
	Format       string `json:"format" yaml:"format"`
	ResourceInfo `json:"resource_info" yaml:"resource_info"`
}
type ResourceInfo struct {
	//Might be S3, Redis, SQS, Azure, Local
	Scheme string `json:"scheme" yaml:"scheme"` //http or https for example
	//Bucket Name, service address, queue, bucket, root directory
	Authority string `json:"authority" yaml:"authority"` // //minio:9001 for example
	//path from bucket, key, omit for sqs?, path from bucket, path from root
	Path     string `json:"path,omitempty" yaml:"path,omitempty"`   //omit empty default value "/"
	Query    string `json:"query,omitempty" yaml:"query,omitempty"` //omit empty
	Fragment string `json:"fragment,omitempty" yaml:"fragment,omitempty"`
	//https://pkg.go.dev/go.lsp.dev/uri  consider this.
	/*
			    foo://example.com:8042/over/there?name=ferret#nose
		         \_/   \______________/\_________/ \_________/ \__/
		          |           |            |            |        |
		       scheme     authority       path        query   fragment
		          |   _____________________|__
		         / \ /                        \
		         urn:example:animal:ferret:nose
	*/
}

//@TODO: think more broadly about types of data, sources of data and how to differentiate them in a DAG
