package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var palabrasInformativas = []string{
	"RECLAMO",
	"ESTAFA",
	"COBRO",
	"GARANTIA",
	"GARANTÍA",
	"FRAUDE",
	"PRODUCTO",
	"SERVICIO",
	"SOLUCION",
	"SOLUCIÓN",
	"INCUMPLIMIENTO",
	"DENUNCIA",
}

// ESTRUCTURAS
type Record struct {
	ODI               string
	SiglasArea        string
	Anio              string
	NroExpediente     string
	TipoExpediente    string
	FechaPresentacion string
	Sector            string
	SubSector         string
	Denunciados       string
	RUC               string
	Texto             string
	Timestamp         string
	IP                string
}

type RejectedRecord struct {
	NroExpediente string
	Texto         string
	IP            string
	Motivo        string
}

// MÉTRICAS
var stats = struct {
	totalLeidos        int
	totalProcesados    int
	totalDescartados   int
	textosNormalizados int
	textosVacios       int
	filasCorruptas     int
	mu                 sync.Mutex
}{}

// Motivos
var motivosRechazo = struct {
	data map[string]int
	mu   sync.Mutex
}{
	data: make(map[string]int),
}

var advertencias = struct {
	data map[string]int
	mu   sync.Mutex
}{
	data: make(map[string]int),
}

// HELPERS
func safeGet(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return row[index]
}

func canonicalHeader(h string) string {
	h = strings.TrimPrefix(h, "\ufeff")
	h = strings.TrimSpace(h)
	h = strings.ToUpper(h)
	h = strings.ReplaceAll(h, " ", "_")
	return h
}

func getByHeader(row []string, indexByHeader map[string]int, fallback int, keys ...string) string {
	for _, key := range keys {
		idx, ok := indexByHeader[canonicalHeader(key)]
		if ok {
			return safeGet(row, idx)
		}
	}

	if fallback >= 0 {
		return safeGet(row, fallback)
	}

	return ""
}

func incrementReason(target *struct {
	data map[string]int
	mu   sync.Mutex
}, reason string) {
	target.mu.Lock()
	target.data[reason]++
	target.mu.Unlock()
}

