package utils_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Taiterbase/vtrips/apps/backend/pkg/utils"
)

func BenchmarkSlicesTest(t *testing.B) {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	t.Run("intersection_10000000", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			a := make([]int, rnd.Intn(10000000))
			b := make([]int, rnd.Intn(10000000))
			// create a new slice for each iteration
			for i := 0; i < len(a); i++ {
				a[i] = rnd.Intn(10000000)
			}
			for i := 0; i < len(b); i++ {
				b[i] = rnd.Intn(10000000)
			}
			t.ResetTimer()
			utils.Intersection(a, b)
		}
	})
	t.Run("intersection_1000000", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			a := make([]int, 1000000)
			b := make([]int, 1000000)
			// create a new slice for each iteration
			for i := 0; i < len(a); i++ {
				a[i] = rnd.Intn(10000000)
			}
			for i := 0; i < len(b); i++ {
				b[i] = rnd.Intn(10000000)
			}
			t.ResetTimer()
			utils.Intersection(a, b)
		}
	})
	t.Run("intersection_100000", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			a := make([]int, 100000)
			b := make([]int, 100000)
			// create a new slice for each iteration
			for i := 0; i < len(a); i++ {
				a[i] = rnd.Intn(10000000)
			}
			for i := 0; i < len(b); i++ {
				b[i] = rnd.Intn(10000000)
			}
			t.ResetTimer()
			utils.Intersection(a, b)
		}
	})
}
