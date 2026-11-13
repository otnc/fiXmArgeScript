package interp

import (
	"math"
	"testing"
	"time"
)

func TestWaitDuration(t *testing.T) {
	cases := []struct {
		ms   float64
		want time.Duration
	}{
		{1.5, 1500 * time.Microsecond},
		{-5, 0},
		{math.NaN(), 0},
		{1e18, maxWait},
		{math.Inf(1), maxWait},
	}
	for _, c := range cases {
		if got := waitDuration(c.ms); got != c.want {
			t.Errorf("waitDuration(%v) = %v, want %v", c.ms, got, c.want)
		}
	}
}