func writeReasonBlock(file *os.File, title string, source *struct {
	data map[string]int
	mu   sync.Mutex
}) {
	file.WriteString(title + "\n")

	source.mu.Lock()
	keys := make([]string, 0, len(source.data))
	for k := range source.data {
		keys = append(keys, k)
	}
	source.mu.Unlock()

	sort.Strings(keys)

	if len(keys) == 0 {
		file.WriteString("- NINGUNA: 0\n\n")
		return
	}

	for _, k := range keys {
		source.mu.Lock()
		v := source.data[k]
		source.mu.Unlock()
		file.WriteString(fmt.Sprintf("- %s: %d\n", k, v))
	}

	file.WriteString("\n")
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

func isOnlyPunctuation(token string) bool {
	trimmed := strings.Trim(token, "!?.,;:-_*/\\|\"'`()[]{}")
	return trimmed == ""
}

func hasInformativeKeyword(text string) bool {
	for _, kw := range palabrasInformativas {
		if strings.Contains(text, kw) {
			return true
		}
	}

	return false
}

func isLowSignalText(text string, words []string) bool {
	if len(words) == 0 {
		return true
	}

	if len(words) == 1 {
		token := words[0]
		r := []rune(token)

		if len(r) <= 2 {
			return true
		}

		if isOnlyPunctuation(token) {
			return true
		}

		if allRunesEqual(token) && len(r) >= 3 {
			return true
		}
	}

	if len(words) <= 3 {
		allEqual := true
		for i := 1; i < len(words); i++ {
			if words[i] != words[0] {
				allEqual = false
				break
			}
		}
		if allEqual {
			return true
		}
	}

	return len([]rune(text)) < 8 && len(words) <= 2
}

func classifyShortTextWarning(texto string) string {
	text := strings.TrimSpace(strings.ToUpper(texto))
	if text == "" {
		return ""
	}

	words := strings.Fields(text)
	charCount := len([]rune(text))
	informative := hasInformativeKeyword(text)

	if charCount < 8 {
		if informative {
			return "TEXTO_CORTO_BAJO"
		}
		return "TEXTO_CORTO_ALTO"
	}

	if len(words) <= 2 {
		if informative {
			return "TEXTO_CORTO_BAJO"
		}

		if isLowSignalText(text, words) {
			return "TEXTO_CORTO_ALTO"
		}

		return "TEXTO_CORTO_MEDIO"
	}

	if charCount < 15 && !informative {
		return "TEXTO_CORTO_MEDIO"
	}

	return ""
}

// NORMALIZER
func normalizer(in <-chan Record, out chan<- Record) {
	for rec := range in {

		original := rec.Texto

		rec.Texto = strings.ToUpper(strings.TrimSpace(rec.Texto))
		rec.Texto = strings.Join(strings.Fields(rec.Texto), " ")

		rec.Sector = strings.ToUpper(strings.TrimSpace(rec.Sector))

		if rec.Texto != original {
			stats.mu.Lock()
			stats.textosNormalizados++
			stats.mu.Unlock()
		}

		if rec.Texto == "" || rec.Texto == "-" {
			rec.Texto = "RECLAMO VACIO"

			stats.mu.Lock()
			stats.textosVacios++
			stats.mu.Unlock()
		}

		out <- rec
	}
}

// VALIDATOR
func validator(in <-chan Record, out chan<- Record, reject chan<- RejectedRecord) {
	for rec := range in {

		texto := strings.TrimSpace(rec.Texto)
		ip := strings.TrimSpace(rec.IP)
		palabras := len(strings.Fields(texto))

		//  1. EXPEDIENTE OBLIGATORIO
		if strings.TrimSpace(rec.NroExpediente) == "" {
			stats.mu.Lock()
			stats.totalDescartados++
			stats.mu.Unlock()

			incrementReason(&motivosRechazo, "SIN_EXPEDIENTE")

			reject <- RejectedRecord{rec.NroExpediente, rec.Texto, rec.IP, "SIN_EXPEDIENTE"}
			continue
		}

		//  2. TEXTO VACÍO REAL
		if texto == "" || texto == "RECLAMO VACIO" {
			stats.mu.Lock()
			stats.totalDescartados++
			stats.mu.Unlock()

			incrementReason(&motivosRechazo, "TEXTO_VACIO")

			reject <- RejectedRecord{rec.NroExpediente, rec.Texto, rec.IP, "TEXTO_VACIO"}
			continue
		}
		
		//  3. TEXTO BASURA (REALMENTE INÚTIL)
		if palabras == 1 && len(texto) <= 3 {
			stats.mu.Lock()
			stats.totalDescartados++
			stats.mu.Unlock()

			incrementReason(&motivosRechazo, "TEXTO_BASURA")

			reject <- RejectedRecord{rec.NroExpediente, rec.Texto, rec.IP, "TEXTO_BASURA"}
			continue
		}

		//  4. TEXTO CORTO: CLASIFICAR POR SEVERIDAD (NO DESCARTAR)
		if shortWarning := classifyShortTextWarning(texto); shortWarning != "" {
			incrementReason(&advertencias, shortWarning)
		}

		//  5. IP INVÁLIDA
		if ip == "" ||
			ip == "IP_INVALIDA" ||
			strings.Contains(ip, "999") ||
			strings.Contains(ip, "abc") {

			stats.mu.Lock()
			stats.totalDescartados++
			stats.mu.Unlock()

			incrementReason(&motivosRechazo, "IP_INVALIDA")

			reject <- RejectedRecord{rec.NroExpediente, rec.Texto, rec.IP, "IP_INVALIDA"}
			continue
		}

		//  6. TIMESTAMP VACÍO → CORREGIR
		if strings.TrimSpace(rec.Timestamp) == "" {
			rec.Timestamp = "SIN_FECHA"
			incrementReason(&advertencias, "TIMESTAMP_CORREGIDO")
		}

		//  válido
		stats.mu.Lock()
		stats.totalProcesados++
		stats.mu.Unlock()

		out <- rec
	}
}

// WRITER CLEAN
func writer(in <-chan Record, done chan<- bool) {

	file, _ := os.Create("../dataset/dataset_clean.csv")
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{
		"ODI", "SIGLAS_AREA", "ANIO", "NRO_EXPEDIENTE",
		"TIPO_EXPEDIENTE", "FECHA_PRESENTACION",
		"SECTOR", "SUB_SECTOR", "DENUNCIADOS", "RUC",
		"TEXTO", "TIMESTAMP", "IP",
	})

	for rec := range in {
		w.Write([]string{
			rec.ODI,
			rec.SiglasArea,
			rec.Anio,
			rec.NroExpediente,
			rec.TipoExpediente,
			rec.FechaPresentacion,
			rec.Sector,
			rec.SubSector,
			rec.Denunciados,
			rec.RUC,
			rec.Texto,
			rec.Timestamp,
			rec.IP,
		})
	}

	done <- true
}

// WRITER REJECTED
func rejectedWriter(in <-chan RejectedRecord, done chan<- bool) {

	file, _ := os.Create("../dataset/rejected_records.csv")
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"NRO_EXPEDIENTE", "TEXTO", "IP", "MOTIVO"})

	count := 0
	limit := 10000

	for rec := range in {
		if count < limit {
			w.Write([]string{
				rec.NroExpediente,
				rec.Texto,
				rec.IP,
				rec.Motivo,
			})
			count++
		}
	}

	done <- true
}

