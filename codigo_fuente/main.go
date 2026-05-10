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

const (
	COL_TIMESTAMP = 10
	COL_TEXTO     = 11
	COL_USUARIO   = 12
	COL_IP        = 13
	COL_ID        = 14
)

type Record struct {
	IDReclamo string
	UsuarioID string
	IP        string
	Timestamp string
	Texto     string
}

type RejectedRecord struct {
	Record
	Motivo string
}

// =========================
// MÉTRICAS GLOBALES
// =========================

var stats = struct {
	totalLeidos      int
	totalProcesados  int
	totalDescartados int
	textosVacios     int
	textosLimpiados  int
	mu               sync.Mutex
}{}

// =========================
// HELPERS
// =========================

func increment(field *int) {
	stats.mu.Lock()
	*field++
	stats.mu.Unlock()
}

// =========================
// READER
// =========================

func reader(path string, out chan<- Record) {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.FieldsPerRecord = -1

	// Leer cabecera
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
			Texto:     row[COL_TEXTO],
		}

		out <- rec
	}

	close(out)
}

// =========================
// NORMALIZER
// =========================

func normalizer(in <-chan Record, out chan<- Record) {

	for rec := range in {

		original := rec.Texto

		// Limpiar texto
		rec.Texto = strings.ToUpper(rec.Texto)
		rec.Texto = strings.TrimSpace(rec.Texto)
		rec.Texto = strings.Join(strings.Fields(rec.Texto), " ")

		// Contar cambios
		if rec.Texto != original {
			increment(&stats.textosLimpiados)
		}

		// Texto vacío
		if rec.Texto == "" || rec.Texto == "-" {
			rec.Texto = "RECLAMO VACIO"
			increment(&stats.textosVacios)
		}

		out <- rec
	}
}

// =========================
// VALIDATOR
// =========================

func validator(
	in <-chan Record,
	valid chan<- Record,
	reject chan<- RejectedRecord,
) {

	for rec := range in {

		// Texto vacío
		if rec.Texto == "RECLAMO VACIO" {

			increment(&stats.totalDescartados)

			reject <- RejectedRecord{
				Record: rec,
				Motivo: "TEXTO_VACIO",
			}

			continue
		}

		// Validar IP
		if net.ParseIP(rec.IP) == nil {

			increment(&stats.totalDescartados)

			reject <- RejectedRecord{
				Record: rec,
				Motivo: "IP_INVALIDA",
			}

			continue
		}

		// Timestamp vacío
		if strings.TrimSpace(rec.Timestamp) == "" {
			rec.Timestamp = "SIN_FECHA"
		}

		increment(&stats.totalProcesados)

		valid <- rec
	}
}

// =========================
// WRITER CLEAN
// =========================

func writer(path string, in <-chan Record, done chan<- bool) {

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// Header
	w.Write([]string{
		"ID_RECLAMO",
		"USUARIO_ID",
		"IP",
		"TIMESTAMP",
		"TEXTO",
	})

	for rec := range in {

		w.Write([]string{
			rec.IDReclamo,
			rec.UsuarioID,
			rec.IP,
			rec.Timestamp,
			rec.Texto,
		})
	}

	done <- true
}

// =========================
// WRITER REJECTED
// =========================

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

	w.Write([]string{
		"ID_RECLAMO",
		"USUARIO_ID",
		"IP",
		"TIMESTAMP",
		"TEXTO",
		"MOTIVO",
	})

	for rec := range in {

		w.Write([]string{
			rec.IDReclamo,
			rec.UsuarioID,
			rec.IP,
			rec.Timestamp,
			rec.Texto,
			rec.Motivo,
		})
	}

	done <- true
}

// =========================
// GENERAR REPORTE TXT
// =========================

func generarReporteTXT() {

	// Crear carpeta logs
	os.MkdirAll("../logs", os.ModePerm)

	file, err := os.Create("../logs/cleaning_summary.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("===== REPORTE DE LIMPIEZA =====\n\n")

	file.WriteString(
		fmt.Sprintf("Total leídos: %d\n", stats.totalLeidos),
	)

	file.WriteString(
		fmt.Sprintf("Procesados: %d\n", stats.totalProcesados),
	)

	file.WriteString(
		fmt.Sprintf("Descartados: %d\n", stats.totalDescartados),
	)

	file.WriteString(
		fmt.Sprintf("Textos limpiados: %d\n", stats.textosLimpiados),
	)

	file.WriteString(
		fmt.Sprintf("Textos vacíos: %d\n", stats.textosVacios),
	)

	if stats.totalLeidos > 0 {

		porcentaje := float64(stats.totalDescartados) /
			float64(stats.totalLeidos) * 100

		file.WriteString(
			fmt.Sprintf("Porcentaje descartado: %.2f%%\n", porcentaje),
		)
	}
}

// =========================
// MAIN
// =========================

func main() {

	start := time.Now()

	// Channels
	rawChan := make(chan Record, 1000)
	cleanChan := make(chan Record, 1000)
	validChan := make(chan Record, 1000)
	rejectChan := make(chan RejectedRecord, 1000)

	// Señales
	doneClean := make(chan bool)
	doneReject := make(chan bool)

	// =========================
	// READER
	// =========================

	go reader("../dataset/dataset_1M_raw.csv", rawChan)

	// =========================
	// NORMALIZERS
	// =========================

	var wgNorm sync.WaitGroup

	for i := 0; i < 4; i++ {

		wgNorm.Add(1)

		go func() {
			defer wgNorm.Done()
			normalizer(rawChan, cleanChan)
		}()
	}

	go func() {
		wgNorm.Wait()
		close(cleanChan)
	}()

	// =========================
	// VALIDATORS
	// =========================

	var wgVal sync.WaitGroup

	for i := 0; i < 4; i++ {

		wgVal.Add(1)

		go func() {
			defer wgVal.Done()
			validator(cleanChan, validChan, rejectChan)
		}()
	}

	go func() {
		wgVal.Wait()
		close(validChan)
		close(rejectChan)
	}()

	// =========================
	// WRITERS
	// =========================

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

	// Esperar writers
	<-doneClean
	<-doneReject
	generarReporteTXT()
	// =========================
	// REPORTE FINAL
	// =========================

	fmt.Println("========== REPORTE ==========")
	fmt.Println("Total leídos: ", stats.totalLeidos)
	fmt.Println("Procesados: ", stats.totalProcesados)
	fmt.Println("Descartados: ", stats.totalDescartados)
	fmt.Println("Textos limpiados: ", stats.textosLimpiados)
	fmt.Println("Textos vacíos: ", stats.textosVacios)

	if stats.totalLeidos > 0 {

		porcentaje := float64(stats.totalDescartados) /
			float64(stats.totalLeidos) * 100

		fmt.Printf("Porcentaje descartado: %.2f%%\n", porcentaje)
	}

	fmt.Println("Tiempo total:", time.Since(start))

	fmt.Println("\nArchivos generados:")
	fmt.Println("- dataset_clean.csv")
	fmt.Println("- rejected_records.csv")
	fmt.Println("- logs/cleaning_summary.txt")
}
