package main

import (
	"fmt"
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
	fmt.Println("lookup", d.location, name)
	name = d.location + "/" + name

	if _, exists := DirCache[name]; exists {
		return &Dir{
			location: name,
		}, nil
	} else if _, exists := FileCache[name]; exists {
		return &File{
			location: name,
		}, nil
	}

	return nil, fuse.ENOENT
}

var _ fs.HandleReadDirAller = (*Dir)(nil)

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	fmt.Println("readdirall", d.location)
	var entries = []fuse.Dirent{}

	if content, exists := DirCache[d.location]; !exists || len(content) == 0 {
		RequestRoute(d.location)
	}

	for _, entry := range DirCache[d.location] {
		e := fuse.Dirent{
			Inode: 2,
			Name:  entry.Name,
			Type:  fuse.DT_Dir,
		}
		entries = append(entries, e)
	}

	for _, entry := range FileCache[d.location] {
		e := fuse.Dirent{
			Inode: 2,
			Name:  entry.Name,
			Type:  fuse.DT_File,
		}
		entries = append(entries, e)
	}

	return entries, nil
}
