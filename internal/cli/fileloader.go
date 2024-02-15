package cli

import (
	"os"
	"path/filepath"

	"github.com/jasoncabot/fabled-story-book/internal/jabl"
)

type fileLoader struct {
	base string
}

func NewFileLoader(base string) *fileLoader {
	return &fileLoader{
		base: base,
	}
}

func (d *fileLoader) LoadSection(identifier jabl.SectionId, onLoad func(code string, err error)) {
	if filepath.IsAbs(string(identifier)) || filepath.Clean(string(identifier)) != string(identifier) {
		onLoad("", os.ErrInvalid)
		return
	}
	contents, err := os.ReadFile(filepath.Join(d.base, string(identifier)))
	if err == nil {
		onLoad(string(contents), nil)
	} else {
		onLoad("", err)
	}
}
