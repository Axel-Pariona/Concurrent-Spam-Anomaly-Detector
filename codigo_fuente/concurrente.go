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

// =====================================
// ESTRUCTURAS
// =====================================

type Record struct {
	Data []string
	Text string
	IP   string
}

type Result struct {
	Data    []string
	IsSpam  bool
	Score   int
	Reasons string
}

// =====================================
// PALABRAS SPAM
// =====================================

var spamWords = []string{
	"URGENTE",
	"ESTAFA",
	"DINERO",
	"GRATIS",
	"PREMIO",
	"AHORA",
	"YA",
	"CLICK",
	"SOLUCION",
	"EXIJO",
	"RESPUESTA",
	"FRAUDE",
	"OFERTA",
	"GANA",
}

// =====================================
// CONECTORES LEGITIMOS
// =====================================

var validConnectors = []string{
	"PORQUE",
	"ADEMAS",
	"SIN EMBARGO",
	"AUNQUE",
	"DEBIDO",
	"YA QUE",
	"POR FAVOR",
	"MIENTRAS",
	"TAMBIEN",
	"ENTONCES",
}

// =====================================
// HELPERS
// =====================================

func containsSpamWords(text string) int {

	score := 0

	upper := strings.ToUpper(text)

	for _, word := range spamWords {

		if strings.Contains(upper, word) {
			score++
		}
	}

	return score
}

func capsRatio(text string) float64 {

	letters := 0
	upper := 0

	for _, c := range text {

		if c >= 'A' && c <= 'Z' {

			letters++
			upper++

		} else if c >= 'a' && c <= 'z' {

			letters++
		}
	}

	if letters == 0 {
		return 0
	}

	return float64(upper) / float64(letters)
}

func lexicalDiversity(text string) float64 {

	words := strings.Fields(
		strings.ToUpper(text),
	)

	if len(words) == 0 {
		return 0
	}

	unique := make(map[string]bool)

	for _, w := range words {
		unique[w] = true
	}

	return float64(len(unique)) /
		float64(len(words))
}

func repeatedWords(text string) bool {

	words := strings.Fields(
		strings.ToUpper(text),
	)

	if len(words) < 4 {
		return false
	}

	freq := make(map[string]int)

	for _, w := range words {
		freq[w]++
	}

	for _, v := range freq {

		if v >= 4 {
			return true
		}
	}

	return false
}

func allRunesEqual(s string) bool {

	r := []rune(s)

	if len(r) <= 1 {
		return false
	}

	for i := 1; i < len(r); i++ {

		if r[i] != r[0] {
			return false
		}
	}

	return true
}

func demasiadosSimbolos(text string) bool {

	count := 0

	for _, c := range text {

		if strings.ContainsRune(
			"!@#$%^&*........",
			c,
		) {
			count++
		}
	}

	return count > 8
}

func hasConnectors(text string) bool {

	upper := strings.ToUpper(text)

	for _, c := range validConnectors {

		if strings.Contains(upper, c) {
			return true
		}
	}

	return false
}

// =====================================
// SCORING NLP SIMPLE
// =====================================

func spamScore(
	rec Record,
	textFreq map[string]int,
	ipFreq map[string]int,
) (int, []string) {

	score := 0

	reasons := []string{}

	words := strings.Fields(rec.Text)

	// TEXTO MUY REPETIDO
	if textFreq[rec.Text] > 50 {

		score += 3

		reasons = append(
			reasons,
			"TEXTO_REPETIDO",
		)
	}

	// IP ABUSIVA
	if ipFreq[rec.IP] > 100 {

		score += 4

		reasons = append(
			reasons,
			"IP_ABUSIVA",
		)
	}

	// MUCHAS MAYUSCULAS
	if capsRatio(rec.Text) > 0.75 {

		score += 2

		reasons = append(
			reasons,
			"EXCESO_MAYUSCULAS",
		)
	}

	// BAJA DIVERSIDAD
	if lexicalDiversity(rec.Text) < 0.45 {

		score += 2

		reasons = append(
			reasons,
			"BAJA_DIVERSIDAD",
		)
	}

	// PALABRAS SPAM
	spamWordsFound :=
		containsSpamWords(rec.Text)

	if spamWordsFound > 0 {

		score += spamWordsFound

		reasons = append(
			reasons,
			"PALABRAS_SPAM",
		)
	}

	// TEXTO CORTO
	if len(words) <= 2 {

		score += 2

		reasons = append(
			reasons,
			"TEXTO_CORTO",
		)
	}

	// REPETICION EXCESIVA
	if repeatedWords(rec.Text) {

		score += 3

		reasons = append(
			reasons,
			"REPETICION_EXCESIVA",
		)
	}

	// TEXTO ARTIFICIAL
	if allRunesEqual(rec.Text) &&
		len(rec.Text) > 8 {

		score += 5

		reasons = append(
			reasons,
			"TEXTO_ARTIFICIAL",
		)
	}

	// EXCESO DE SIMBOLOS
	if demasiadosSimbolos(rec.Text) {

		score += 2

		reasons = append(
			reasons,
			"EXCESO_SIMBOLOS",
		)
	}

	// SIN CONECTORES
	if !hasConnectors(rec.Text) &&
		len(words) > 5 {

		score += 2

		reasons = append(
			reasons,
			"SIN_CONECTORES",
		)
	}

	// SPAM DIRECTO
	if spamWordsFound >= 2 &&
		len(words) <= 4 {

		score += 3

		reasons = append(
			reasons,
			"SPAM_DIRECTO",
		)
	}

	return score, reasons
}

