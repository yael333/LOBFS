package babel

import (
	crand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strings"
	"strconv"
)

const (
	PAGE_LENGTH = 40 * 80
	TITLE_LENGTH = 20
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

// String returns the string representation of an Address.
func (a Address) String() string {
	return fmt.Sprintf("%d:%d:%d:%d:%s", a.Wall, a.Shelf, a.Volume, a.Page, a.Hex)
}

// ValidateWall checks if the wall number is valid.
func ValidateWall(wall uint32) error {
	if wall >= WALLS {
		return errors.New("invalid wall number: must be between 0 and 3")
	}
	return nil
}

// ValidateShelf checks if the shelf number is valid.
func ValidateShelf(shelf uint32) error {
	if shelf >= SHELVES {
		return errors.New("invalid shelf number: must be between 0 and 4")
	}
	return nil
}

// ValidateVolume checks if the volume number is valid.
func ValidateVolume(volume uint32) error {
	if volume >= VOLUMES {
		return errors.New("invalid volume number: must be between 0 and 31")
	}
	return nil
}

// ValidatePage checks if the page number is valid.
func ValidatePage(page uint32) error {
	if page >= PAGES {
		return errors.New("invalid page number: must be between 0 and 409")
	}
	return nil
}

// ValidateHex checks if the hex string is valid.
func ValidateHex(hex string) error {
	if len(hex) > MAX_HEX_LEN {
		return errors.New("invalid hex: length must be at most 3260 characters")
	}
	return nil
}

// GeneratePage generates a page content for the given address.
func GeneratePage(addr Address) ([]rune, error) {
    if err := ValidateWall(addr.Wall); err != nil {
        return nil, err
    }
    if err := ValidateShelf(addr.Shelf); err != nil {
        return nil, err
    }
    if err := ValidateVolume(addr.Volume); err != nil {
        return nil, err
    }
    if err := ValidatePage(addr.Page); err != nil {
        return nil, err
    }

    // Construct unique identifier for this volume
	uniqueId := addr.Wall * SHELVES * VOLUMES * PAGES + addr.Shelf * VOLUMES * PAGES + addr.Volume * PAGES + addr.Page

    // Use the unique identifier as the seed for a new PRNG
    rng := mrand.New(mrand.NewSource(int64(uniqueId)))

    // Generate the page content
    page := []rune{}
    for i := 0; i < PAGE_LENGTH; i++ {
        letterIdx := rng.Intn(len(BABEL_SET))
        page = append(page, BABEL_SET[letterIdx])
    }

    return page, nil
}

// GenerateTitle generates a title for the given address.
func GenerateTitle(addr Address) ([]rune, error) {
    if err := ValidateWall(addr.Wall); err != nil {
        return nil, err
    }
    if err := ValidateShelf(addr.Shelf); err != nil {
        return nil, err
    }
    if err := ValidateVolume(addr.Volume); err != nil {
        return nil, err
    }
    if err := ValidatePage(addr.Page); err != nil {
        return nil, err
    }

    // Construct unique identifier for this volume
    uniqueId := addr.Wall * SHELVES * VOLUMES + addr.Shelf * VOLUMES + addr.Volume

    // Use the unique identifier as the seed for a new PRNG
    rng := mrand.New(mrand.NewSource(int64(uniqueId)))

    // Generate the page content
    page := []rune{}
    for i := 0; i < TITLE_LENGTH; i++ {
        letterIdx := rng.Intn(len(BABEL_SET))
        page = append(page, BABEL_SET[letterIdx])
    }

    return page, nil
}

// AddressToBigInt converts an Address to a big integer.
func AddressToBigInt(addr Address) *big.Int {
	bi := FromHex([]rune(addr.Hex))
	multiplier := big.NewInt(WALLS * SHELVES * VOLUMES * PAGES)
	bi.Mul(bi, multiplier)
	bi.Add(bi, big.NewInt(int64(addr.Wall * SHELVES * VOLUMES * PAGES + addr.Shelf * VOLUMES * PAGES + addr.Volume * PAGES + addr.Page)))
	return bi
}

// BigIntToAddress converts a big integer to an Address.
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

// Search finds an address for the given content.
func Search(content string) Address {
	// Generate a random number between 0 and the number of remaining positions after placing the content
	randPos, err := crand.Int(crand.Reader, big.NewInt(int64(PAGE_LENGTH - len(content))))
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

// ToHex converts a big integer to a hexadecimal representation.
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

// FromHex converts a hexadecimal representation to a big integer.
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

// ParseAddress parses an address string into an Address struct.
func ParseAddress(addrStr string) (Address, error) {
	// Split the string into parts
	parts := strings.Split(addrStr, "/")

	// Create an address with all fields set to zero
	addr := Address{}

	// Set the hex field, if it exists
	if len(parts) > 0 && len(parts[0]) > 0 {
		addr.Hex = parts[0]
	}

	// Set the wall field, if it exists
	if len(parts) > 1 && len(parts[1]) > 0 {
		wall, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			return Address{}, err
		}
		addr.Wall = uint32(wall)
	}

	// Set the shelf field, if it exists
	if len(parts) > 2 && len(parts[2]) > 0 {
		shelf, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			return Address{}, err
		}
		addr.Shelf = uint32(shelf)
	}

	// Set the volume field, if it exists
	if len(parts) > 3 && len(parts[3]) > 0 {
		volume, err := strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			return Address{}, err
		}
		addr.Volume = uint32(volume)
	}

	// Set the page field, if it exists
	if len(parts) > 4 && len(parts[4]) > 0 {
		page, err := strconv.ParseUint(parts[4], 10, 32)
		if err != nil {
			return Address{}, err
		}
		addr.Page = uint32(page)
	}

	return addr, nil
}

// RotateBigInt rotates a big integer by a certain number of bits.
func RotateBigInt(bi *big.Int, bitLen int, pageNum uint32) *big.Int {
    rotationAmount := int(new(big.Int).Add(bi, big.NewInt(int64(pageNum))).Mod(bi, big.NewInt(int64(bitLen))).Int64())
    shifted := new(big.Int).Lsh(bi, uint(rotationAmount))
    msb := new(big.Int).Rsh(bi, uint(bitLen-rotationAmount))
    rotated := new(big.Int).Or(shifted, msb)
    return rotated
}

// shuffleRunes shuffles the order of runes using a random number generator.
func shuffleRunes(runes []rune, bi *big.Int) []rune {
    r := mrand.New(mrand.NewSource(bi.Int64()))

    // Fisher-Yates shuffle
    for i := len(runes) - 1; i > 0; i-- {
        j := r.Intn(i + 1)
        runes[i], runes[j] = runes[j], runes[i]
    }

    return runes
}

// indexOf returns the index of a rune in a rune slice.
func indexOf(runes []rune, target rune) int {
	for i, v := range runes {
		if v == target {
			return i
		}
	}
	return -1
}
