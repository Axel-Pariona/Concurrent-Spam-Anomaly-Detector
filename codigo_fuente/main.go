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

type Reclamo struct {
	ID              string
	Timestamp       string
	UsuarioID       string
	IPAddress       string
	Empresa         string
	Departamento    string
	Servicio        string
	MedioPresent    string
	TipoQueja       string
	TextoReclamo    string
	IsSyntheticSpam string
}

type CleaningEvent struct {
	RecordID      string
	Action        string
	Field         string
	OriginalValue string
	NewValue      string
	Reason        string
}

var (
	totalRead          int
	totalClean         int
	totalNormalized    int
	totalInvalidIP     int
	totalInvalidTS     int
	totalMissingFields int
	totalDuplicates    int
	statsMutex         sync.Mutex
)

func fixEncoding(s string) string {
	replacements := map[string]string{
		"Ã‰": "É",
		"Ã“": "Ó",
		"Ãš": "Ú",
		"Ã": "Á",
		"Ã": "Í",
		"Ã‘": "Ñ",
		"Ã¡": "á",
		"Ã©": "é",
		"Ã­": "í",
		"Ã³": "ó",
		"Ãº": "ú",
		"Ã±": "ñ",
	}

	for wrong, correct := range replacements {
		s = strings.ReplaceAll(s, wrong, correct)
	}

	return s
}

func normalizeField(s string) string {
	s = fixEncoding(s)
	s = strings.TrimSpace(s)

	if s == "-" || s == "null" || s == "N/A" {
		return ""
	}

	return strings.ToUpper(s)
}

func normalizeTimestamp(ts string) (string, error) {
	layouts := []string{
		"2/01/2006 15:04",
		"2006-01-02 15:04:05",
		"02-01-2006 15:04",
		"2006/01/02 15:04",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, ts); err == nil {
			return parsed.Format("2006-01-02 15:04:05"), nil
		}
	}

	return "", fmt.Errorf("invalid timestamp")
}

func normalizerWorker(
	in <-chan Reclamo,
	out chan<- Reclamo,
	audit chan<- CleaningEvent,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for r := range in {

		fields := map[string]*string{
			"empresa":            &r.Empresa,
			"departamento":       &r.Departamento,
			"servicio":           &r.Servicio,
			"medio_presentacion": &r.MedioPresent,
			"tipo_queja":         &r.TipoQueja,
			"texto_reclamo":      &r.TextoReclamo,
		}

		for fieldName, fieldPtr := range fields {
			original := *fieldPtr
			normalized := normalizeField(original)

			if original != normalized {
				*fieldPtr = normalized

				audit <- CleaningEvent{
					RecordID:      r.ID,
					Action:        "NORMALIZED",
					Field:         fieldName,
					OriginalValue: original,
					NewValue:      normalized,
					Reason:        "trim/uppercase/encoding_fix",
				}

				statsMutex.Lock()
				totalNormalized++
				statsMutex.Unlock()
			}
		}

		originalTS := r.Timestamp
		ts, err := normalizeTimestamp(r.Timestamp)
		if err == nil && ts != originalTS {
			r.Timestamp = ts

			audit <- CleaningEvent{
				RecordID:      r.ID,
				Action:        "NORMALIZED",
				Field:         "timestamp",
				OriginalValue: originalTS,
				NewValue:      ts,
				Reason:        "timestamp_standardization",
			}

			statsMutex.Lock()
			totalNormalized++
			statsMutex.Unlock()
		}

		out <- r
	}
}

func validatorWorker(
	in <-chan Reclamo,
	out chan<- Reclamo,
	audit chan<- CleaningEvent,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for r := range in {

		if r.UsuarioID == "" || r.TextoReclamo == "" {
			audit <- CleaningEvent{
				RecordID:      r.ID,
				Action:        "REMOVED",
				Field:         "required_fields",
				OriginalValue: "",
				NewValue:      "",
				Reason:        "missing_required_field",
			}

			statsMutex.Lock()
			totalMissingFields++
			statsMutex.Unlock()

			continue
		}

		if _, err := time.Parse("2006-01-02 15:04:05", r.Timestamp); err != nil {
			audit <- CleaningEvent{
				RecordID:      r.ID,
				Action:        "REMOVED",
				Field:         "timestamp",
				OriginalValue: r.Timestamp,
				NewValue:      "",
				Reason:        "invalid_timestamp",
			}

			statsMutex.Lock()
			totalInvalidTS++
			statsMutex.Unlock()

			continue
		}

		if net.ParseIP(r.IPAddress) == nil {
			audit <- CleaningEvent{
				RecordID:      r.ID,
				Action:        "REMOVED",
				Field:         "ip_address",
				OriginalValue: r.IPAddress,
				NewValue:      "",
				Reason:        "invalid_ip",
			}

			statsMutex.Lock()
			totalInvalidIP++
			statsMutex.Unlock()

			continue
		}

		out <- r
	}
}

