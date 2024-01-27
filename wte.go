package wte

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// createPenaltyMatrix creates the penalty matrix for the Whittaker-Eilers smoother.
// `spacing` is a slice of the same length as `data`, representing the distances between data points.
func createPenaltyMatrix(data, spacing []float64, d int) (*mat.Dense, error) {
	if d < 1 {
		return nil, fmt.Errorf("order must be at least 1")
	}

	n := len(data)

	if len(spacing) != n {
		return nil, fmt.Errorf("spacing slice must be the same length as data slice")
	}

	if n < d {
		return nil, fmt.Errorf("data slice must be at least as long as order")
	}

	D := mat.NewDense(n-d, n, nil)

	for i := 0; i < n-d; i++ {
		for j := 0; j <= d; j++ {
			if i+j < n {
				// Adjusting the penalty based on the spacing
				penalty := binomialCoefficient(d, j) * float64(powInt(-1, j))
				if j > 0 && j < d {
					penalty *= spacing[i+j-1] / spacing[i+j]
				}
				D.Set(i, i+j, penalty)
			}
		}
	}

	return D, nil
}

// binomialCoefficient calculates the binomial coefficient (n choose k).
func binomialCoefficient(n, k int) float64 {
	if k == 0 || k == n {
		return 1
	}

	return binomialCoefficient(n-1, k-1) + binomialCoefficient(n-1, k)
}

// powInt calculates the integer power of a base.
func powInt(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}

	return result
}

// Smooth performs Whittaker-Eilers smoothing on a slice of float64.
// It applies a smoothing operation to the input data using the given lambda (smoothing parameter),
// order of differences, and spacing between data points.
//
// Parameters:
// data    - A slice of float64 representing the data to be smoothed.
// lambda  - A float64 representing the smoothing parameter. Higher values increase smoothness.
// order   - An int specifying the order of differences to be used in smoothing.
// spacing - A slice of float64 representing the spacing between consecutive data points.
//
//	Pass nil if the data is equally spaced. In this case, the function assumes a unit spacing.
//
// Returns:
// A slice of float64 containing the smoothed data. The length of this slice will be the same as the input data.
// An error is returned if the function encounters issues in creating the penalty matrix or solving the system.
//
// Example:
// smoothedData, err := Smooth([]float64{1, 2, 3, 4, 5}, 10.0, 2, nil)
//
//	if err != nil {
//	    // handle error
//	}
//
// fmt.Println(smoothedData) // Output: [Smoothed Data]
func Smooth(data []float64, lambda float64, order int, spacing []float64) ([]float64, error) {
	n := len(data)

	// NOTE(njern): If spacing is nil, we assume that the data is equally spaced.
	if spacing == nil {
		spacing = make([]float64, n)
		for i := range spacing {
			spacing[i] = 1
		}
	}

	// Create the penalty matrix D based on the order
	d, err := createPenaltyMatrix(data, spacing, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create penalty matrix: %w", err)
	}

	// Create the identity matrix
	I := mat.NewDiagDense(n, nil)
	for i := 0; i < n; i++ {
		I.SetDiag(i, 1)
	}

	// Compute D'D
	DtD := mat.NewDense(n, n, nil)
	DtD.Product(d.T(), d)

	// Compute I + lambda * D'D
	systemMatrix := mat.NewDense(n, n, nil)
	systemMatrix.Scale(lambda, DtD)
	systemMatrix.Add(I, systemMatrix)

	// Convert the data slice to a gonum Vector
	y := mat.NewVecDense(n, data)

	// Solve the system (I + lambda * D'D) * x = y
	var x mat.VecDense
	err = x.SolveVec(systemMatrix, y)
	if err != nil {
		return nil, fmt.Errorf("failed to solve the system: %w", err)
	}

	// Extract the smoothed data from the solution vector
	smoothedData := make([]float64, n)
	for i := 0; i < n; i++ {
		smoothedData[i] = x.AtVec(i)
	}

	return smoothedData, nil
}