// =====================================
// WORKER
// =====================================

func worker(
	input <-chan Record,
	output chan<- Result,
	textFreq map[string]int,
	ipFreq map[string]int,
	wg *sync.WaitGroup,
) {

	defer wg.Done()

	for rec := range input {

		score, reasons := spamScore(
			rec,
			textFreq,
			ipFreq,
		)

		isSpam := score >= 6

		rowCopy := append(
			[]string{},
			rec.Data...,
		)

		newRow := append(
			rowCopy,
			fmt.Sprintf("%d", score),
			strings.Join(reasons, "|"),
		)

		output <- Result{
			Data:    newRow,
			IsSpam:  isSpam,
			Score:   score,
			Reasons: strings.Join(reasons, "|"),
		}
	}
}

// =====================================
// MAIN
// =====================================

func main() {

	start := time.Now()

	numWorkers := 16 //runtime.NumCPU()

	fmt.Println("Workers:", numWorkers)

	file, err := os.Open(
		"../dataset/dataset_clean.csv",
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	header, _ := reader.Read()

	cleanFile, _ := os.Create(
		"../dataset/dataset_final_concurrente.csv",
	)

	spamFile, _ := os.Create(
		"../dataset/spam_detected_concurrente.csv",
	)

	defer cleanFile.Close()
	defer spamFile.Close()

	cleanWriter := csv.NewWriter(cleanFile)
	spamWriter := csv.NewWriter(spamFile)

	defer cleanWriter.Flush()
	defer spamWriter.Flush()

	newHeader := append(
		header,
		"SPAM_SCORE",
		"SPAM_REASON",
	)

	cleanWriter.Write(newHeader)
	spamWriter.Write(newHeader)

	// =================================
	// FRECUENCIAS
	// =================================

	textFreq := make(map[string]int)
	ipFreq := make(map[string]int)

	var records []Record

	total := 0

	for {

		row, err := reader.Read()

		if err != nil {
			break
		}

		if len(row) < 5 {
			continue
		}

		text := strings.TrimSpace(row[4])

		ip := strings.TrimSpace(row[2])

		textFreq[text]++
		ipFreq[ip]++

		records = append(records, Record{
			Data: row,
			Text: text,
			IP:   ip,
		})

		total++
	}

	// =================================
	// CHANNELS
	// =================================

	inputChan := make(chan Record, 1000)
	outputChan := make(chan Result, 1000)

	// =================================
	// WORKERS
	// =================================

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {

		wg.Add(1)

		go worker(
			inputChan,
			outputChan,
			textFreq,
			ipFreq,
			&wg,
		)
	}

	// =================================
	// PRODUCER
	// =================================

	go func() {

		for _, rec := range records {
			inputChan <- rec
		}

		close(inputChan)
	}()

	// =================================
	// CIERRE OUTPUT
	// =================================

	go func() {

		wg.Wait()

		close(outputChan)
	}()

	// =================================
	// CONSUMER
	// =================================

	spamCount := 0
	cleanCount := 0

	for result := range outputChan {

		if result.IsSpam {

			spamWriter.Write(result.Data)
			spamCount++

		} else {

			cleanWriter.Write(result.Data)
			cleanCount++
		}
	}

	// =================================
	// REPORTE
	// =================================

	fmt.Println(
		"\n===== DETECCION DE SPAM CONCURRENTE =====",
	)

	fmt.Println("\nTotal registros:", total)

	fmt.Println("Spam detectado:", spamCount)

	fmt.Println("Registros validos:", cleanCount)

	fmt.Printf(
		"Porcentaje spam: %.2f%%\n",
		float64(spamCount)/float64(total)*100,
	)

	fmt.Println(
		"\nTiempo total:",
		time.Since(start),
	)

	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	fmt.Printf(
		"Memoria usada: %.2f MB\n",
		float64(m.Alloc)/1024/1024,
	)

	fmt.Println("\nArchivos generados:")
	fmt.Println("- dataset_final_concurrente.csv")
	fmt.Println("- spam_detected_concurrente.csv")
}
