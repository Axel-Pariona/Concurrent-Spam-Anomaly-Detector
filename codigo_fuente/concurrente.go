package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Record struct {
	Data []string
	Text string
	IP   string
}

// FRECUENCIAS (con mutex)
var textFreq = struct {
	data map[string]int
	mu   sync.Mutex
}{
	data: make(map[string]int),
}

var ipFreq = struct {
	data map[string]int
	mu   sync.Mutex
}{
	data: make(map[string]int),
}

// HELPERS
func allRunesEqual(s string) bool {
	r := []rune(s)
	if len(r) <= 1 {
		return true
	}
	for i := 1; i < len(r); i++ {
		if r[i] != r[0] {
			return false
		}
	}
	return true
}

// WORKER FASE 1 (conteo)
func counterWorker(in <-chan Record, wg *sync.WaitGroup) {
	defer wg.Done()

	for rec := range in {

		textFreq.mu.Lock()
		textFreq.data[rec.Text]++
		textFreq.mu.Unlock()

		ipFreq.mu.Lock()
		ipFreq.data[rec.IP]++
		ipFreq.mu.Unlock()
	}
}

// WORKER FASE 2 
func classifierWorker(
	in <-chan Record,
	clean chan<- []string,
	spam chan<- []string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for rec := range in {

		isSpam := false

		// Regla 1: BOT
		if textFreq.data[rec.Text] > 100 && ipFreq.data[rec.IP] > 20 {
			isSpam = true
		}

		// Regla 2: IP abusiva
		if ipFreq.data[rec.IP] > 200 {
			isSpam = true
		}

		// Regla 3: texto artificial
		if allRunesEqual(rec.Text) && len(rec.Text) > 10 {
			isSpam = true
		}

		if isSpam {
			spam <- rec.Data
		} else {
			clean <- rec.Data
		}
	}
}

// MAIN
func main() {

	start := time.Now()

	file, err := os.Open("../dataset/dataset_clean.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	header, _ := reader.Read()

	// Salida
	cleanFile, _ := os.Create("../dataset/dataset_final_concurrent.csv")
	spamFile, _ := os.Create("../dataset/spam_detected_concurrent.csv")

	defer cleanFile.Close()
	defer spamFile.Close()

	cleanWriter := csv.NewWriter(cleanFile)
	spamWriter := csv.NewWriter(spamFile)

	defer cleanWriter.Flush()
	defer spamWriter.Flush()

	cleanWriter.Write(header)
	spamWriter.Write(header)

	// FASE 1: LECTURA + CONTEO
	recordChan := make(chan Record, 1000)

	var wgCounter sync.WaitGroup
	numWorkers := 4

	for i := 0; i < numWorkers; i++ {
		wgCounter.Add(1)
		go counterWorker(recordChan, &wgCounter)
	}

	// Leer CSV
	var allRecords []Record

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		text := strings.TrimSpace(strings.ToUpper(row[10]))
		ip := strings.TrimSpace(row[12])

		rec := Record{
			Data: row,
			Text: text,
			IP:   ip,
		}

		allRecords = append(allRecords, rec)
		recordChan <- rec
	}

	close(recordChan)
	wgCounter.Wait()

	// FASE 2: CLASIFICACIÓN
	classifyChan := make(chan Record, 1000)
	cleanChan := make(chan []string, 1000)
	spamChan := make(chan []string, 1000)

	var wgClassifier sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wgClassifier.Add(1)
		go classifierWorker(classifyChan, cleanChan, spamChan, &wgClassifier)
	}

	// Writer goroutines
	var wgWriter sync.WaitGroup
	wgWriter.Add(2)

	go func() {
		defer wgWriter.Done()
		for row := range cleanChan {
			cleanWriter.Write(row)
		}
	}()

	go func() {
		defer wgWriter.Done()
		for row := range spamChan {
			spamWriter.Write(row)
		}
	}()

	// Enviar datos a clasificación
	for _, rec := range allRecords {
		classifyChan <- rec
	}

	close(classifyChan)
	wgClassifier.Wait()

	close(cleanChan)
	close(spamChan)
	wgWriter.Wait()

	// RESULTADOS
	
	fmt.Println("\n DETECCIÓN DE SPAM (CONCURRENTE)")
	fmt.Println("--------------------------------")
	fmt.Println("Total registros:", len(allRecords))
	fmt.Println("Tiempo:", time.Since(start))

	fmt.Println("\n Archivos generados:")
	fmt.Println("- dataset_final_concurrent.csv")
	fmt.Println("- spam_detected_concurrent.csv")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Memoria usada: %.2f MB\n", float64(m.Alloc)/1024/1024)
}
