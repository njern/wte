package wte

import (
	"math"
	"testing"
)

func TestSmooth(t *testing.T) {
	testData := []float64{1.1, 1.9, 3.1, 3.91, 5.0, 6.02, 7.01, 7.7, 9.0, 10.0}

	type args struct {
		data    []float64
		lambda  float64
		order   int
		spacing []float64
	}
	tests := []struct {
		name    string
		args    args
		want    []float64
		wantErr bool
	}{
		{"invalid: order=0", args{data: testData, lambda: 2e1, order: 0, spacing: nil}, nil, true},
		{"invalid: order > len(data)", args{data: testData, lambda: 2e2, order: 13, spacing: nil}, nil, true},
		{"invalid: len(spacing) != len(data)", args{data: testData, lambda: 2e2, order: 13, spacing: []float64{1}}, nil, true},

		{"L=2e1, order=1, no spacing", args{data: testData, lambda: 2e1, order: 1, spacing: nil}, []float64{4.120121432424289, 4.271127504045504, 4.540689950868993, 4.882286895235932, 5.272498184364668, 5.676334382711638, 6.06298730019419, 6.402289582686452, 6.676706344313036, 6.834958423155272}, false},
		{"L=2e1, order=2, no spacing", args{data: testData, lambda: 2e1, order: 2, spacing: nil}, []float64{1.0376249600205127, 2.0182069230969337, 3.001907638172329, 3.9859355110908243, 4.972403565787932, 5.959628050644621, 6.947305035752462, 7.938149188670791, 8.938009925171329, 9.940829201592251}, false},
		{"L=2e2, order=2, no spacing", args{data: testData, lambda: 2e2, order: 2, spacing: nil}, []float64{1.0275647382281692, 2.014761087020907, 3.0023196121225038, 3.9900286844067137, 4.978165076686679, 5.966605418353506, 6.95533551341487, 7.944608138786678, 8.934949393817767, 9.925662337163043}, false},
		{"L=2e2, order=3, no spacing", args{data: testData, lambda: 2e2, order: 3, spacing: nil}, []float64{1.0583913481986855, 2.025523751412608, 2.997874061833645, 3.9752342362027893, 4.957607763500082, 5.945118610751834, 6.938246340018855, 7.9378492642022005, 8.944813195902658, 9.959341427980325}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Smooth(tt.args.data, tt.args.lambda, tt.args.order, tt.args.spacing)
			if (err != nil) != tt.wantErr {
				t.Errorf("Smooth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// NOTE(njern): To account for floating point errors, compare length of response
			// and that values are within 0.01% of each other instead of exact match.
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d items, got %d", len(tt.want), len(got))
			}

			for i := range got {
				if !within(got[i], tt.want[i], 0.01) {
					t.Errorf("got[%d] = %f, want[%d] = %f", i, got[i], i, tt.want[i])
				}
			}
		})
	}
}

// within checks if two float64 numbers are within `perc` of each other.
func within(a, b, perc float64) bool {
	if a == b {
		return true
	}

	average := (a + b) / 2
	if average == 0 {
		return false
	}

	percentDifference := math.Abs(a-b) / average * 100
	return percentDifference <= perc
}
