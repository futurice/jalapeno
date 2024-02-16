package recipe

import (
	"crypto/sha256"
	"fmt"
)

const (
	HashPrefix = "sha256:"
)

type File struct {
	Checksum string `yaml:"checksum"` // e.g. "sha256:xxxxxxxxx" w. default algo
	Content  []byte `yaml:"-"`
}

func NewFile(content []byte) File {
	f := File{Content: content}
	f.Checksum = f.hash()

	return f
}

func (f File) HasBeenModifiedByUser() bool {
	return f.Checksum != f.hash()
}

func (f File) hash() string {
	return fmt.Sprintf("%s%x", HashPrefix, sha256.Sum256(f.Content))
}
