package main

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// =====================================
// COLUMNAS DEL DATASET ORIGINAL
// =====================================

const (
	COL_TIMESTAMP = 10
	COL_TEXTO     = 11
	COL_USUARIO   = 12
	COL_IP        = 13
	COL_ID        = 14
)

// =====================================
// ESTRUCTURAS
// =====================================

type Record struct {
	IDReclamo string
	UsuarioID string
	IP        string
	Timestamp string

	TextoOriginal string
	TextoLimpio   string
}

type RejectedRecord struct {
	Record
	Motivo string
}

// =====================================
// METRICAS GLOBALES
// =====================================

var stats = struct {
	totalLeidos      int
	totalProcesados  int
	totalDescartados int
	textosVacios     int
	textosLimpiados  int
	mu               sync.Mutex
}{}

// =====================================
// HELPER METRICAS
// =====================================

func increment(field *int) {

	stats.mu.Lock()
	*field++
	stats.mu.Unlock()
}

// =====================================
// HELPERS NLP SIMPLE
// =====================================

func demasiadosSimbolos(text string) bool {

	count := 0

	for _, c := range text {

		if strings.ContainsRune("!@#$%^&*........", c) {
			count++
		}
	}

	return count > 8
}

func demasiadasRepeticiones(text string) bool {

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

// =====================================
// READER
// =====================================

func reader(path string, out chan<- Record) {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	r := csv.NewReader(file)
	r.FieldsPerRecord = -1

	// Saltar header
	_, err = r.Read()
	if err != nil {
		panic(err)
	}

	for {

		row, err := r.Read()

		if err != nil {
			break
		}

		increment(&stats.totalLeidos)

		// Validar columnas mínimas
		if len(row) <= COL_ID {
			continue
		}

		rec := Record{
			IDReclamo: row[COL_ID],
			UsuarioID: row[COL_USUARIO],
			IP:        row[COL_IP],
			Timestamp: row[COL_TIMESTAMP],

			TextoOriginal: row[COL_TEXTO],
		}

		out <- rec
	}

	close(out)
}

// =====================================
// NORMALIZER
// =====================================

func normalizer(
	in <-chan Record,
	out chan<- Record,
) {

	for rec := range in {

		original := rec.TextoOriginal

		// =============================
		// NORMALIZACION
		// =============================

		texto := strings.TrimSpace(original)

		// eliminar espacios dobles
		texto = strings.Join(
			strings.Fields(texto),
			" ",
		)

		rec.TextoLimpio = texto

		// =============================
		// METRICAS
		// =============================

		if texto != original {
			increment(&stats.textosLimpiados)
		}

		// =============================
		// TEXTO VACIO
		// =============================

		if texto == "" || texto == "-" {

			rec.TextoLimpio = "RECLAMO VACIO"

			increment(&stats.textosVacios)
		}

		out <- rec
	}
}

// =====================================
// VALIDATOR
// =====================================

func validator(
	in <-chan Record,
	valid chan<- Record,
	reject chan<- RejectedRecord,
) {

	for rec := range in {

		// =================================
		// TEXTO VACIO
		// =================================

		if rec.TextoLimpio == "RECLAMO VACIO" {

			increment(&stats.totalDescartados)

			reject <- RejectedRecord{
				Record: rec,
				Motivo: "TEXTO_VACIO",
			}

			continue
		}

		// =================================
		// TEXTO MUY CORTO
		// =================================

		if len(strings.Fields(rec.TextoLimpio)) < 3 {

			increment(&stats.totalDescartados)

			reject <- RejectedRecord{
				Record: rec,
				Motivo: "TEXTO_MUY_CORTO",
			}

			continue
		}

		// =================================
		// EXCESO DE SIMBOLOS
		// =================================

		if demasiadosSimbolos(rec.TextoLimpio) {

			increment(&stats.totalDescartados)

			reject <- RejectedRecord{
				Record: rec,
				Motivo: "EXCESO_SIMBOLOS",
			}

			continue
		}

		// =================================
		// REPETICION EXCESIVA
		// =================================

		if demasiadasRepeticiones(rec.TextoLimpio) {

			increment(&stats.totalDescartados)

			reject <- RejectedRecord{
				Record: rec,
				Motivo: "REPETICION_EXCESIVA",
			}

			continue
		}

		// =================================
		// VALIDAR IP
		// =================================

		if net.ParseIP(rec.IP) == nil {

			increment(&stats.totalDescartados)

			reject <- RejectedRecord{
				Record: rec,
				Motivo: "IP_INVALIDA",
			}

			continue
		}

		// =================================
		// TIMESTAMP VACIO
		// =================================

		if strings.TrimSpace(rec.Timestamp) == "" {
			rec.Timestamp = "SIN_FECHA"
		}

		// =================================
		// REGISTRO VALIDO
		// =================================

		increment(&stats.totalProcesados)

		valid <- rec
	}
}

// =====================================
// WRITER CLEAN
// =====================================

func writer(
	path string,
	in <-chan Record,
	done chan<- bool,
) {

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// =================================
	// HEADER
	// =================================

	w.Write([]string{
		"ID_RECLAMO",
		"USUARIO_ID",
		"IP",
		"TIMESTAMP",
		"TEXTO",
	})

	// =================================
	// ESCRITURA
	// =================================

	for rec := range in {

		w.Write([]string{
			rec.IDReclamo,
			rec.UsuarioID,
			rec.IP,
			rec.Timestamp,
			rec.TextoLimpio,
		})
	}

	done <- true
}

// =====================================
// WRITER REJECTED
// =====================================

func rejectedWriter(
	path string,
	in <-chan RejectedRecord,
	done chan<- bool,
) {

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// =================================
	// HEADER
	// =================================

	w.Write([]string{
		"ID_RECLAMO",
		"USUARIO_ID",
		"IP",
		"TIMESTAMP",
		"TEXTO",
		"MOTIVO",
	})

	// =================================
	// ESCRITURA
	// =================================

	for rec := range in {

		w.Write([]string{
			rec.IDReclamo,
			rec.UsuarioID,
			rec.IP,
			rec.Timestamp,
			rec.TextoLimpio,
			rec.Motivo,
		})
	}

	done <- true
}

// =====================================
// REPORTE TXT
// =====================================

func generarReporteTXT() {

	os.MkdirAll("../logs", os.ModePerm)

	file, err := os.Create(
		"../logs/cleaning_summary.txt",
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.WriteString(
		"===== REPORTE DE LIMPIEZA =====\n\n",
	)

	file.WriteString(
		fmt.Sprintf(
			"Total leidos: %d\n",
			stats.totalLeidos,
		),
	)

	file.WriteString(
		fmt.Sprintf(
			"Procesados: %d\n",
			stats.totalProcesados,
		),
	)

	file.WriteString(
		fmt.Sprintf(
			"Descartados: %d\n",
			stats.totalDescartados,
		),
	)

	file.WriteString(
		fmt.Sprintf(
			"Textos limpiados: %d\n",
			stats.textosLimpiados,
		),
	)

	file.WriteString(
		fmt.Sprintf(
			"Textos vacios: %d\n",
			stats.textosVacios,
		),
	)

	if stats.totalLeidos > 0 {

		porcentaje :=
			float64(stats.totalDescartados) /
				float64(stats.totalLeidos) * 100

		file.WriteString(
			fmt.Sprintf(
				"Porcentaje descartado: %.2f%%\n",
				porcentaje,
			),
		)
	}
}

// =====================================
// MAIN
// =====================================

func main() {

	start := time.Now()

	// =================================
	// CHANNELS
	// =================================

	rawChan := make(chan Record, 1000)

	cleanChan := make(chan Record, 1000)

	validChan := make(chan Record, 1000)

	rejectChan := make(chan RejectedRecord, 1000)

	// =================================
	// SIGNALS
	// =================================

	doneClean := make(chan bool)

	doneReject := make(chan bool)

	// =================================
	// READER
	// =================================

	go reader(
		"../dataset/dataset_1M_raw.csv",
		rawChan,
	)

	// =================================
	// NORMALIZERS
	// =================================

	var wgNorm sync.WaitGroup

	for i := 0; i < 4; i++ {

		wgNorm.Add(1)

		go func() {

			defer wgNorm.Done()

			normalizer(
				rawChan,
				cleanChan,
			)
		}()
	}

	go func() {

		wgNorm.Wait()

		close(cleanChan)
	}()

	// =================================
	// VALIDATORS
	// =================================

	var wgVal sync.WaitGroup

	for i := 0; i < 4; i++ {

		wgVal.Add(1)

		go func() {

			defer wgVal.Done()

			validator(
				cleanChan,
				validChan,
				rejectChan,
			)
		}()
	}

	go func() {

		wgVal.Wait()

		close(validChan)
		close(rejectChan)
	}()

	// =================================
	// WRITERS
	// =================================

	go writer(
		"../dataset/dataset_clean.csv",
		validChan,
		doneClean,
	)

	go rejectedWriter(
		"../dataset/rejected_records.csv",
		rejectChan,
		doneReject,
	)

	// =================================
	// ESPERAR WRITERS
	// =================================

	<-doneClean
	<-doneReject

	// =================================
	// REPORTE TXT
	// =================================

	generarReporteTXT()

	// =================================
	// REPORTE FINAL
	// =================================

	fmt.Println(
		"\n========== REPORTE ==========",
	)

	fmt.Println(
		"Total leidos:",
		stats.totalLeidos,
	)

	fmt.Println(
		"Procesados:",
		stats.totalProcesados,
	)

	fmt.Println(
		"Descartados:",
		stats.totalDescartados,
	)

	fmt.Println(
		"Textos limpiados:",
		stats.textosLimpiados,
	)

	fmt.Println(
		"Textos vacios:",
		stats.textosVacios,
	)

	if stats.totalLeidos > 0 {

		porcentaje :=
			float64(stats.totalDescartados) /
				float64(stats.totalLeidos) * 100

		fmt.Printf(
			"Porcentaje descartado: %.2f%%\n",
			porcentaje,
		)
	}

	fmt.Println(
		"Tiempo total:",
		time.Since(start),
	)

	fmt.Println("\nArchivos generados:")

	fmt.Println(
		"- dataset_clean.csv",
	)

	fmt.Println(
		"- rejected_records.csv",
	)

	fmt.Println(
		"- logs/cleaning_summary.txt",
	)
}
