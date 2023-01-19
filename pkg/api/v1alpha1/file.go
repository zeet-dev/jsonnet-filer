package v1alpha1

const (
	ApiVersion string = "jsonnet-filer.zeet.co/v1alpha1"
	Kind       string = "File"
)

type ObjectMeta struct {
	Name string `json:"name"`
}

type File struct {
	ApiVersion       string     `json:"apiVersion"`
	Kind             string     `json:"kind"`
	Metadata         ObjectMeta `json:"metadata"`
	Content          any        `json:"content"`
	EncodingStrategy string     `json:"encodingStrategy"`
}

func NewFile(name string, content any) File {
	return File{
		ApiVersion: ApiVersion,
		Kind:       Kind,
		Metadata: ObjectMeta{
			Name: name,
		},
		Content:          content,
		EncodingStrategy: "yaml",
	}
}

func (f File) Ok() bool {
	if f.ApiVersion != ApiVersion {
		return false
	}

	if f.Kind != Kind {
		return false
	}

	return true
}
