package polynomial

import (
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"strings"

	"github.com/SadPencil/go-lagrange-interpolation/field"
)

type Polynomial struct {
	Modulus *big.Int

	Coefficients []*field.Field
}

func NewZero(modulus *big.Int) *Polynomial {
	return &Polynomial{
		Modulus:      modulus,
		Coefficients: make([]*field.Field, 0),
	}
}

func RandomPolynomial(rand *rand.Rand, degree int, modulus *big.Int) *Polynomial {
	coeffs := make([]*field.Field, 0)
	for i := 0; i < degree+1; i++ {
		coeffs = append(coeffs, field.RandomField(rand, modulus))
	}
	for coeffs[degree].IsZero() {
		coeffs[degree] = field.RandomField(rand, modulus)
	}
	return &Polynomial{
		Modulus:      modulus,
		Coefficients: coeffs,
	}
}

func CryptoRandomPolynomial(rand io.Reader, degree int, modulus *big.Int) (*Polynomial, error) {
	coeffs := make([]*field.Field, 0)
	for i := 0; i < degree+1; i++ {
		value, err := field.CryptoRandomField(rand, modulus)
		if err != nil {
			return nil, err
		}
		coeffs = append(coeffs, value)
	}
	for coeffs[degree].IsZero() {
		value, err := field.CryptoRandomField(rand, modulus)
		if err != nil {
			return nil, err
		}
		coeffs[degree] = value
	}
	return &Polynomial{
		Modulus:      modulus,
		Coefficients: coeffs,
	}, nil
}

func (p *Polynomial) Clone() *Polynomial {
	coeffs := make([]*field.Field, 0)
	for _, coeff := range p.Coefficients {
		coeffs = append(coeffs, coeff.Clone())
	}
	return &Polynomial{
		Modulus:      p.Modulus,
		Coefficients: coeffs,
	}
}

func assertSameModulus(firstValue *Polynomial, restValues ...*Polynomial) *big.Int {
	for _, rest := range restValues {
		if firstValue.Modulus.Cmp(rest.Modulus) != 0 {
			panic("modulus mismatch")
		}
	}
	return firstValue.Modulus
}

// Degree returns the degree of the polynomial.
// Note: this method returns `-1` for a zero polynomial, instead of `-∞`.
func (p *Polynomial) Degree() int {
	d := len(p.Coefficients) - 1

	for d >= 0 && p.Coefficients[d].IsZero() {
		d--
	}

	return d
}

func (p *Polynomial) IsZero() bool {
	return p.Degree() == -1
}

func (p *Polynomial) LeadingTermCoefficient() *field.Field {
	degree := p.Degree()
	if degree == -1 {
		return field.NewZero(p.Modulus)
	} else {
		return p.Coefficients[degree]
	}
}

func (p *Polynomial) Shrink() *Polynomial {
	length := p.Degree() + 1
	if length < 1 {
		length = 1
	}
	p.Coefficients = p.Coefficients[:length]
	return p
}

func (p *Polynomial) Negative(poly *Polynomial) *Polynomial {
	modulus := poly.Modulus
	coeffs := make([]*field.Field, 0)

	for _, coeff := range poly.Coefficients {
		coeffs = append(coeffs, new(field.Field).Negative(coeff))
	}

	p.Modulus = modulus
	p.Coefficients = coeffs
	return p
}

func (p *Polynomial) CoefficientAt(index int) *field.Field {
	if index < len(p.Coefficients) {
		return p.Coefficients[index]
	} else {
		return field.NewZero(p.Modulus)
	}
}

func (p *Polynomial) Add(addendA *Polynomial, addendB *Polynomial) *Polynomial {
	modulus := assertSameModulus(addendA, addendB)
	coeffs := make([]*field.Field, 0)
	for i := 0; i < max(len(addendA.Coefficients), len(addendB.Coefficients)); i++ {
		result := new(field.Field).Add(addendA.CoefficientAt(i), addendB.CoefficientAt(i))
		coeffs = append(coeffs, result)
	}

	p.Modulus = modulus
	p.Coefficients = coeffs
	// TODO: p.Shrink() might not be needed here
	p.Shrink()
	return p
}

func (p *Polynomial) Subtract(minuend *Polynomial, subtrahend *Polynomial) *Polynomial {
	modulus := assertSameModulus(minuend, subtrahend)
	coeffs := make([]*field.Field, 0)
	for i := 0; i < max(len(minuend.Coefficients), len(subtrahend.Coefficients)); i++ {
		result := new(field.Field).Subtract(minuend.CoefficientAt(i), subtrahend.CoefficientAt(i))
		coeffs = append(coeffs, result)
	}

	p.Modulus = modulus
	p.Coefficients = coeffs
	p.Shrink()
	return p
}

