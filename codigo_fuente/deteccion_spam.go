package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
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

	file, err := os.Open("../dataset/dataset_clean.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	header, _ := reader.Read()

	// Archivos de salida
	cleanFile, _ := os.Create("../dataset/dataset_final.csv")
	spamFile, _ := os.Create("../dataset/spam_detected.csv")

	defer cleanFile.Close()
	defer spamFile.Close()

	cleanWriter := csv.NewWriter(cleanFile)
	spamWriter := csv.NewWriter(spamFile)

	defer cleanWriter.Flush()
	defer spamWriter.Flush()

	cleanWriter.Write(header)
	spamWriter.Write(header)

	// Frecuencias
	textFreq := make(map[string]int)
	ipFreq := make(map[string]int)

	var total, spamCount, cleanCount int

	// Leer todo primero para contar frecuencias
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

		total++
	}

	// Clasificación
	for _, rec := range records {

		isSpam := false

		// ✅ Regla 1: BOT (texto repetido + misma IP)
		if textFreq[rec.Text] > 100 && ipFreq[rec.IP] > 20 {
			isSpam = true
		}

		// ✅ Regla 2: IP abusiva
		if ipFreq[rec.IP] > 200 {
			isSpam = true
		}

		// ✅ Regla 3: texto artificial
		if allRunesEqual(rec.Text) && len(rec.Text) > 10 {
			isSpam = true
		}

		if isSpam {
			spamWriter.Write(rec.Data)
			spamCount++
		} else {
			cleanWriter.Write(rec.Data)
			cleanCount++
		}
	}

	// Resultados
	fmt.Println("\n📊 DETECCIÓN DE SPAM (SECUENCIAL)")
	fmt.Println("--------------------------------")
	fmt.Println("Total registros:", total)
	fmt.Println("Spam detectado:", spamCount)
	fmt.Println("Registros válidos:", cleanCount)
	fmt.Println("Tiempo:", time.Since(start))

	fmt.Println("\n📁 Archivos generados:")
	fmt.Println("- dataset_final.csv")
	fmt.Println("- spam_detected.csv")
}
