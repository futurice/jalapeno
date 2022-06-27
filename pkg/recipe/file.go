package recipe

type File struct {
	// Name is the path-like name of the template.
	Name string `json:"name"`
	// Data is the template as byte data.
	Data []byte `json:"data"`
}
