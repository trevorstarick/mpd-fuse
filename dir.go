package main

import (
	"net/url"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Dir struct {
	Filename string
	Location string
}

var _ fs.Node = (*Dir)(nil)

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0755

	return nil
}

var _ fs.NodeStringLookuper = (*Dir)(nil)

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	filename := d.Filename + "/" + name
	location := d.Location + "/" + url.PathEscape(name)

	if content, exists := Tree[location]; exists {
		switch content.DType {
		case fuse.DT_File:
			return &File{
				Filename: filename,
				Location: location,
				Size:     uint64(content.Size),
			}, nil
		case fuse.DT_Dir:
			return &Dir{
				Filename: filename,
				Location: location,
			}, nil
		default:
			break
		}
	}

	return nil, fuse.ENOENT
}

var _ fs.HandleReadDirAller = (*Dir)(nil)

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	if content, exists := Tree[d.Location]; !exists || len(content.Children) == 0 {
		RequestRoute(d.Location)
	}

	return Tree[d.Location].Children, nil
}
