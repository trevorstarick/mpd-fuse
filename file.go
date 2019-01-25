package main

import (
	"io"
	"net/http"
	"os"
	"path"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type File struct {
	fs.NodeRef
	fuse     *fs.Server
	Filename string
	Location string
	Size     uint64
	Body     *io.ReadCloser
	File     *os.File
}

func (f *File) Download() {
	res, err := http.Get(ROOT + f.Location)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	_, err = io.Copy(f.File, res.Body)
	if err != nil {
		panic(err)
	}
}

var _ fs.Node = (*File)(nil)

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = f.Size
	return nil
}

var _ fs.NodeOpener = (*File)(nil)

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	// todo: check if exists in cache
	// todo: stream from remote

	if _, err := os.Stat(f.Filename); os.IsNotExist(err) {
		dirPath := path.Dir(f.Filename)
		os.MkdirAll(dirPath, 0777)

		var err error
		f.File, err = os.Create(f.Filename)
		if err != nil {
			panic(err)
		}

		f.Download()

	} else {
		if f.File == nil {
			f.File, err = os.Open(f.Filename)
			if err != nil {
				panic(err)
			}
		}
	}

	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	resp.Flags |= fuse.OpenKeepCache
	return f, nil
}

var _ fs.Handle = (*File)(nil)

var _ fs.HandleReader = (*File)(nil)

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	data := make([]byte, req.Size)

	_, err := f.File.Seek(req.Offset, io.SeekStart)
	if err != nil && err != io.EOF {
		panic(err)
	}

	_, err = f.File.Read(data)
	if err != nil && err != io.EOF {
		panic(err)
	}

	resp.Data = data
	return nil
}
