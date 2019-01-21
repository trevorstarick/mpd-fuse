# mpd-fuse

an extension of [mpd-cache](https://github.com/trevorstarick/mpd-cache) that replaces the sshfs mechanism in exchange for a 100% Go virtual filesystem backed by a remote disk. the goals are to create a FUSE filesystem to let the user connect to a remote server, mount a directory from the remote server to a local mounting point. extra features that are interesting to consider is the ability to allow for content caching to support less than reliable connections, unique folders that allow for bi-directional file changes.
  
reading material:  
[libfuse/fuse](https://github.com/libfuse/libfuse)  
[@wiki/FUSE](https://en.wikipedia.org/wiki/Filesystem_in_Userspace)  
[bazil/fuse](https://github.com/bazil/fuse)  
