package main

import (
	"encoding/json"
	"net/http"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"github.com/davecgh/go-spew/spew"
)

// todo: provide basic fs
// todo: provide access to remote machine
// todo: implement basic caching
// todo: implement limited caching (512MB)
// todo: implement least-accessed pruning

type HTTPEntry struct {
	Name  string
	Type  string // directory, file
	MTime string
	Size  uint
}

var DirCache = make(map[string][]HTTPEntry)
var FileCache = make(map[string][]HTTPEntry)

var ROOT = "http://azure.shivver.io/music"

func RequestRoute(route string) {
	res, err := http.Get(ROOT + route)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(res.Body)

	entries := []HTTPEntry{}
	decoder.Decode(&entries)

	for _, entry := range entries {
		if entry.Type == "directory" {
			DirCache[route] = append(DirCache[route], entry)
			DirCache[route+"/"+entry.Name] = []HTTPEntry{}
		} else if entry.Type == "file" {
			FileCache[route] = append(FileCache[route], entry)
			FileCache[route+"/"+entry.Name] = []HTTPEntry{}
		} else {
			spew.Dump(entry)
		}
	}
}

func main() {
	mountpoint := "mountpoint"
	RequestRoute("/")

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("mpd"),
		fuse.Subtype("mpdfs"),
		fuse.LocalVolume(),
		fuse.VolumeName("MPD"),
	)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	srv := fs.New(c, nil)
	filesys := &FS{}

	if err := srv.Serve(filesys); err != nil {
		panic(err)
	}

	<-c.Ready
	if err := c.MountError; err != nil {
		panic(err)
	}
}
