package main

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Dir struct {
	location string
}

var _ fs.Node = (*Dir)(nil)

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0755

	return nil
}

var _ fs.NodeStringLookuper = (*Dir)(nil)

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	name = d.location + "/" + name

	if content, exists := Tree[name]; exists {
		switch content.DType {
		case fuse.DT_File:
			return &File{
				location: name,
				Size:     uint64(content.Size),
			}, nil
		case fuse.DT_Dir:
			return &Dir{
				location: name,
			}, nil
		default:
			break
		}
	}

	return nil, fuse.ENOENT
}

var _ fs.HandleReadDirAller = (*Dir)(nil)

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	if content, exists := Tree[d.location]; !exists || len(content.Children) == 0 {
		RequestRoute(d.location)
	}

	return Tree[d.location].Children, nil
}
