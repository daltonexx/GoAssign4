package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	// Configurações
	casas := [6]int{999, 999999, 999999999, 999999999999, 999999999999999, 999999999999999999}
	N := 50 // Número de primos a gerar
	processadores := []int{2, 4, 6, 8, 10, 12}

	// Itera para cada tamanho de primo
	for _, tam := range casas {
		fmt.Printf("Tamanho: %d\n", tam)
		for _, P := range processadores {
			speedup := calcularSpeedup(N, tam, P)
			fmt.Printf("Processadores: %d, Speedup: %.2f\n", P, speedup)
		}
	}
}

// Função que calcula o speedup
func calcularSpeedup(N, tam, P int) float64 {
	// Tempo sequencial
	runtime.GOMAXPROCS(1) // Força execução com 1 processador
	sequencial := medirTempoSequencial(N, tam)

	// Tempo paralelo
	runtime.GOMAXPROCS(P) // Configura o número de processadores
	paralelo := medirTempoParalelo(N, tam, P)

	// Calcula e retorna o speedup
	return float64(sequencial) / float64(paralelo)
}

// Medir tempo na execução sequencial
func medirTempoSequencial(N, tam int) time.Duration {
	start := time.Now()
	for i := 0; i < N; i++ {
		genPrime(tam)
	}
	return time.Since(start)
}

// Medir tempo na execução paralela
func medirTempoParalelo(N, tam, P int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup
	tarefas := make(chan int, N)

	// Envia N tarefas
	for i := 0; i < N; i++ {
		tarefas <- tam
	}
	close(tarefas)

	// Cria P workers
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for tam := range tarefas {
				genPrime(tam)
			}
		}()
	}

	// Aguarda conclusão
	wg.Wait()
	return time.Since(start)
}

// Função que gera um número primo
func genPrime(tam int) {
	notPrimo := true
	for notPrimo {
		v := rand.Intn(tam)
		notPrimo = !isPrime(v)
	}
}

// Verifica se um número é primo
func isPrime(p int) bool {
	if p < 2 || p%2 == 0 {
		return false
	}
	for i := 3; i*i <= p; i += 2 {
		if p%i == 0 {
			return false
		}
	}
	return true
}
