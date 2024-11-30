package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// merge sort sequencial
func mergeSortSeq(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := mergeSortSeq(arr[:mid])
	right := mergeSortSeq(arr[mid:])
	return merge(left, right)
}

// merge sort paralelo
func mergeSortParallel(arr []int, maxDepth int, currentDepth int) []int {
	if len(arr) <= 1 {
		return arr
	}
	if currentDepth >= maxDepth {
		// se atingir a profundidade máxima, executa o merge sort sequencial
		return mergeSortSeq(arr)
	}

	mid := len(arr) / 2
	var leftResult, rightResult []int

	var wg sync.WaitGroup // waitGroup para sincronizar as goroutines
	wg.Add(2)             // adiciona 2 goroutines ao WaitGroup

	// executa merge sort paralelo na metade esquerda
	go func() {
		defer wg.Done() // decrementa o contador do WaitGroup
		leftResult = mergeSortParallel(arr[:mid], maxDepth, currentDepth+1)
	}()

	// executa merge sort paralelo na metade direita
	go func() {
		defer wg.Done()
		rightResult = mergeSortParallel(arr[mid:], maxDepth, currentDepth+1)
	}()

	wg.Wait() // espera as duas goroutines terminarem

	return merge(leftResult, rightResult)
}

// calculo de speedup
// retorna o tempo sequencial, o tempo paralelo e o speedup representados em float64
func benchmark(vectorSize int, maxDepth int) (float64, float64, float64) {
	arr := make([]int, vectorSize)
	for i := range arr {
		arr[i] = rand.Intn(1_000_000)
	}

	// execução sequencial
	startSeq := time.Now()
	mergeSortSeq(arr)
	timeSeq := time.Since(startSeq).Seconds()

	// execução paralela
	startPar := time.Now()
	mergeSortParallel(arr, maxDepth, 0)
	timePar := time.Since(startPar).Seconds()

	// calculo de speedup
	speedup := timeSeq / timePar
	return timeSeq, timePar, speedup
}

// avaliação de desempenho com diferentes tamanhos de vetor e número de processadores
func evaluate() {
	vectorSizes := []int{1_000, 10_000, 100_000, 1_000_000, 10_000_000} // granularidade
	numProcessors := []int{2, 4, 6, 8, 10, 12}                   // num de processadores

	for _, size := range vectorSizes {
		fmt.Printf("Vector Size: %d\n", size)
		for _, procs := range numProcessors {
			runtime.GOMAXPROCS(procs)
			fmt.Printf("  Processors: %d\n", procs)

			// profundidade máxima para o merge sort paralelo
			// isso é feito para evitar muitas goroutines e o paralelo não ser eficiente
			maxDepth := procs * 2

			timeSeq, timePar, speedup := benchmark(size, maxDepth)
			fmt.Printf("    Sequential Time: %.2fs, Parallel Time: %.2fs, Speedup: %.2f\n",
				timeSeq, timePar, speedup)
		}
	}
}

func main() {
	evaluate()
}
