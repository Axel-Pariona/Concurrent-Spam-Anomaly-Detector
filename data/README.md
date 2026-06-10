# Data

Esta carpeta contiene el dataset base y los archivos de datos generados durante la ejecución del proyecto **Concurrent Spam Anomaly Detector**.

## Archivo versionado

El archivo principal que sí se conserva en GitHub es:

```txt
dataset_base.csv
```

Este archivo sirve como entrada inicial para generar un dataset ampliado de mayor tamaño.

## Archivos generados

Durante la ejecución del proyecto se generan varios archivos derivados:

```txt
dataset_1M_raw.csv
dataset_clean.csv
rejected_records.csv
dataset_final_secuencial.csv
dataset_final_concurrente.csv
spam_detected_secuencial.csv
spam_detected_concurrente.csv
```

Estos archivos no se suben al repositorio porque pueden ser grandes y pueden regenerarse ejecutando los scripts del proyecto.

## Flujo de datos

```txt
dataset_base.csv
  ↓
scripts/generate_dataset.py
  ↓
dataset_1M_raw.csv
  ↓
src/cleaner_pipeline.go
  ↓
dataset_clean.csv
  ↓
src/spam_sequential.go / src/spam_concurrent.go
  ↓
dataset_final_*.csv
spam_detected_*.csv
```

## Descripción de archivos

### `dataset_base.csv`

Dataset inicial usado como fuente para generar registros adicionales.

### `dataset_1M_raw.csv`

Dataset ampliado generado por el script de Python.

Contiene registros válidos, registros con ruido, duplicados, entradas incompletas, direcciones IP inválidas y posibles patrones de spam.

Se genera con:

```bash
python scripts/generate_dataset.py
```

### `dataset_clean.csv`

Dataset limpio generado por el pipeline concurrente en Go.

Se genera con:

```bash
cd src
go run cleaner_pipeline.go
```

### `rejected_records.csv`

Archivo con registros descartados durante la limpieza.

Puede incluir registros con:

- Texto vacío.
- Texto muy corto.
- IP inválida.
- Timestamp inválido.
- Exceso de símbolos.
- Repeticiones artificiales.

### `dataset_final_secuencial.csv`

Resultado final del detector secuencial.

Se genera con:

```bash
cd src
go run spam_sequential.go
```

### `spam_detected_secuencial.csv`

Registros clasificados como posible spam por la versión secuencial.

### `dataset_final_concurrente.csv`

Resultado final del detector concurrente.

Se genera con:

```bash
cd src
go run spam_concurrent.go
```

### `spam_detected_concurrente.csv`

Registros clasificados como posible spam por la versión concurrente.

## Columnas usadas por el sistema

El sistema trabaja principalmente con columnas del dataset relacionadas con:

- Timestamp.
- Texto del reclamo.
- Usuario.
- IP.
- Identificador del reclamo.

En el pipeline de limpieza, estas columnas son extraídas por posición desde el CSV original.

## Por qué no se suben los archivos generados

Los archivos generados pueden tener cientos de miles o millones de registros. Por ello:

- Aumentan demasiado el peso del repositorio.
- Pueden regenerarse desde el dataset base.
- No son necesarios para revisar el código.
- Dificultan el mantenimiento del proyecto.

El `.gitignore` conserva `dataset_base.csv` y excluye los archivos generados.

## Reproducción de datos

Para reconstruir los archivos generados desde cero:

```bash
python scripts/generate_dataset.py
cd src
go run cleaner_pipeline.go
go run spam_sequential.go
go run spam_concurrent.go
```

## Logs asociados

Los resúmenes de generación y limpieza se guardan en:

```txt
logs/generation_summary.txt
logs/cleaning_summary.txt
```

Estos archivos permiten revisar cuántos registros fueron generados, procesados, limpiados o rechazados.

## Nota

Esta carpeta forma parte de un proyecto académico de programación concurrente y distribuida. Los datos generados son utilizados para evaluar rendimiento, concurrencia y escalabilidad del sistema.
