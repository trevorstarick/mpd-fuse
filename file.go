package main

import (
	"fmt"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"golang.org/x/net/context"
)

type File struct {
	fs.NodeRef
	location string
	fuse     *fs.Server
	content  []byte
	count    uint64
}

var _ fs.Node = (*File)(nil)

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	fmt.Println("attr", f.location)
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len([]byte("test\n")))
	return nil
}

var _ fs.NodeOpener = (*File)(nil)

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	fmt.Println("open", f.location)
	// todo: check if exists in cache
	// todo: stream from remote
	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	resp.Flags |= fuse.OpenKeepCache
	return f, nil
}

var _ fs.Handle = (*File)(nil)

var _ fs.HandleReader = (*File)(nil)

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	fmt.Println("read", f.location)
	fuseutil.HandleRead(req, resp, []byte("test\n"))
	return nil
}
