package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// =====================
// EJECUCIÓN DE PROGRAMAS
// =====================
func runProgram(name string, args ...string) (time.Duration, string) {

	start := time.Now()

	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	duration := time.Since(start)

	return duration, string(output)
}

// =====================
// EXTRAER MEMORIA
// =====================
func extractMemoryMB(output string) float64 {

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, "Memoria usada") {

			// Ej: "Memoria usada: 736.40 MB"
			parts := strings.Fields(line)

			for i, p := range parts {
				if p == "usada:" && i+1 < len(parts) {
					value, err := strconv.ParseFloat(parts[i+1], 64)
					if err == nil {
						return value
					}
				}
			}
		}
	}

	return 0.0
}

// =====================
// MEDIA RECORTADA (TIEMPO)
// =====================
func trimmedMeanDurations(values []time.Duration) time.Duration {

	if len(values) <= 2 {
		return 0
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	var sum time.Duration

	for i := 1; i < len(values)-1; i++ {
		sum += values[i]
	}

	return sum / time.Duration(len(values)-2)
}

// =====================
// MEDIA RECORTADA (MEMORIA)
// =====================
func trimmedMeanFloat(values []float64) float64 {

	if len(values) <= 2 {
		return 0
	}

	sort.Float64s(values)

	var sum float64

	for i := 1; i < len(values)-1; i++ {
		sum += values[i]
	}

	return sum / float64(len(values)-2)
}

// =====================
// MAIN
// =====================
func main() {

	runs := 10

	var seqTimes []time.Duration
	var concTimes []time.Duration

	var seqMems []float64
	var concMems []float64

	fmt.Println("🚀 INICIANDO BENCHMARK\n")

	for i := 1; i <= runs; i++ {

		fmt.Printf("🔹 Iteración %d\n", i)

		// 🔹 Secuencial
		seqTime, seqOutput := runProgram("go", "run", "secuencial.go")
		seqMem := extractMemoryMB(seqOutput)

		fmt.Printf("Secuencial: %v | Memoria: %.2f MB\n", seqTime, seqMem)

		// 🔹 Concurrente
		concTime, concOutput := runProgram("go", "run", "concurrente.go")
		concMem := extractMemoryMB(concOutput)

		fmt.Printf("Concurrente: %v | Memoria: %.2f MB\n", concTime, concMem)

		fmt.Println("-----------------------------")

		// Guardar para media recortada
		seqTimes = append(seqTimes, seqTime)
		seqMems = append(seqMems, seqMem)

		concTimes = append(concTimes, concTime)
		concMems = append(concMems, concMem)
	}

	// =====================
	// RESULTADOS
	// =====================

	trimSeq := trimmedMeanDurations(seqTimes)
	trimConc := trimmedMeanDurations(concTimes)

	trimMemSeq := trimmedMeanFloat(seqMems)
	trimMemConc := trimmedMeanFloat(concMems)

	speedup := float64(trimSeq) / float64(trimConc)

	fmt.Println("\n📊 RESULTADOS (MEDIA RECORTADA)")
	fmt.Println("-----------------------------")
	fmt.Printf("Secuencial: %v | Memoria: %.2f MB\n", trimSeq, trimMemSeq)
	fmt.Printf("Concurrente: %v | Memoria: %.2f MB\n", trimConc, trimMemConc)
	fmt.Printf("Speedup: %.2f\n", speedup)
}
