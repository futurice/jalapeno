package recipe

type Recipe struct {
	Metadata  *Metadata
	Variables []*Variable
	Templates []*File
}
