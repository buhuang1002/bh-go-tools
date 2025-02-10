package bigint

import (
	"fmt"
	"math/big"
)

// BigIntMaxSerializedLen is the max length of a byte slice representing a CBOR serialized big.
const BigIntMaxSerializedLen = 128

func NewInt(i int64) *big.Int {
	return big.NewInt(0).SetInt64(i)
}

func NewIntUnsigned(i uint64) *big.Int {
	return big.NewInt(0).SetUint64(i)
}

func NewFromGo(i *big.Int) *big.Int {
	return big.NewInt(0).Set(i)
}

func Zero() *big.Int {
	return NewInt(0)
}

// PositiveFromUnsignedBytes interprets b as the bytes of a big-endian unsigned
// integer and returns a positive Int with this absolute value.
func PositiveFromUnsignedBytes(b []byte) *big.Int {
	i := big.NewInt(0).SetBytes(b)
	return i
}

// MustFromString convers dec string into big integer and panics if conversion
// is not sucessful.
func MustFromString(s string) *big.Int {
	v, err := FromString(s)
	if err != nil {
		panic(err)
	}
	return v
}

func FromString(s string) (*big.Int, error) {
	v, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse string as a big int")
	}

	return v, nil
}

func Product(ints ...*big.Int) *big.Int {
	p := NewInt(1)
	for _, i := range ints {
		p = Mul(p, i)
	}
	return p
}

func Mul(a, b *big.Int) *big.Int {
	return big.NewInt(0).Mul(a, b)
}

func Div(a, b *big.Int) *big.Int {
	return big.NewInt(0).Div(a, b)
}

func Mod(a, b *big.Int) *big.Int {
	return big.NewInt(0).Mod(a, b)
}

func Add(a, b *big.Int) *big.Int {
	return big.NewInt(0).Add(a, b)
}

func Sum(ints ...*big.Int) *big.Int {
	sum := Zero()
	for _, i := range ints {
		sum = Add(sum, i)
	}
	return sum
}

func Subtract(num1 *big.Int, ints ...*big.Int) *big.Int {
	sub := num1
	for _, i := range ints {
		sub = Sub(sub, i)
	}
	return sub
}

func Sub(a, b *big.Int) *big.Int {
	return big.NewInt(0).Sub(a, b)
}

func Neg(a *big.Int) *big.Int {
	return big.NewInt(0).Neg(a)
}

// Returns a**e unless e <= 0 (in which case returns 1).
func Exp(a *big.Int, e *big.Int) *big.Int {
	return big.NewInt(0).Exp(a, e, nil)
}

// Returns x << n
func Lsh(a *big.Int, n uint) *big.Int {
	return big.NewInt(0).Lsh(a, n)
}

// Returns x >> n
func Rsh(a *big.Int, n uint) *big.Int {
	return big.NewInt(0).Rsh(a, n)
}

func BitLen(a *big.Int) uint {
	return uint(a.BitLen())
}

func Max(x, y *big.Int) *big.Int {
	if Eq(x, Zero()) && Eq(x, y) {
		if x.Sign() != 0 {
			return y
		}
		return x
	}

	if Gt(x, y) {
		return x
	}
	return y
}

func Min(x, y *big.Int) *big.Int {
	if Eq(x, Zero()) && Eq(x, y) {
		if x.Sign() != 0 {
			return x
		}
		return y
	}

	if Lt(x, y) {
		return x
	}
	return y
}

func Cmp(a, b *big.Int) int {
	return a.Cmp(b)
}

func Eq(a, b *big.Int) bool {
	return Cmp(a, b) == 0
}

func Ne(a, b *big.Int) bool {
	return Cmp(a, b) != 0
}

func Gt(a, b *big.Int) bool {
	return Cmp(a, b) > 0
}

func Lt(a, b *big.Int) bool {
	return Cmp(a, b) < 0
}

func Ge(a, b *big.Int) bool {
	return Cmp(a, b) >= 0
}

func Le(a, b *big.Int) bool {
	return Cmp(a, b) <= 0
}

func FromBytes(buf []byte) (*big.Int, error) {
	if len(buf) == 0 {
		return NewInt(0), nil
	}

	var negative bool
	switch buf[0] {
	case 0:
		negative = false
	case 1:
		negative = true
	default:
		return Zero(), fmt.Errorf("big int prefix should be either 0 or 1, got %d", buf[0])
	}

	i := big.NewInt(0).SetBytes(buf[1:])
	if negative {
		i.Neg(i)
	}

	return i, nil
}
