package polynomial

import (
	"math/big"
	"testing"

	"github.com/SadPencil/go-lagrange-interpolation/field"
)

func TestPolynomial(t *testing.T) {
	modulus := big.NewInt(7)
	poly1 := &Polynomial{
		Modulus: modulus,
		Coefficients: []*field.Field{
			{Modulus: modulus, Value: big.NewInt(1)},
			{Modulus: modulus, Value: big.NewInt(1)},
			{Modulus: modulus, Value: big.NewInt(4)},
			{Modulus: modulus, Value: big.NewInt(5)},
			{Modulus: modulus, Value: big.NewInt(1)},
			{Modulus: modulus, Value: big.NewInt(4)},
		},
	}

	{
		p := poly1.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(1)})
		if p.Value.Cmp(big.NewInt(2)) != 0 {
			t.Fatalf("incorrect computation")
		}
	}
	{

		p := poly1.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(6)})
		if p.Value.Cmp(big.NewInt(3)) != 0 {
			t.Fatalf("incorrect computation")
		}
	}

	poly2 := &Polynomial{
		Modulus: modulus,
		Coefficients: []*field.Field{
			{Modulus: modulus, Value: big.NewInt(1)},
			{Modulus: modulus, Value: big.NewInt(9)},
			{Modulus: modulus, Value: big.NewInt(1)},
			{Modulus: modulus, Value: big.NewInt(9)},
			{Modulus: modulus, Value: big.NewInt(8)},
			{Modulus: modulus, Value: big.NewInt(1)},
			{Modulus: modulus, Value: big.NewInt(0)},
		},
	}

	poly3 := new(Polynomial).Add(poly1, poly2)
	poly3Ans := &Polynomial{
		Modulus: modulus,
		Coefficients: []*field.Field{
			{Modulus: modulus, Value: big.NewInt(2)},
			{Modulus: modulus, Value: big.NewInt(3)},
			{Modulus: modulus, Value: big.NewInt(5)},
			{Modulus: modulus, Value: big.NewInt(0)},
			{Modulus: modulus, Value: big.NewInt(2)},
			{Modulus: modulus, Value: big.NewInt(5)},
		},
	}

	if !poly3.Equals(poly3Ans) {
		t.Fatalf("incorrect computation")
	}

	poly4 := new(Polynomial).Multiply(poly1, poly2)

	for i := int64(0); i < 7; i++ {
		poly4p := poly4.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(i)})
		poly1p := poly1.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(i)})
		poly2p := poly2.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(i)})
		poly4pans := new(field.Field).Multiply(poly1p, poly2p)
		if !poly4pans.Equals(poly4p) {
			t.Fatalf("incorrect computation")
		}
	}

	poly5 := new(Polynomial).DivideBy(poly4, poly1)
	poly6 := new(Polynomial).Modulo(poly4, poly1)

	for i := int64(0); i < 7; i++ {
		poly4p := poly4.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(i)})
		poly1p := poly1.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(i)})
		poly5p := poly5.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(i)})
		poly6p := poly6.EvaluateAt(&field.Field{Modulus: modulus, Value: big.NewInt(i)})

		poly4pans := new(field.Field).Multiply(poly1p, poly5p)
		poly4pans.Add(poly4pans, poly6p)
		if !poly4pans.Equals(poly4p) {
			t.Fatalf("incorrect computation")
		}
	}

}
