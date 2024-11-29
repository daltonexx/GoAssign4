package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// quicksort sequencial
func quicksortSeq(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	pivot := arr[len(arr)/2]
	var left, middle, right []int
	for _, v := range arr {
		if v < pivot {
			left = append(left, v)
		} else if v == pivot {
			middle = append(middle, v)
		} else {
			right = append(right, v)
		}
	}
	left = quicksortSeq(left)
	right = quicksortSeq(right)
	return append(append(left, middle...), right...)
}

// quicksort paralelo
func quicksortParallel(arr []int, maxDepth int, currentDepth int) []int {
	if len(arr) <= 1 {
		return arr
	}
	if currentDepth >= maxDepth {
		return quicksortSeq(arr)
	}

	pivot := arr[len(arr)/2]
	var left, middle, right []int
	for _, v := range arr {
		if v < pivot {
			left = append(left, v)
		} else if v == pivot {
			middle = append(middle, v)
		} else {
			right = append(right, v)
		}
	}

	var wg sync.WaitGroup // waitGroup para sincronizar as goroutines
	wg.Add(2)             // adiciona 2 goroutines ao WaitGroup

	var leftResult, rightResult []int

	// executa quicksort paralelo nas duas metades
	go func() {
		defer wg.Done() // decrementa o contador do WaitGroup
		leftResult = quicksortParallel(left, maxDepth, currentDepth+1)
	}()

	go func() {
		defer wg.Done()
		rightResult = quicksortParallel(right, maxDepth, currentDepth+1)
	}()

	wg.Wait()

	// junta os resultados
	return append(append(leftResult, middle...), rightResult...)
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
	quicksortSeq(arr)
	timeSeq := time.Since(startSeq).Seconds()

	// execução paralela
	startPar := time.Now()
	quicksortParallel(arr, maxDepth, 0)
	timePar := time.Since(startPar).Seconds()

	// calcula o speedup
	speedup := timeSeq / timePar
	return timeSeq, timePar, speedup
}

// avaliação de desempenho com diferentes tamanhos de vetor e número de processadores
func evaluate() {
	vectorSizes := []int{10_000, 100_000, 1_000_000, 10_000_000} // granularidade
	numProcessors := []int{1, 2, 4, 8, 16}                       // num de processadores

	// executa o benchmark para cada tamanho de vetor e número de processadores
	for _, size := range vectorSizes {
		fmt.Printf("Tamanho do vetor: %d\n", size)
		for _, procs := range numProcessors {
			runtime.GOMAXPROCS(procs)
			fmt.Printf("  Processadores: %d\n", procs)
			timeSeq, timePar, speedup := benchmark(size, procs)
			fmt.Printf("    Tempo sequencial: %.2fs, Tempo paralelo: %.2fs, Speedup: %.2f\n",
				timeSeq, timePar, speedup)
		}
	}
}

func main() {
	evaluate()
}
