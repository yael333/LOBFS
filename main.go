package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/winfsp/cgofuse/fuse"
	"github.com/konata-chan404/LOBFS/babel"
)

type BabelFS struct {
	fuse.FileSystemBase
}

// find retrieves the address from the given path and returns a boolean indicating if it's a file.
func (self *BabelFS) find(path string) (address babel.Address, isFile bool) {
	parts := strings.Split(path, "/")

	if len(parts) >= 6 {
		// Extract the address components from the path
		hex := parts[1]
		wall, _ := strconv.ParseUint(parts[2], 10, 32)
		shelf, _ := strconv.ParseUint(parts[3], 10, 32)
		volume, _ := strconv.ParseUint(parts[4], 10, 32)
		page, _ := strconv.ParseUint(parts[5], 10, 32)

		// Create and return the address
		return babel.Address{Hex: hex, Wall: uint32(wall - 1), Shelf: uint32(shelf - 1), Volume: uint32(volume - 1), Page: uint32(page - 1)}, true
	}

	log.Printf("[</3] Path not found: %s\n", path)
	return babel.Address{}, false
}

// Open tries to open a file with the given path and flags. It returns an error code and a file handle.
func (self *BabelFS) Open(path string, flags int) (errc int, fh uint64) {
	_, isFile := self.find(path)

	if isFile {
		return 0, 0
	}

	log.Printf("[</3] File not found: %s\n", path)
	return -fuse.ENOENT, ^uint64(0)
}

// Getattr retrieves file attributes for a given path and file handle. It returns an error code.
func (self *BabelFS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	// Count the number of slashes in the path
	switch strings.Count(path, "/") {
	case 0:
		*stat = fuse.Stat_t{Mode: fuse.S_IFDIR | 0555} // Root directory
	case 1:
		*stat = fuse.Stat_t{Mode: fuse.S_IFDIR | 0555} // Hex directory
	case 2:
		*stat = fuse.Stat_t{Mode: fuse.S_IFDIR | 0555} // Wall directory
	case 3:
		*stat = fuse.Stat_t{Mode: fuse.S_IFDIR | 0555} // Shelf directory
	case 4:
		*stat = fuse.Stat_t{Mode: fuse.S_IFDIR | 0555} // Volume directory
	case 5:
		*stat = fuse.Stat_t{Mode: fuse.S_IFREG | 0444, Size: int64(babel.PAGE_LENGTH)} // Page file
	default:
		log.Printf("[</3] Can't find Attributes for Invalid path: %s\n", path)
		return -fuse.ENOENT
	}
	return 0
}

// Read reads data from a file into a buffer. It returns the number of bytes read.
func (self *BabelFS) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	address, isFile := self.find(path)

	if isFile {
		page, err := babel.GeneratePage(address)
		if err != nil {
			return 0
		}

		pageStr := string(page)
		end := ofst + int64(len(buff))
		if end > int64(len(pageStr)) {
			end = int64(len(pageStr))
		}
		if ofst < end {
			copy(buff, pageStr[ofst:end])
			log.Printf("[<3] Successfully read from path: %s\n", path)
			return int(end - ofst)
		}
	}
	else {
		log.Printf("[</3] Failed to read from path: %s\n", path)
	}
	return 0
}

// Readdir reads the contents of a directory. It returns an error code.
func (self *BabelFS) Readdir(path string, fill func(name string, stat *fuse.Stat_t, ofst int64) bool, ofst int64, fh uint64) (errc int) {
	// Count the number of slashes in the path
	switch strings.Count(path, "/") {
	case 1:
		if len(path) == 1 { // Root directory
			// Iterate through all hex combinations
			for i := 0; i < len(babel.HEX_SET); i++ {
				for j := 0; j < len(babel.HEX_SET); j++ {
					fill(string(babel.HEX_SET[i])+string(babel.HEX_SET[j]), &fuse.Stat_t{Mode: fuse.S_IFDIR | 0555}, 0)
				}
			}
		} else { // Hex directory
			// Iterate through all wall numbers
			for i := 1; i <= babel.WALLS; i++ {
				fill(strconv.Itoa(i), &fuse.Stat_t{Mode: fuse.S_IFDIR | 0555}, 0)
			}
		}
	case 2:
		// Iterate through all shelf numbers
		for i := 1; i <= babel.SHELVES; i++ {
			fill(strconv.Itoa(i), &fuse.Stat_t{Mode: fuse.S_IFDIR | 0555}, 0)
		}
	case 3:
		// Iterate through all volume numbers
		for i := 1; i <= babel.VOLUMES; i++ {
			fill(strconv.Itoa(i), &fuse.Stat_t{Mode: fuse.S_IFDIR | 0555}, 0)
		}
	case 4:
		// Iterate through all page numbers
		for i := 1; i <= babel.PAGES; i++ {
			fill(strconv.Itoa(i), &fuse.Stat_t{Mode: fuse.S_IFREG | 0444}, 0)
		}
	case 5:
		log.Printf("[</3] Invalid directory: %s\n", path)
		return -fuse.ENOENT
	default:
		log.Printf("[</3] Invalid directory: %s\n", path)
		return -fuse.ENOENT
	}
	return 0
}

func main() {
	log.Println("[*] Starting LOBFS...")
	fs := &BabelFS{}
	host := fuse.NewFileSystemHost(fs)
	host.Mount("", os.Args[1:])
}
