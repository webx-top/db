package factory

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestSelecterWeightedRoundRobin(t *testing.T) {
	length := 10
	note := make([]int, length)
	percent := make([]float32, length)
	wrr := NewWeightedRoundRobin(length) //[1]
	var total int
	for i := range note {
		note[i] = rand.Intn(length) + 1
		wrr.SetWeight(i, int32(note[i])) //[2]
		total += note[i]
	}
	for i, weight := range note {
		percent[i] = float32(weight) / float32(total)
	}
	times := make([]int, length)
	testTimes := 100000
	for i := 0; i < testTimes; i++ {
		s := wrr.Select(length) //[3]
		times[s]++
	}
	fmt.Println("--------------------------------------------------")

	for k, v := range note {
		fmt.Println(k, v, percent[k], `->`, times[k], float32(times[k])/float32(testTimes))
	}
}

func TestSelecterMathNGCD(t *testing.T) {
	for _, a := range [][]int32{
		[]int32{
			2, 4, 6, 8,
		}, []int32{
			4, 12, 16, 32,
		}, []int32{
			8, 6, 4, 2, 0,
		}, []int32{
			300, 150, 100, 0,
		},
	} {
		fmt.Println(MathNGCD(a, len(a)))
	}
}
