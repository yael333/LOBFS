package main

import (
	"os"
	"fmt"
	
	"github.com/konata-chan404/LOBFS/babel"
	"github.com/winfsp/cgofuse/fuse"
)

const (
	filename = "hello"
	contents = "hello, world\n"
)

type Hellofs struct {
	fuse.FileSystemBase
}

func (self *Hellofs) Open(path string, flags int) (errc int, fh uint64) {
	switch path {
	case "/" + filename:
		return 0, 0
	default:
		return -fuse.ENOENT, ^uint64(0)
	}
}

func (self *Hellofs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	switch path {
	case "/":
		stat.Mode = fuse.S_IFDIR | 0555
		return 0
	case "/" + filename:
		stat.Mode = fuse.S_IFREG | 0444
		stat.Size = int64(len(contents))
		return 0
	default:
		return -fuse.ENOENT
	}
}

func (self *Hellofs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	endofst := ofst + int64(len(buff))
	if endofst > int64(len(contents)) {
		endofst = int64(len(contents))
	}
	if endofst < ofst {
		return 0
	}
	n = copy(buff, contents[ofst:endofst])
	return
}

func (self *Hellofs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	fill(filename, nil, 0)
	return 0
}

func main() {
	lib := make(babel.Library)
	loc := babel.NewLocation(1, 2, 3, 4, 5)
	page, err := lib.GetPageAtLocation(loc)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(page)
	
	vol := lib.GetHexagon(1).GetWall(2).GetShelf(3).GetVolume(loc)
	fmt.Println("Volume name:", vol.Name)

	hellofs := &Hellofs{}
	host := fuse.NewFileSystemHost(hellofs)
	host.Mount("", os.Args[1:])
}