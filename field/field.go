package field

import (
	cryptoRand "crypto/rand"
	"fmt"
	"io"
	"math/big"
	"math/rand"
)

type Field struct {
	Value   *big.Int
	Modulus *big.Int
}

func NewZero(modulus *big.Int) *Field {
	return &Field{Modulus: modulus, Value: big.NewInt(0)}
}

func RandomField(rand *rand.Rand, modulus *big.Int) *Field {
	value := new(big.Int).Rand(rand, modulus)
	return &Field{Modulus: modulus, Value: value}
}

func CryptoRandomField(rand io.Reader, modulus *big.Int) (*Field, error) {
	value, err := cryptoRand.Int(rand, modulus)
	if err != nil {
		return nil, err
	}
	return &Field{Modulus: modulus, Value: value}, nil
}

func (fp *Field) String() string {
	return fmt.Sprintf("%v (mod %v)", fp.Value.String(), fp.Modulus.String())
}

func (fp *Field) Clone() *Field {
	return &Field{
		Value:   fp.Value,
		Modulus: fp.Modulus,
	}
}

func assertSameModulus(firstValue *Field, restValues ...*Field) *big.Int {
	for _, rest := range restValues {
		if firstValue.Modulus.Cmp(rest.Modulus) != 0 {
			panic("modulus mismatch")
		}
	}
	return firstValue.Modulus
}

// Add sets fp as "addendA + addendB".
func (fp *Field) Add(addendA *Field, addendB *Field) *Field {
	modulus := assertSameModulus(addendA, addendB)

	result := new(big.Int).Add(addendA.Value, addendB.Value)
	result.Mod(result, addendA.Modulus)

	fp.Modulus = modulus
	fp.Value = result
	return fp
}

// Subtract sets fp as "minuend - subtrahend".
func (fp *Field) Subtract(minuend *Field, subtrahend *Field) *Field {
	modulus := assertSameModulus(minuend, subtrahend)

	result := new(big.Int).Add(minuend.Modulus, minuend.Value)
	result.Sub(result, subtrahend.Value)
	result.Mod(result, minuend.Modulus)

	fp.Modulus = modulus
	fp.Value = result
	return fp
}

// Multiply sets fp as "factorA * factorB".
func (fp *Field) Multiply(factorA *Field, factorB *Field) *Field {
	modulus := assertSameModulus(factorA, factorB)

	result := new(big.Int).Mul(factorA.Value, factorB.Value)
	result.Mod(result, factorA.Modulus)

	fp.Modulus = modulus
	fp.Value = result
	return fp
}

// DivideBy sets fp as "dividend / divisor".
func (fp *Field) DivideBy(dividend *Field, divisor *Field) *Field {
	modulus := assertSameModulus(dividend, divisor)

	inv := new(Field).Inverse(divisor)
	result := new(big.Int).Mul(dividend.Value, inv.Value)
	result.Mod(result, dividend.Modulus)

	fp.Modulus = modulus
	fp.Value = result
	return fp
}

// Inverse sets fp the multiplicative inverse of field element fp, i.e., "fp^{-1}".
func (fp *Field) Inverse(element *Field) *Field {
	result := new(big.Int).ModInverse(element.Value, element.Modulus)
	fp.Modulus = element.Modulus
	fp.Value = result
	return fp
}

// Negative sets fp the additive inverse of a field element fp, i.e., "-fp".
func (fp *Field) Negative(element *Field) *Field {
	result := new(big.Int).Sub(element.Modulus, element.Value)
	result.Mod(result, element.Modulus)
	fp.Modulus = element.Modulus
	fp.Value = result
	return fp
}

func (fp *Field) IsZero() bool {
	return fp.Value.BitLen() == 0
}

func (fp *Field) Equals(other *Field) bool {
	return fp.Modulus == other.Modulus && fp.Value.Cmp(other.Value) == 0
}