// READER
func reader(path string, out chan<- Record) {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return
	}

	indexByHeader := make(map[string]int, len(header))
	for idx, h := range header {
		indexByHeader[canonicalHeader(h)] = idx
	}

	for {
		row, err := r.Read()
		if err != nil {
			break
		}

		stats.mu.Lock()
		stats.totalLeidos++
		stats.mu.Unlock()

		if len(row) < 13 {
			stats.mu.Lock()
			stats.filasCorruptas++
			stats.mu.Unlock()
			continue
		}

		rec := Record{
			ODI:               getByHeader(row, indexByHeader, 0, "ODI"),
			SiglasArea:        getByHeader(row, indexByHeader, 1, "SIGLAS_AREA"),
			Anio:              getByHeader(row, indexByHeader, 2, "ANIO"),
			NroExpediente:     getByHeader(row, indexByHeader, 3, "NRO_EXPEDIENTE"),
			TipoExpediente:    getByHeader(row, indexByHeader, 4, "TIPO_EXPEDIENTE"),
			FechaPresentacion: getByHeader(row, indexByHeader, 5, "FECHA_PRESENTACION"),
			Sector:            getByHeader(row, indexByHeader, 6, "SECTOR"),
			SubSector:         getByHeader(row, indexByHeader, 7, "SUB_SECTOR"),
			Denunciados:       getByHeader(row, indexByHeader, 8, "DENUNCIADOS"),
			RUC:               getByHeader(row, indexByHeader, 9, "RUC", "RUC_DENUNCIADOS"),
			Texto:             getByHeader(row, indexByHeader, 10, "TEXTO", "TEXTO_RECLAMO"),
			Timestamp:         getByHeader(row, indexByHeader, 11, "TIMESTAMP"),
			IP:                getByHeader(row, indexByHeader, 12, "IP", "IP_ADDRESS"),
		}

		out <- rec
	}

	close(out)
}

// REPORTE TXT
func generarReporte() {

	file, _ := os.Create("../dataset/reporte_limpieza.txt")
	defer file.Close()

	file.WriteString("REPORTE DE LIMPIEZA\n")
	file.WriteString("----------------------------\n")

	file.WriteString(fmt.Sprintf("Total leídos: %d\n", stats.totalLeidos))
	file.WriteString(fmt.Sprintf("Procesados: %d\n", stats.totalProcesados))
	file.WriteString(fmt.Sprintf("Descartados: %d\n", stats.totalDescartados))
	file.WriteString(fmt.Sprintf("Filas corruptas: %d\n", stats.filasCorruptas))
	file.WriteString(fmt.Sprintf("Textos normalizados: %d\n", stats.textosNormalizados))
	file.WriteString(fmt.Sprintf("Textos vacíos corregidos: %d\n\n", stats.textosVacios))

	writeReasonBlock(file, "MOTIVOS DE RECHAZO (Definitivos):", &motivosRechazo)
	writeReasonBlock(file, "ADVERTENCIAS (Procesa pero marca):", &advertencias)

	porcentaje := float64(stats.totalDescartados) / float64(stats.totalLeidos) * 100
	file.WriteString(fmt.Sprintf("\nPorcentaje descartado: %.2f%%\n", porcentaje))
}

// MAIN
func main() {

	start := time.Now()

	rawChan := make(chan Record, 200)
	normChan := make(chan Record, 200)
	validChan := make(chan Record, 200)
	rejectChan := make(chan RejectedRecord, 200)

	done := make(chan bool)
	doneReject := make(chan bool)

	// Reader
	go reader("../dataset/dataset_1M_raw.csv", rawChan)

	// Normalizers
	var wgNorm sync.WaitGroup
	wgNorm.Add(4)

	for i := 0; i < 4; i++ {
		go func() {
			defer wgNorm.Done()
			normalizer(rawChan, normChan)
		}()
	}

	go func() {
		wgNorm.Wait()
		close(normChan)
	}()

	// Validators
	var wgVal sync.WaitGroup
	wgVal.Add(4)

	for i := 0; i < 4; i++ {
		go func() {
			defer wgVal.Done()
			validator(normChan, validChan, rejectChan)
		}()
	}

	go func() {
		wgVal.Wait()
		close(validChan)
		close(rejectChan)
	}()

	// Writers
	go writer(validChan, done)
	go rejectedWriter(rejectChan, doneReject)

	<-done
	<-doneReject

	// Reporte
	generarReporte()

	fmt.Println("\n REPORTE DE LIMPIEZA")
	fmt.Println("-----------------------------")
	fmt.Println("Total leídos:        ", stats.totalLeidos)
	fmt.Println("Procesados:          ", stats.totalProcesados)
	fmt.Println("Descartados:         ", stats.totalDescartados)
	fmt.Println("Filas corruptas:     ", stats.filasCorruptas)
	fmt.Println("Textos normalizados: ", stats.textosNormalizados)
	fmt.Println("Textos vacíos:       ", stats.textosVacios)

	fmt.Printf("Porcentaje descartado: %.2f%%\n",
		float64(stats.totalDescartados)/float64(stats.totalLeidos)*100)

	fmt.Println("\n⏱ Tiempo total:", time.Since(start))

	fmt.Println("\n Archivos generados:")
	fmt.Println("- dataset_clean.csv")
	fmt.Println("- rejected_records.csv")
	fmt.Println("- reporte_limpieza.txt")

}
