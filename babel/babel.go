package babel

import (
	"crypto/sha256"
	crand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strings"
)

const (
	PAGE_LENGTH = 40 * 80
	WALLS       = 4
	SHELVES     = 5
	VOLUMES     = 32
	PAGES       = 410
	MAX_HEX_LEN = 3260
)

var BABEL_SET = []rune{' ', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', ',', '.'}

var HEX_SET = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}


type Address struct {
	Hex    string
	Wall   uint32
	Shelf  uint32
	Volume uint32
	Page   uint32
}

func (a Address) String() string {
	return fmt.Sprintf("%d:%d:%d:%d:%s", a.Wall, a.Shelf, a.Volume, a.Page, a.Hex)
}


func GeneratePage(addr Address) ([]rune, error) {
    if addr.Wall >= WALLS {
        return nil, errors.New("wall number must be between 0 and 3")
    }
    if addr.Shelf >= SHELVES {
        return nil, errors.New("shelf number must be between 0 and 4")
    }
    if addr.Volume >= VOLUMES {
        return nil, errors.New("volume number must be between 0 and 31")
    }
    if addr.Page >= PAGES {
        return nil, errors.New("page number must be between 0 and 409")
    }
    if len(addr.Hex) > MAX_HEX_LEN {
        return nil, errors.New("hex address must be at most 3260 characters")
    }

    // Convert the address to a big integer
    bi := AddressToBigInt(addr)

    // Hash the address
    hash := sha256.Sum256(bi.Bytes())
    bi = new(big.Int).SetBytes(hash[:])

    // Generate the page content
    bitLen := bi.BitLen()
    page := []rune{}
    base := big.NewInt(int64(len(BABEL_SET)))
    for i := 0; i < PAGE_LENGTH; i++ {
        letterIdx := new(big.Int).Mod(bi, base).Int64()
        page = append(page, BABEL_SET[letterIdx])
        bi = RotateBigInt(bi, bitLen)
    }

    // Apply transformation to page
    page = shuffleRunes(page, bi)

    return page, nil
}


func AddressToBigInt(addr Address) *big.Int {
	bi := FromHex([]rune(addr.Hex))
	multiplier := big.NewInt(WALLS * SHELVES * VOLUMES * PAGES)
	bi.Mul(bi, multiplier)
	bi.Add(bi, big.NewInt(int64(addr.Wall*SHELVES*VOLUMES*PAGES+addr.Shelf*VOLUMES*PAGES+addr.Volume*PAGES+addr.Page)))
	return bi
}


func BigIntToAddress(bi *big.Int) Address {
    pageBI := new(big.Int).Set(bi)
    pageBI.Mod(pageBI, big.NewInt(PAGES))
    page := uint32(pageBI.Int64())
    bi.Div(bi, big.NewInt(PAGES))

    volumeBI := new(big.Int).Set(bi)
    volumeBI.Mod(volumeBI, big.NewInt(VOLUMES))
    volume := uint32(volumeBI.Int64())
    bi.Div(bi, big.NewInt(VOLUMES))

    shelfBI := new(big.Int).Set(bi)
    shelfBI.Mod(shelfBI, big.NewInt(SHELVES))
    shelf := uint32(shelfBI.Int64())
    bi.Div(bi, big.NewInt(SHELVES))

    wallBI := new(big.Int).Set(bi)
    wallBI.Mod(wallBI, big.NewInt(WALLS))
    wall := uint32(wallBI.Int64())
    bi.Div(bi, big.NewInt(WALLS))

    // The remaining part of the BigInt is the Hex
    hex := ToHex(bi)

    return Address{Hex: string(hex), Wall: wall, Shelf: shelf, Volume: volume, Page: page}
}

func Search(content string) Address {
	// Generate a random number between 0 and the number of remaining positions after placing the content
	randPos, err := crand.Int(crand.Reader, big.NewInt(int64(PAGE_LENGTH-len(content))))
	if err != nil {
		panic(err)
	}

	// Split the page into two parts at the random position
	padBeforeLen := randPos.Uint64()
	padAfterLen := PAGE_LENGTH - padBeforeLen - uint64(len(content))

	// Pad content to PAGE_LENGTH with spaces at random position
	paddedContent := strings.Repeat(" ", int(padBeforeLen)) + content + strings.Repeat(" ", int(padAfterLen))

	// Convert paddedContent into a big integer in the base of BABEL_SET
	bi := new(big.Int)
	for _, char := range paddedContent {
		bi.Mul(bi, big.NewInt(int64(len(BABEL_SET))))
		bi.Add(bi, big.NewInt(int64(indexOf(BABEL_SET, char))))
	}

	// Convert big integer into an address
	addr := BigIntToAddress(bi)
	return addr
}


func ToHex(bi *big.Int) []rune {
	base := big.NewInt(int64(len(HEX_SET)))
	zero := big.NewInt(0)
	res := []rune{}
	for bi.Cmp(zero) > 0 {
		mod := &big.Int{}
		bi.DivMod(bi, base, mod)
		res = append([]rune{HEX_SET[mod.Int64()]}, res...)
	}
	return res
}

func FromHex(r []rune) *big.Int {
	base := big.NewInt(int64(len(HEX_SET)))
	bi := big.NewInt(0)
	for _, v := range r {
		index := indexOf(HEX_SET, v)
		if index == -1 {
			panic(fmt.Sprintf("invalid character %c in hex string", v))
		}
		bi.Mul(bi, base)
		value := big.NewInt(int64(index))
		bi.Add(bi, value)
	}
	return bi
}


func RotateBigInt(bi *big.Int, bitLen int) *big.Int {
    rotationAmount := int(bi.Mod(bi, big.NewInt(int64(bitLen))).Int64())
    shifted := new(big.Int).Lsh(bi, uint(rotationAmount))
    msb := new(big.Int).Rsh(bi, uint(bitLen-rotationAmount))
    rotated := new(big.Int).Or(shifted, msb)
    return rotated
}

func shuffleRunes(runes []rune, bi *big.Int) []rune {
    r := mrand.New(mrand.NewSource(bi.Int64()))

    // Fisher-Yates shuffle
    for i := len(runes) - 1; i > 0; i-- {
        j := r.Intn(i + 1)
        runes[i], runes[j] = runes[j], runes[i]
    }

    return runes
}

func indexOf(runes []rune, target rune) int {
	for i, v := range runes {
		if v == target {
			return i
		}
	}
	return -1
}
