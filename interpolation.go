package interpolation

import (
	"fmt"
	"math/big"

	"github.com/SadPencil/go-lagrange-interpolation/field"
	"github.com/SadPencil/go-lagrange-interpolation/polynomial"
)

type XYPoint struct {
	X *field.Field
	Y *field.Field
}

func LagrangeInterpolation(points []*XYPoint) (*polynomial.Polynomial, error) {
	n := len(points)
	if n == 0 {
		return nil, fmt.Errorf("at least 1 point is expected to interpolate")
	}

	modulus := points[0].X.Modulus
	for _, point := range points {
		if point.X.Modulus.Cmp(modulus) != 0 || point.Y.Modulus.Cmp(modulus) != 0 {
			return nil, fmt.Errorf("modulus mismatch")
		}
	}

	oneField := &field.Field{Modulus: modulus, Value: big.NewInt(1)}
	M := &polynomial.Polynomial{Modulus: modulus, Coefficients: []*field.Field{oneField}}
	f := polynomial.NewZero(modulus)

	for i := 0; i < n; i++ {
		// Don't know how to name this polynomial.
		temp := &polynomial.Polynomial{
			Modulus: modulus,
			Coefficients: []*field.Field{
				new(field.Field).Negative(points[i].X),
				oneField,
			},
		}
		M.Multiply(M, temp)
	}

	for i := 0; i < n; i++ {
		// `temp` is the same with above. Don't know how to name this polynomial. As a result I can't name the function as `GetTemp()`.
		temp := &polynomial.Polynomial{
			Modulus: modulus,
			Coefficients: []*field.Field{
				new(field.Field).Negative(points[i].X),
				oneField,
			},
		}

		m := new(polynomial.Polynomial).DivideBy(M, temp)
		mx := m.EvaluateAt(points[i].X)

		factor := &polynomial.Polynomial{
			Modulus: modulus,
			Coefficients: []*field.Field{
				new(field.Field).DivideBy(points[i].Y, mx),
			},
		}

		item := new(polynomial.Polynomial).Multiply(factor, m)
		f.Add(f, item)
	}

	return f, nil
}