func (p *Polynomial) Multiply(factorA *Polynomial, factorB *Polynomial) *Polynomial {
	modulus := assertSameModulus(factorA, factorB)

	if factorA.IsZero() || factorB.IsZero() {
		p.Modulus = modulus
		p.Coefficients = []*field.Field{
			field.NewZero(modulus),
		}
		return p
	}

	degreeA := factorA.Degree()
	degreeB := factorB.Degree()

	coeffs := make([]*field.Field, 0)
	for i := 0; i < degreeA+degreeB+1; i++ {
		coeffs = append(coeffs, field.NewZero(modulus))
	}

	for i := 0; i <= degreeA; i++ {
		for j := 0; j <= degreeB; j++ {
			item := new(field.Field).Multiply(factorA.CoefficientAt(i), factorB.CoefficientAt(j))
			coeffs[i+j].Add(coeffs[i+j], item)
		}
	}

	p.Modulus = modulus
	p.Coefficients = coeffs
	p.Shrink()
	return p
}

func DivMod(dividend *Polynomial, divisor *Polynomial) (quotient *Polynomial, remainder *Polynomial) {
	modulus := assertSameModulus(dividend, divisor)

	if divisor.IsZero() {
		panic("polynomial divided by zero")
	}

	n := dividend.Degree()
	m := divisor.Degree()

	q := n - m
	if q <= -1 {
		// quotient: zero, remainder: dividend (cloned)
		quotient = NewZero(modulus)
		remainder = dividend.Clone()

		return quotient, remainder
	}

	iv := new(field.Field).Inverse(divisor.LeadingTermCoefficient())

	quotientCoefficients := make([]*field.Field, 0)
	for i := 0; i < q+1; i++ {
		quotientCoefficients = append(quotientCoefficients, field.NewZero(modulus))
	}
	remainderCoefficients := dividend.Clone().Coefficients

	for i := q; i >= 0; i-- {
		quotientCoefficients[i].Multiply(remainderCoefficients[n], iv)
		n--
		if !quotientCoefficients[i].IsZero() {
			for j := 0; j <= m; j++ {
				item := new(field.Field).Multiply(quotientCoefficients[i], divisor.CoefficientAt(j))
				remainderCoefficients[i+j].Subtract(remainderCoefficients[i+j], item)
			}
		}
	}

	quotient = &Polynomial{
		Modulus:      modulus,
		Coefficients: quotientCoefficients,
	}
	remainder = &Polynomial{
		Modulus:      modulus,
		Coefficients: remainderCoefficients,
	}
	quotient.Shrink()
	remainder.Shrink()

	return quotient, remainder
}

func (p *Polynomial) DivideBy(dividend *Polynomial, divisor *Polynomial) *Polynomial {
	quotient, _ := DivMod(dividend, divisor)
	p.Modulus = quotient.Modulus
	p.Coefficients = quotient.Coefficients
	return p
}

func (p *Polynomial) Modulo(dividend *Polynomial, divisor *Polynomial) *Polynomial {
	_, remainder := DivMod(dividend, divisor)
	p.Modulus = remainder.Modulus
	p.Coefficients = remainder.Coefficients
	return p
}

func (p *Polynomial) EvaluateAt(x *field.Field) *field.Field {
	result := field.NewZero(p.Modulus)
	for i := p.Degree(); i >= 0; i-- {
		item1 := new(field.Field).Multiply(result, x)
		item2 := p.CoefficientAt(i)
		result.Add(item1, item2)
	}
	return result
}

func (p *Polynomial) Equals(other *Polynomial) bool {
	otherClone := other.Clone()
	otherClone.Shrink()
	p.Shrink()

	if len(p.Coefficients) != len(other.Coefficients) {
		return false
	}
	for i := 0; i < len(p.Coefficients); i++ {
		if !p.Coefficients[i].Equals(other.Coefficients[i]) {
			return false
		}
	}

	return true
}

func (p *Polynomial) String() string {
	sb := new(strings.Builder)
	for i := 0; i < len(p.Coefficients); i++ {
		sb.WriteString(fmt.Sprintf("Polynomial[%v]: %v\n", i, p.Coefficients[i]))
	}
	return sb.String()
}
