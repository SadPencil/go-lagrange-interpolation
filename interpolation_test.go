package interpolation

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/SadPencil/go-lagrange-interpolation/field"
	"github.com/SadPencil/go-lagrange-interpolation/polynomial"
)

func TestLagrangeInterpolation(t *testing.T) {
	// As of Go 1.20 there is no reason to call Seed with a random value
	// rand.Seed(time.Now().UnixNano())
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	degree := 100
	modulusHex := "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab"
	modulus := new(big.Int)
	_, success := modulus.SetString(modulusHex, 0)
	if !success {
		t.Fatalf("failed to parse modulus string")
	}

	t.Logf("Generating a polynomial with degree %v, whose coefficients are large field elements...", degree)
	poly := polynomial.RandomPolynomial(randSource, degree, modulus)
	t.Logf("Polynomial generated.")

	t.Logf("Randomly selecting %v points...", degree+1)
	points := make([]*XYPoint, 0)
	for i := 0; i < degree+1; i++ {
		x := field.RandomField(randSource, modulus)
		y := poly.EvaluateAt(x)

		points = append(points, &XYPoint{X: x, Y: y})
	}
	t.Logf("Points selected.")

	t.Logf("Performing lagrange interpolation...")
	resultPoly, err := LagrangeInterpolation(points)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("Interpolation done.")

	t.Logf("Checking equivalent of two polynomials...")
	equal := poly.Equals(resultPoly)
	if equal {
		t.Logf("Check passed.")
	} else {
		t.Fatalf("Polynomials are not equal.")
	}
}
