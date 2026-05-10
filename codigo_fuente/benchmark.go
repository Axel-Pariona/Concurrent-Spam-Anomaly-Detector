package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// =====================================
// EJECUTAR PROGRAMA
// =====================================

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

// =====================================
// EXTRAER MEMORIA
// =====================================

func extractMemoryMB(output string) float64 {

	lines := strings.Split(output, "\n")

	for _, line := range lines {

		if strings.Contains(line, "Memoria usada") {

			parts := strings.Fields(line)

			for i, p := range parts {

				if p == "usada:" && i+1 < len(parts) {

					value, err := strconv.ParseFloat(
						parts[i+1],
						64,
					)

					if err == nil {
						return value
					}
				}
			}
		}
	}

	return 0
}

// =====================================
// MEDIA RECORTADA TIEMPO
// =====================================

func trimmedMeanDurations(
	values []time.Duration,
) time.Duration {

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

// =====================================
// MEDIA RECORTADA FLOAT
// =====================================

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

// =====================================
// MAIN
// =====================================

func main() {

	runs := 5

	var seqTimes []time.Duration
	var concTimes []time.Duration

	var seqMems []float64
	var concMems []float64

	fmt.Println("===================================")
	fmt.Println("BENCHMARK SPAM DETECTOR")
	fmt.Println("===================================")

	// =================================
	// SECUENCIAL
	// =================================

	fmt.Println("\nEJECUTANDO SECUENCIAL...\n")

	for i := 1; i <= runs; i++ {

		fmt.Printf("Secuencial Run %d\n", i)

		duration, output := runProgram(
			"go",
			"run",
			"secuencial.go",
		)

		mem := extractMemoryMB(output)

		fmt.Printf(
			"Tiempo: %v | Memoria: %.2f MB\n",
			duration,
			mem,
		)

		fmt.Println("-------------------------")

		seqTimes = append(seqTimes, duration)
		seqMems = append(seqMems, mem)
	}

	// =================================
	// CONCURRENTE
	// =================================

	fmt.Println("\nEJECUTANDO CONCURRENTE...\n")

	for i := 1; i <= runs; i++ {

		fmt.Printf("Concurrente Run %d\n", i)

		duration, output := runProgram(
			"go",
			"run",
			"concurrente.go",
		)

		mem := extractMemoryMB(output)

		fmt.Printf(
			"Tiempo: %v | Memoria: %.2f MB\n",
			duration,
			mem,
		)

		fmt.Println("-------------------------")

		concTimes = append(concTimes, duration)
		concMems = append(concMems, mem)
	}

	// =================================
	// MEDIA RECORTADA
	// =================================

	seqAvg := trimmedMeanDurations(seqTimes)
	concAvg := trimmedMeanDurations(concTimes)

	seqMemAvg := trimmedMeanFloat(seqMems)
	concMemAvg := trimmedMeanFloat(concMems)

	// =================================
	// SPEEDUP
	// =================================

	speedup := float64(seqAvg) / float64(concAvg)

	// =================================
	// RESULTADOS
	// =================================

	fmt.Println("\n===================================")
	fmt.Println("RESULTADOS FINALES")
	fmt.Println("MEDIA RECORTADA")
	fmt.Println("===================================")

	fmt.Printf(
		"\nSecuencial:\nTiempo promedio: %v\nMemoria promedio: %.2f MB\n",
		seqAvg,
		seqMemAvg,
	)

	fmt.Printf(
		"\nConcurrente:\nTiempo promedio: %v\nMemoria promedio: %.2f MB\n",
		concAvg,
		concMemAvg,
	)

	fmt.Println("\n===================================")

	fmt.Printf("SPEEDUP: %.2fx\n", speedup)

	fmt.Println("===================================")
}
