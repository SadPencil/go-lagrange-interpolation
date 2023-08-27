# go-lagrange-interpolation

A Lagrange interpolation implementation in Go. 

Given a list of points in $\mathbb{F}_p$, outputs all coefficients of interpolating polynomial $P(x)$.

I am shocked by the fact that the majority of Lagrange interpolation implementations only calculate $P(0)$ and not all the coefficients. However, I have recently discovered a well-written implementation in C++ that computes all coefficients of $P(x)$, available on [OI Wiki](https://oi-wiki.org/math/numerical/lagrange/). Therefore, I have personally created a reimplementation of this in Go, which is available in this repository.

## Example

The following [example codes](example/example.go) domonstrate Lagrange interpolation, resulting in a polynomial with coefficients `1, 1, 4, 5, 1, 4`.

```go
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
```
