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
	dirPath := path.Dir(CACHE + f.location)
	os.MkdirAll(dirPath, 0777)

	out, err := os.Create(CACHE + f.location + ".tmp")
	if err != nil {
		panic(err)
	}

	defer out.Close()

	res, err := http.Get(ROOT + f.location)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		panic(err)
	}

	os.Rename(CACHE+f.location+".tmp", CACHE+f.location)

	f.File, _ = os.Open(CACHE + f.location)

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

	f.File.Seek(req.Offset, io.SeekStart)

	_, err := (*f.File).Read(data)
	if err != nil && err != io.EOF {
		panic(err)
		return err
	}

	resp.Data = data
	return nil
}
