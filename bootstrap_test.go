package bootstrap

import (
	"math"
	"math/rand"
	"testing"
)

// TestSum tests SumAggregator.
func TestSum(t *testing.T) {
	sum := NewSumAggregator()
	aggregate := sum.Aggregate([]float64{0, 1, 2, 3, 4})
	if aggregate != 10.0 {
		t.Errorf("expected aggregate %f; got %f", 10.0, aggregate)
	}
}

// TestAverage tests AverageAggregator.
func TestAverage(t *testing.T) {
	avg := NewAverageAggregator()
	aggregate := avg.Aggregate([]float64{0, 1, 2, 3, 4})
	if aggregate != 2.0 {
		t.Errorf("expected aggregate %f; got %f", 2.0, aggregate)
	}
	aggregate = avg.Aggregate([]float64{})
	if aggregate != 0.0 {
		t.Errorf("expected aggregate %f; got %f", 0.0, aggregate)
	}
}

// TestQuantile tests QuantileAggregator.
func TestQuantile(t *testing.T) {
	qnt := NewQuantileAggregator(0.5)
	aggregate := qnt.Aggregate([]float64{4, 3, 2, 1, 0})
	if aggregate != 2.0 {
		t.Errorf("expected aggregate %f; got %f", 2.0, aggregate)
	}
	qnt = NewQuantileAggregator(1.0)
	aggregate = qnt.Aggregate([]float64{0, 1, 2, 3, 4})
	if aggregate != 4.0 {
		t.Errorf("expected aggregate %f; got %f", 4.0, aggregate)
	}
	qnt = NewQuantileAggregator(0.33)
	aggregate = qnt.Aggregate([]float64{0, 1, 2, 3, 4})
	if aggregate != 1.0 {
		t.Errorf("expected aggregate %f; got %f", 1.0, aggregate)
	}
	aggregate = qnt.Aggregate([]float64{})
	if !math.IsNaN(aggregate) {
		t.Errorf("expected aggregate NaN; got %f", aggregate)
	}
}

// TestBasicResampler tests the resampler.
func TestBasicResampler(t *testing.T) {
	resampler := NewBasicResampler(NewSumAggregator(), 2000)
	resampler.r.Seed(0)
	resampler.Resample([]float64{0, 1, 2, 3, 4})
	if min := resampler.Quantile(0); min != 0.0 {
		t.Errorf("expected min to be %f; got %f", 0.0, min)
	}
	if median := resampler.Quantile(0.5); median != 10.0 {
		t.Errorf("expected median to be %f; got %f", 10.0, median)
	}
	if max := resampler.Quantile(1); max != 20.0 {
		t.Errorf("expected max to be %f; got %f", 20.0, max)
	}
	resampler.Reset()
	if nan := resampler.Quantile(1); !math.IsNaN(nan) {
		t.Errorf("expected nan to be %f; got %f", math.NaN(), nan)
	}
}

// TestPresampledResampler tests the resampler.
func TestPresampledResampler(t *testing.T) {
	resampler := NewPresampledResampler(NewSumAggregator(), 2000, 5)
	resampler.Resample([]float64{0, 1, 2, 3, 4})
	if min := resampler.Quantile(0); min != 0.0 {
		t.Errorf("expected min to be %f; got %f", 0.0, min)
	}
	if median := resampler.Quantile(0.5); median != 10.0 {
		t.Errorf("expected median to be %f; got %f", 10.0, median)
	}
	if max := resampler.Quantile(1); max != 20.0 {
		t.Errorf("expected max to be %f; got %f", 20.0, max)
	}
	resampler.Reset()
	if nan := resampler.Quantile(1); !math.IsNaN(nan) {
		t.Errorf("expected nan to be %f; got %f", math.NaN(), nan)
	}
}

func BenchmarkResampler(b *testing.B) {
	resampler := NewBasicResampler(NewSumAggregator(), b.N)
	resampler.r.Seed(0)
	sampleData := make([]float64, 1000)
	for i := range sampleData {
		sampleData[i] = rand.Float64()
	}
	b.ResetTimer()
	resampler.Resample(sampleData)
}

func BenchmarkPresampledResampler(b *testing.B) {
	resampler := NewPresampledResampler(NewSumAggregator(), b.N, 1000)
	sampleData := make([]float64, 1000)
	for i := range sampleData {
		sampleData[i] = rand.Float64()
	}
	b.ResetTimer()
	resampler.Resample(sampleData)
}
