package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Função de merge sort
func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])
	return merge(left, right)
}

// Função para mesclar duas metades ordenadas
func merge(left, right []int) []int {
	result := []int{}
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
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

// Função para merge sort paralelo
func mergeSortParallel(arr []int, nProc int) []int {
	var wg sync.WaitGroup

	// Dividir o arr em nProc partes
	chunkSize := len(arr) / nProc
	chunks := make([][]int, nProc)

	for i := 0; i < nProc; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == nProc-1 {
			end = len(arr) // Último pedaço
		}
		chunks[i] = arr[start:end]
	}

	// Executar o merge sort em paralelo
	for i := 0; i < nProc; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			chunks[i] = mergeSort(chunks[i])
		}(i)
	}

	wg.Wait()

	// Mesclar as partes ordenadas
	for len(chunks) > 1 {
		var temp [][]int
		for i := 0; i < len(chunks); i += 2 {
			if i+1 < len(chunks) {
				temp = append(temp, merge(chunks[i], chunks[i+1]))
			} else {
				temp = append(temp, chunks[i])
			}
		}
		chunks = temp
	}

	return chunks[0]
}

func main() {
	// Exemplo de tamanhos de vetores
	sizes := []int{1000, 10000, 100000, 1000000, 10000000}

	for _, size := range sizes {
		// Gerar vetor aleatório
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = rand.Intn(10000)
		}

		// Rodar o merge sort sequencial para calcular o tempo de referência
		var sequentialTime time.Duration
		iterations := 5
		for i := 0; i < iterations; i++ {
			start := time.Now()
			_ = mergeSort(arr)
			sequentialTime += time.Since(start)
		}
		sequentialTime /= time.Duration(iterations) // Média do tempo sequencial

		// Testar para diferentes números de processadores
		for nProc := 1; nProc <= 8; nProc++ {
			var parallelTime time.Duration
			for i := 0; i < iterations; i++ {
				start := time.Now()
				_ = mergeSortParallel(arr, nProc)
				parallelTime += time.Since(start)
			}
			parallelTime /= time.Duration(iterations) // Média do tempo paralelo

			// Verificar se o tempo paralelo não é muito pequeno
			if parallelTime > 0 {
				speedup := float64(sequentialTime) / float64(parallelTime)
				// Printar os resultados
				fmt.Printf("Tamanho: %d, Processadores: %d, Speedup: %.2f\n", size, nProc, speedup)
			} else {
				// Caso não tenha ocorrido paralelização eficaz, exibir uma mensagem
				fmt.Printf("Tamanho: %d, Processadores: %d, Speedup: Não calculado\n", size, nProc)
			}
		}
	}
}