func deduplicator(
	in <-chan Reclamo,
	out chan<- Reclamo,
	audit chan<- CleaningEvent,
) {
	seen := make(map[string]bool)

	for r := range in {
		key := r.UsuarioID + "|" + r.Timestamp + "|" + r.TextoReclamo

		if seen[key] {
			audit <- CleaningEvent{
				RecordID:      r.ID,
				Action:        "REMOVED",
				Field:         "duplicate",
				OriginalValue: "",
				NewValue:      "",
				Reason:        "duplicate_record",
			}

			statsMutex.Lock()
			totalDuplicates++
			statsMutex.Unlock()

			continue
		}

		seen[key] = true
		out <- r
	}

	close(out)
}

func main() {

	file, err := os.Open("../Dataset/dataset_1M_raw.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	outputFile, _ := os.Create("../Dataset/dataset_clean.csv")
	defer outputFile.Close()

	reportFile, _ := os.Create("../Dataset/cleaning_report.csv")
	defer reportFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	reportWriter := csv.NewWriter(reportFile)
	defer reportWriter.Flush()

	header, _ := reader.Read()
	writer.Write(header)

	reportWriter.Write([]string{
		"record_id",
		"action",
		"field",
		"original_value",
		"new_value",
		"reason",
	})

	rawChan := make(chan Reclamo, 1000)
	normChan := make(chan Reclamo, 1000)
	validChan := make(chan Reclamo, 1000)
	finalChan := make(chan Reclamo, 1000)
	auditChan := make(chan CleaningEvent, 10000)

	var auditWG sync.WaitGroup
	auditWG.Add(1)

	go func() {
		defer auditWG.Done()

		for event := range auditChan {
			reportWriter.Write([]string{
				event.RecordID,
				event.Action,
				event.Field,
				event.OriginalValue,
				event.NewValue,
				event.Reason,
			})
		}
	}()

	var normWG sync.WaitGroup
	for i := 0; i < 4; i++ {
		normWG.Add(1)
		go normalizerWorker(rawChan, normChan, auditChan, &normWG)
	}

	var valWG sync.WaitGroup
	for i := 0; i < 4; i++ {
		valWG.Add(1)
		go validatorWorker(normChan, validChan, auditChan, &valWG)
	}

	go deduplicator(validChan, finalChan, auditChan)

	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				break
			}

			statsMutex.Lock()
			totalRead++
			statsMutex.Unlock()

			rawChan <- Reclamo{
				ID:              record[0],
				Timestamp:       record[1],
				UsuarioID:       record[2],
				IPAddress:       record[3],
				Empresa:         record[4],
				Departamento:    record[5],
				Servicio:        record[6],
				MedioPresent:    record[7],
				TipoQueja:       record[8],
				TextoReclamo:    record[9],
				IsSyntheticSpam: record[10],
			}
		}

		close(rawChan)
	}()

	go func() {
		normWG.Wait()
		close(normChan)
	}()

	go func() {
		valWG.Wait()
		close(validChan)
	}()

	cleanID := 1

	for r := range finalChan {
		writer.Write([]string{
			fmt.Sprintf("%d", cleanID),
			r.Timestamp,
			r.UsuarioID,
			r.IPAddress,
			r.Empresa,
			r.Departamento,
			r.Servicio,
			r.MedioPresent,
			r.TipoQueja,
			r.TextoReclamo,
			r.IsSyntheticSpam,
		})

		cleanID++
		totalClean++
	}

	close(auditChan)
	auditWG.Wait()

	fmt.Println("\n===== Resumen de la limpieza =====")
	fmt.Println("Total Read:", totalRead)
	fmt.Println("Final Clean:", totalClean)
	fmt.Println("Normalized Fields:", totalNormalized)
	fmt.Println("Invalid IP Removed:", totalInvalidIP)
	fmt.Println("Invalid Timestamp Removed:", totalInvalidTS)
	fmt.Println("Missing Fields Removed:", totalMissingFields)
	fmt.Println("Duplicates Removed:", totalDuplicates)
}
