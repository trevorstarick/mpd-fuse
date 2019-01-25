package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)

// todo: provide basic fs
// todo: provide access to remote machine
// todo: implement basic caching
// todo: implement limited caching (512MB)
// todo: implement least-accessed pruning

const CRAWL_ENABLED = false

type Entry struct {
	Name     string
	DType    fuse.DirentType
	Type     string // directory, file
	MTime    string
	Size     uint
	Children []fuse.Dirent
}

var Tree = make(map[string]Entry)

var ROOT = "http://azure.shivver.io/media/music"
var CACHE = "cache"

var crawlChan = make(chan string, 8192)

var mutex sync.Mutex

func RequestRoute(route string) {
	u := ROOT + route + "/"
	res, err := http.Get(u)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(res.Body)

	entries := []Entry{}
	decoder.Decode(&entries)

	var children = make([]fuse.Dirent, len(entries))

	for i, entry := range entries {
		children[i] = fuse.Dirent{
			Inode: 2,
			Name:  entry.Name,
		}

		if entry.Type == "directory" {
			// recursively index directories
			if CRAWL_ENABLED {
				crawlChan <- (route + "/" + url.PathEscape(entry.Name))
			}

			entry.DType = fuse.DT_Dir
			children[i].Type = fuse.DT_Dir
		} else if entry.Type == "file" {
			entry.DType = fuse.DT_File
			children[i].Type = fuse.DT_File
		} else {
			entry.DType = fuse.DT_File
			children[i].Type = fuse.DT_Unknown
		}

		mutex.Lock()
		Tree[route+"/"+url.PathEscape(entry.Name)] = entry
		mutex.Unlock()
	}

	e := Tree[route]
	e.Children = children
	mutex.Lock()
	Tree[route] = e
	mutex.Unlock()
}

func main() {
	// used for concurrent indexing
	if CRAWL_ENABLED {
		for i := 0; i < 64; i++ {
			go func() {
				for u := range crawlChan {
					RequestRoute(u)
				}
			}()
		}
	}

	mountpoint := "mountpoint"
	RequestRoute("")

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
