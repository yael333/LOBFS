package babel

import (
	"errors"
	"fmt"
	"math/rand"
)

const (
	VolumeSize  = 410
	PageSize    = 3200
	ShelfSize   = 32
	WallSize    = 5
	HexagonSize = 4
	Chars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789,.-;:_ "
	CharCount   = len(Chars)
)

type Location struct {
	Hexagon int
	Wall    int
	Shelf   int
	Volume  int
	Page    int
}

type Volume struct {
	Name  string
	Pages [VolumeSize]string
}

type Shelf [ShelfSize]*Volume
type Wall [WallSize]*Shelf
type Hexagon [HexagonSize]*Wall
type Library map[int]*Hexagon // int can be replaced with other identifier type

// NewLocation constructs a new Location instance based on the given values.
func NewLocation(h, w, s, v, p int) Location {
	return Location{h, w, s, v, p}
}

func (l Location) Seed() int64 {
	return int64(l.Hexagon<<24 | l.Wall<<16 | l.Shelf<<8 | l.Volume | l.Page)
}

func (l Location) String() string {
	return fmt.Sprintf("%d-w%d-s%d-v%d-p%d", l.Hexagon, l.Wall, l.Shelf, l.Volume, l.Page)
}

// NewVolume returns a new volume with pages generated based on the location.
func NewVolume(l Location) *Volume {
	v := &Volume{Name: l.String()}
	for i := 0; i < VolumeSize; i++ {
		pageLoc := Location{l.Hexagon, l.Wall, l.Shelf, l.Volume, i}
		v.Pages[i] = GeneratePage(pageLoc.Seed())
	}
	return v
}

func GeneratePage(seed int64) string {
	rng := rand.New(rand.NewSource(seed))
	var page string
	for i := 0; i < PageSize; i++ {
		page += string(Chars[rng.Intn(CharCount)])
	}
	return page
}

// GetHexagon retrieves the hexagon with the given ID, creating it if it doesn't exist.
func (lib *Library) GetHexagon(id int) *Hexagon {
	h, ok := (*lib)[id]
	if !ok {
		h = &Hexagon{}
		(*lib)[id] = h
	}
	return h
}

// GetWall retrieves the wall with the given ID from the hexagon, creating it if it doesn't exist.
func (h *Hexagon) GetWall(id int) *Wall {
	w := (*h)[id]
	if w == nil {
		w = &Wall{}
		(*h)[id] = w
	}
	return w
}

// GetShelf retrieves the shelf with the given ID from the wall, creating it if it doesn't exist.
func (w *Wall) GetShelf(id int) *Shelf {
	s := (*w)[id]
	if s == nil {
		s = &Shelf{}
		(*w)[id] = s
	}
	return s
}

// GetVolume retrieves the volume with the given ID from the shelf, creating it if it doesn't exist.
func (s *Shelf) GetVolume(l Location) *Volume {
	v := (*s)[l.Volume]
	if v == nil {
		v = NewVolume(l)
		(*s)[l.Volume] = v
	}
	return v
}

// GetPage retrieves the page with the given ID from the volume.
func (v *Volume) GetPage(id int) (string, error) {
	if id < 0 || id >= VolumeSize {
		return "", errors.New("invalid page number")
	}
	return v.Pages[id], nil
}

// GetPageAtLocation retrieves the page at the given location in the library.
func (lib *Library) GetPageAtLocation(l Location) (string, error) {
	h := lib.GetHexagon(l.Hexagon)
	w := h.GetWall(l.Wall)
	s := w.GetShelf(l.Shelf)
	v := s.GetVolume(l)
	return v.GetPage(l.Page)
}
