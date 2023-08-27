package main

import (
	"fmt"
	"math/big"

	interpolation "github.com/SadPencil/go-lagrange-interpolation"
	"github.com/SadPencil/go-lagrange-interpolation/field"
)

func main() {
	modulus := big.NewInt(11)
	points := []*interpolation.XYPoint{
		{X: &field.Field{Modulus: modulus, Value: big.NewInt(0)}, Y: &field.Field{Modulus: modulus, Value: big.NewInt(1)}},
		{X: &field.Field{Modulus: modulus, Value: big.NewInt(1)}, Y: &field.Field{Modulus: modulus, Value: big.NewInt(5)}},
		{X: &field.Field{Modulus: modulus, Value: big.NewInt(8)}, Y: &field.Field{Modulus: modulus, Value: big.NewInt(9)}},
		{X: &field.Field{Modulus: modulus, Value: big.NewInt(2)}, Y: &field.Field{Modulus: modulus, Value: big.NewInt(5)}},
		{X: &field.Field{Modulus: modulus, Value: big.NewInt(4)}, Y: &field.Field{Modulus: modulus, Value: big.NewInt(0)}},
		{X: &field.Field{Modulus: modulus, Value: big.NewInt(10)}, Y: &field.Field{Modulus: modulus, Value: big.NewInt(7)}},
		{X: &field.Field{Modulus: modulus, Value: big.NewInt(6)}, Y: &field.Field{Modulus: modulus, Value: big.NewInt(4)}},
	}
	poly, err := interpolation.LagrangeInterpolation(points)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(poly.Coefficients); i++ {
		fmt.Printf("f[%v]: %v\n", i, poly.Coefficients[i].String())
	}
	for i := 0; i < 11; i++ {
		x := &field.Field{Modulus: modulus, Value: big.NewInt(int64(i))}
		y := poly.EvaluateAt(x)
		fmt.Printf("(%v, %v)\n", x, y)
	}
}
