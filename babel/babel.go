package babel

import (
	"errors"
	"math/big"
	"math/rand"
	"strings"
)

const (
	Chars             = "abcdefghijklmnopqrstuvwxyz, ."
	HexBase           = 36
	LocationMultiplier = 1000
	PageMultiplier    = 1 << 62
	MaxLocationValue  = 1000
	PaddingLength     = 10
)

var (
	charToNum map[rune]*big.Int
	numToChar map[int32]rune
)

type Hex string

type Location struct {
	Wall   int
	Shelf  int
	Volume int
	Page   int
}

type Page string

type Library struct {
	rand *rand.Rand
}

var ErrInvalidChar = errors.New("invalid character")

func init() {
	charToNum = make(map[rune]*big.Int, len(Chars))
	numToChar = make(map[int32]rune, len(Chars))
	for i, c := range Chars {
		charToNum[c] = big.NewInt(int64(i))
		numToChar[int32(i)] = c
	}
}

func (h Hex) ToBigInt() (*big.Int, error) {
	n := new(big.Int)
	_, ok := n.SetString(string(h), HexBase)
	if !ok {
		return nil, ErrInvalidChar
	}
	return n, nil
}

func (l Location) ToBigInt() *big.Int {
	n := big.NewInt(int64(l.Wall))
	n.Mul(n, big.NewInt(LocationMultiplier))
	n.Add(n, big.NewInt(int64(l.Shelf)))
	n.Mul(n, big.NewInt(LocationMultiplier))
	n.Add(n, big.NewInt(int64(l.Volume)))
	n.Mul(n, big.NewInt(LocationMultiplier))
	n.Add(n, big.NewInt(int64(l.Page)))
	return n
}

func (p Page) ToBigInt() (*big.Int, error) {
	n := new(big.Int)
	for _, c := range p {
		v, ok := charToNum[c]
		if !ok {
			return nil, ErrInvalidChar
		}
		n.Mul(n, big.NewInt(int64(len(Chars))))
		n.Add(n, v)
	}
	return n, nil
}

func (p *Page) FromBigInt(n *big.Int) {
	var s strings.Builder
	mod := new(big.Int)
	for n.BitLen() > 0 {
		n, mod = n.DivMod(n, big.NewInt(int64(len(Chars))), mod)
		s.WriteRune(numToChar[int32(mod.Int64())])
	}
	*p = Page(s.String())
}

func NewLibrary(seed int64) *Library {
	return &Library{rand: rand.New(rand.NewSource(seed))}
}

func (lib *Library) GeneratePage(h Hex, l Location) (Page, error) {
	hexNum, err := h.ToBigInt()
	if err != nil {
		return "", err
	}
	locNum := l.ToBigInt()

	pageNum := new(big.Int).Mul(locNum, big.NewInt(PageMultiplier))
	pageNum.Add(pageNum, hexNum)

	p := new(Page)
	p.FromBigInt(pageNum)
	return *p, nil
}

func (lib *Library) SearchPage(text string) (h Hex, l Location, p Page, err error) {
	text = lib.padText(text)

	p = Page(text)
	pageNum, err := p.ToBigInt()
	if err != nil {
		return
	}

	l = Location{
		Wall:   lib.rand.Intn(MaxLocationValue),
		Shelf:  lib.rand.Intn(MaxLocationValue),
		Volume: lib.rand.Intn(MaxLocationValue),
		Page:   lib.rand.Intn(MaxLocationValue),
	}
	locNum := l.ToBigInt()

	hexNum := new(big.Int).Add(pageNum, new(big.Int).Mul(locNum, big.NewInt(PageMultiplier)))
	h = Hex(hexNum.Text(HexBase))

	return
}



func (lib *Library) randString(n int) string {
	var s strings.Builder
	for i := 0; i < n; i++ {
		s.WriteRune(rune(Chars[lib.rand.Intn(len(Chars))]))
	}
	return s.String()
}

func (lib *Library) padText(text string) string {
	prefix := lib.randString(PaddingLength)
	suffix := lib.randString(PaddingLength)
	return prefix + text + suffix
}

