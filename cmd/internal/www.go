package internal

import (
	"go/token"
)

type ImportMapEntry struct {
	Key  string
	Path string
}

type FileImports map[string]ImportMapEntry

type FolderContextInformation struct {
	Position token.Pos   `json:"pos,omitempty"`
	FileName string      `json:"file_name,omitempty"`
	AbsPath  string      `json:"abs_path,omitempty"`
	RelPath  string      `json:"rel_path,omitempty"`
	Imports  FileImports `json:"imports,omitempty"`
}

func (fci *FolderContextInformation) Pos() token.Pos {
	return fci.Position
}

func (fci *FolderContextInformation) End() token.Pos {
	return fci.Position
}
