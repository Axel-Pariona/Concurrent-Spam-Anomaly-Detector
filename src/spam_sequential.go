package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// =====================================
// ESTRUCTURA
// =====================================

type Record struct {
	Data []string

	Text string
	IP   string

	Score   int
	Reasons []string
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

// -------------------------------------
// PALABRAS SPAM
// -------------------------------------

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

// -------------------------------------
// CAPS RATIO
// -------------------------------------

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

// -------------------------------------
// DIVERSIDAD LEXICA
// -------------------------------------

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

// -------------------------------------
// REPETICION EXCESIVA
// -------------------------------------

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

// -------------------------------------
// TEXTO ARTIFICIAL
// -------------------------------------

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

// -------------------------------------
// MUCHOS SIMBOLOS
// -------------------------------------

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

// -------------------------------------
// CONECTORES
// -------------------------------------

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

	// =================================
	// 1. TEXTO MUY REPETIDO
	// =================================

	if textFreq[rec.Text] > 50 {

		score += 3

		reasons = append(
			reasons,
			"TEXTO_REPETIDO",
		)
	}

	// =================================
	// 2. IP ABUSIVA
	// =================================

	if ipFreq[rec.IP] > 100 {

		score += 4

		reasons = append(
			reasons,
			"IP_ABUSIVA",
		)
	}

	// =================================
	// 3. MUCHAS MAYUSCULAS
	// =================================

	if capsRatio(rec.Text) > 0.75 {

		score += 2

		reasons = append(
			reasons,
			"EXCESO_MAYUSCULAS",
		)
	}

	// =================================
	// 4. BAJA DIVERSIDAD LEXICA
	// =================================

	if lexicalDiversity(rec.Text) < 0.45 {

		score += 2

		reasons = append(
			reasons,
			"BAJA_DIVERSIDAD",
		)
	}

	// =================================
	// 5. PALABRAS SPAM
	// =================================

	spamWordsFound :=
		containsSpamWords(rec.Text)

	if spamWordsFound > 0 {

		score += spamWordsFound

		reasons = append(
			reasons,
			"PALABRAS_SPAM",
		)
	}

	// =================================
	// 6. TEXTO MUY CORTO
	// =================================

	if len(words) <= 2 {

		score += 2

		reasons = append(
			reasons,
			"TEXTO_CORTO",
		)
	}

	// =================================
	// 7. REPETICION EXCESIVA
	// =================================

	if repeatedWords(rec.Text) {

		score += 3

		reasons = append(
			reasons,
			"REPETICION_EXCESIVA",
		)
	}

	// =================================
	// 8. TEXTO ARTIFICIAL
	// =================================

	if allRunesEqual(rec.Text) &&
		len(rec.Text) > 8 {

		score += 5

		reasons = append(
			reasons,
			"TEXTO_ARTIFICIAL",
		)
	}

	// =================================
	// 9. EXCESO DE SIMBOLOS
	// =================================

	if demasiadosSimbolos(rec.Text) {

		score += 2

		reasons = append(
			reasons,
			"EXCESO_SIMBOLOS",
		)
	}

	// =================================
	// 10. SIN CONECTORES
	// =================================

	if !hasConnectors(rec.Text) &&
		len(words) > 5 {

		score += 2

		reasons = append(
			reasons,
			"SIN_CONECTORES",
		)
	}

	// =================================
	// 11. SPAM DIRECTO
	// =================================

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
// MAIN
// =====================================

func main() {

	start := time.Now()

	// =================================
	// ABRIR DATASET
	// =================================

	file, err := os.Open(
		"../data/dataset_clean.csv",
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	reader.FieldsPerRecord = -1

	header, _ := reader.Read()

	// =================================
	// OUTPUTS
	// =================================

	cleanFile, _ := os.Create(
		"../data/dataset_final_secuencial.csv",
	)

	spamFile, _ := os.Create(
		"../data/spam_detected_secuencial.csv",
	)

	defer cleanFile.Close()
	defer spamFile.Close()

	cleanWriter := csv.NewWriter(cleanFile)

	spamWriter := csv.NewWriter(spamFile)

	defer cleanWriter.Flush()
	defer spamWriter.Flush()

	// =================================
	// HEADERS
	// =================================

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

	// =================================
	// RECORDS
	// =================================

	var records []Record

	// =================================
	// METRICAS
	// =================================

	total := 0
	spamCount := 0
	cleanCount := 0

	// =================================
	// LECTURA
	// =================================

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
	// CLASIFICACION
	// =================================

	for _, rec := range records {

		score, reasons := spamScore(
			rec,
			textFreq,
			ipFreq,
		)

		isSpam := false

		if score >= 6 {
			isSpam = true
		}

		rowCopy := append(
			[]string{},
			rec.Data...,
		)

		newRow := append(
			rowCopy,
			fmt.Sprintf("%d", score),
			strings.Join(reasons, "|"),
		)

		if isSpam {

			spamWriter.Write(newRow)

			spamCount++

		} else {

			cleanWriter.Write(newRow)

			cleanCount++
		}
	}

	// =================================
	// REPORTE
	// =================================

	fmt.Println(
		"\n===== DETECCION DE SPAM (SECUENCIAL) =====",
	)

	fmt.Println(
		"\nTotal registros:",
		total,
	)

	fmt.Println(
		"Spam detectado:",
		spamCount,
	)

	fmt.Println(
		"Registros validos:",
		cleanCount,
	)

	fmt.Printf(
		"Porcentaje spam: %.2f%%\n",
		float64(spamCount)/
			float64(total)*100,
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

	fmt.Println(
		"- dataset_final_secuencial.csv",
	)

	fmt.Println(
		"- spam_detected_secuencial.csv",
	)
}
