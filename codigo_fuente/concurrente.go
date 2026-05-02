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

func main() {

	start := time.Now()

	numWorkers := runtime.NumCPU() / 2
	if numWorkers < 2 {
		numWorkers = 2
	}

	// FASE 1: LECTURA + CONTEO (SECUENCIAL)
	file, _ := os.Open("../dataset/dataset_clean.csv")
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	header, _ := reader.Read()

	textFreq := make(map[string]int)
	ipFreq := make(map[string]int)

	var records []Record

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		text := strings.TrimSpace(strings.ToUpper(row[10]))
		ip := strings.TrimSpace(row[12])

		textFreq[text]++
		ipFreq[ip]++

		records = append(records, Record{
			Data: row,
			Text: text,
			IP:   ip,
		})
	}

	file.Close()

	// FASE 2: PIPELINE CONCURRENTE
	inputChan := make(chan Record, 200)
	cleanChan := make(chan []string, 200)
	spamChan := make(chan []string, 200)

	var wg sync.WaitGroup

	// WORKERS
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for rec := range inputChan {

				isSpam := false

				if textFreq[rec.Text] > 100 && ipFreq[rec.IP] > 20 {
					isSpam = true
				}

				if ipFreq[rec.IP] > 200 {
					isSpam = true
				}

				if allRunesEqual(rec.Text) && len(rec.Text) > 10 {
					isSpam = true
				}

				if isSpam {
					spamChan <- rec.Data
				} else {
					cleanChan <- rec.Data
				}
			}
		}()
	}

	// PRODUCER
	go func() {
		for _, rec := range records {
			inputChan <- rec
		}
		close(inputChan)
	}()

	// CIERRE CONTROLADO
	go func() {
		wg.Wait()
		close(cleanChan)
		close(spamChan)
	}()

	// FASE 3: ESCRITURA CONCURRENTE
	cleanFile, _ := os.Create("../dataset/dataset_final_concurrent.csv")
	spamFile, _ := os.Create("../dataset/spam_detected_concurrent.csv")

	cleanWriter := csv.NewWriter(cleanFile)
	spamWriter := csv.NewWriter(spamFile)

	cleanWriter.Write(header)
	spamWriter.Write(header)

	var wgWriter sync.WaitGroup
	wgWriter.Add(2)

	// consumidor CLEAN
	go func() {
		defer wgWriter.Done()
		for row := range cleanChan {
			cleanWriter.Write(row)
		}
	}()

	// consumidor SPAM
	go func() {
		defer wgWriter.Done()
		for row := range spamChan {
			spamWriter.Write(row)
		}
	}()

	wgWriter.Wait()

	cleanWriter.Flush()
	spamWriter.Flush()

	cleanFile.Close()
	spamFile.Close()

	// RESULTADOS
	fmt.Println("\nDETECCIÓN DE SPAM ")
	fmt.Println("----------------------------------------------------")
	fmt.Println("Workers:", numWorkers)
	fmt.Println("Tiempo:", time.Since(start))

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Memoria usada: %.2f MB\n", float64(m.Alloc)/1024/1024)
}
