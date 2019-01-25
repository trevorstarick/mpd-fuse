package main

import (
	"bazil.org/fuse/fs"
)

// FS implements the hello world file system.
type FS struct {
	Filename string
	Location string
}

var _ fs.FS = (*FS)(nil)

func (f *FS) Root() (fs.Node, error) {
	return &Dir{
		Filename: f.Filename,
		Location: "",
	}, nil
}
