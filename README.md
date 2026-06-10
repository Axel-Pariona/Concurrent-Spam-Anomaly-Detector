# Concurrent Spam Anomaly Detector

Concurrent Spam Anomaly Detector es un proyecto académico desarrollado en Go para procesar grandes volúmenes de registros de reclamos y detectar posibles patrones de spam o anomalías mediante reglas heurísticas, procesamiento concurrente y comparación de rendimiento frente a una versión secuencial.

El proyecto forma parte del curso **Programación Concurrente y Distribuida** y combina tres componentes principales:

- Generación y ampliación de dataset.
- Limpieza y procesamiento concurrente de datos.
- Detección de spam/anomalías con versiones secuencial y concurrente.
- Benchmark de rendimiento.
- Modelado formal con Promela/SPIN.

## Objetivo del proyecto

El objetivo principal es implementar y comparar un sistema de detección de spam y anomalías en registros de reclamos usando programación concurrente en Go.

El proyecto busca demostrar:

- Uso de goroutines.
- Comunicación mediante channels.
- Sincronización con `sync.Mutex` y `sync.WaitGroup`.
- Pipeline concurrente de limpieza de datos.
- Worker pool para clasificación concurrente.
- Comparación entre procesamiento secuencial y concurrente.
- Medición de tiempo, memoria y speedup.
- Modelado formal de concurrencia con Promela/SPIN.

## Contexto del problema

Los registros de reclamos pueden contener información incompleta, texto repetitivo, entradas inválidas, direcciones IP sospechosas o mensajes con características similares a spam.

Procesar estos datos de forma secuencial puede ser costoso cuando el volumen es alto. Por ello, el proyecto propone una solución concurrente capaz de dividir el trabajo en etapas y procesar registros en paralelo.

## Tecnologías utilizadas

- Go
- Python
- Pandas
- Promela
- SPIN
- CSV
- Goroutines
- Channels
- Worker Pool
- Pipeline concurrente

## Estructura del proyecto

```txt
Concurrent-Spam-Anomaly-Detector/
  README.md
  .gitignore
  go.mod

  data/
    README.md
    dataset_base.csv

  docs/
    ai_gap_analysis.md
    methodology.md

  formal/
    README.md
    modelo.pml
    modelo_deadlock.pml

  logs/
    cleaning_summary.txt
    generation_summary.txt

  scripts/
    generate_dataset.py

  src/
    benchmark.go
    cleaner_pipeline.go
    spam_concurrent.go
    spam_sequential.go
```

## Descripción de carpetas

### `data/`

Contiene el dataset base y documentación sobre los archivos generados.

- `dataset_base.csv`: dataset inicial usado como base para generar datos ampliados.
- `README.md`: explicación de los archivos de entrada, salida y archivos ignorados.

Los archivos grandes generados por el sistema no se versionan en GitHub.

### `scripts/`

Contiene scripts auxiliares.

- `generate_dataset.py`: genera un dataset ampliado a partir de `data/dataset_base.csv`.

### `src/`

Contiene el código principal en Go.

- `cleaner_pipeline.go`: pipeline concurrente para limpiar y validar registros.
- `spam_sequential.go`: versión secuencial del detector de spam.
- `spam_concurrent.go`: versión concurrente del detector de spam usando worker pool.
- `benchmark.go`: ejecuta ambas versiones varias veces y calcula promedios, memoria y speedup.

### `formal/`

Contiene el modelado formal en Promela.

- `modelo.pml`: modelo formal principal del pipeline concurrente.
- `modelo_deadlock.pml`: modelo utilizado para análisis de posibles bloqueos.
- `README.md`: instrucciones para ejecutar verificación con SPIN.

### `logs/`

Contiene resúmenes generados durante el proceso.

- `generation_summary.txt`: resumen de generación del dataset.
- `cleaning_summary.txt`: resumen de limpieza del dataset.

### `docs/`

Contiene documentación complementaria.

- `methodology.md`: explicación metodológica del proyecto.
- `ai_gap_analysis.md`: análisis técnico de brechas, mejoras y revisión del proyecto.

## Flujo general del sistema

```txt
dataset_base.csv
  ↓
generate_dataset.py
  ↓
dataset_1M_raw.csv
  ↓
cleaner_pipeline.go
  ↓
dataset_clean.csv
  ↓
spam_sequential.go / spam_concurrent.go
  ↓
dataset_final_*.csv
spam_detected_*.csv
  ↓
benchmark.go
  ↓
comparación de rendimiento
```

## Arquitectura concurrente

El proyecto utiliza dos enfoques principales de concurrencia.

### Pipeline concurrente

El archivo `cleaner_pipeline.go` implementa un pipeline con etapas especializadas:

```txt
Reader
  ↓
Normalizer
  ↓
Validator
  ↓
Deduplicator
  ↓
Output
```

Cada etapa procesa registros y se comunica con la siguiente mediante channels.

### Worker pool

El archivo `spam_concurrent.go` utiliza un conjunto de workers para clasificar registros en paralelo.

```txt
Input records
  ↓
Jobs channel
  ↓
Workers
  ↓
Results channel
  ↓
Output files
```

## Reglas de detección de spam

El detector usa reglas heurísticas basadas en señales como:

- Texto repetido muchas veces.
- IP con alta frecuencia de registros.
- Exceso de mayúsculas.
- Baja diversidad léxica.
- Palabras asociadas a spam.
- Texto demasiado corto.
- Repetición excesiva de palabras.
- Exceso de símbolos.
- Texto artificial.
- Ausencia de conectores legítimos.

Cada regla suma un puntaje. Si el puntaje supera un umbral, el registro se clasifica como posible spam.

## Requisitos

### Go

Se recomienda Go 1.22 o superior.

Verificar instalación:

```bash
go version
```

### Python

Se recomienda Python 3.10 o superior para generar el dataset.

Instalar dependencias:

```bash
pip install pandas
```

### SPIN

Para validar los modelos Promela, se requiere SPIN instalado.

Verificar instalación:

```bash
spin -V
```

## Ejecución del proyecto

### 1. Generar dataset ampliado

Desde la raíz del proyecto:

```bash
python scripts/generate_dataset.py
```

Esto genera:

```txt
data/dataset_1M_raw.csv
logs/generation_summary.txt
```

### 2. Ejecutar limpieza concurrente

Entrar a la carpeta `src`:

```bash
cd src
```

Ejecutar:

```bash
go run cleaner_pipeline.go
```

Esto genera:

```txt
data/dataset_clean.csv
data/rejected_records.csv
logs/cleaning_summary.txt
```

### 3. Ejecutar detector secuencial

Desde `src`:

```bash
go run spam_sequential.go
```

Esto genera:

```txt
data/dataset_final_secuencial.csv
data/spam_detected_secuencial.csv
```

### 4. Ejecutar detector concurrente

Desde `src`:

```bash
go run spam_concurrent.go
```

Esto genera:

```txt
data/dataset_final_concurrente.csv
data/spam_detected_concurrente.csv
```

### 5. Ejecutar benchmark

Desde `src`:

```bash
go run benchmark.go
```

El benchmark ejecuta varias veces la versión secuencial y concurrente, calcula una media recortada y reporta:

- Tiempo promedio.
- Memoria promedio.
- Speedup.
- Comparación de rendimiento.

## Verificación formal con SPIN

Los modelos Promela se encuentran en:

```txt
formal/
```

Para ejecutar una simulación:

```bash
cd formal
spin modelo.pml
```

Para generar el verificador:

```bash
spin -a modelo.pml
gcc -o pan pan.c
./pan
```

En Windows, si usas MinGW:

```powershell
spin -a modelo.pml
gcc -o pan.exe pan.c
.\pan.exe
```

Para más detalle:

[`formal/README.md`](formal/README.md)

## Archivos generados

Durante la ejecución se pueden generar archivos grandes:

```txt
data/dataset_1M_raw.csv
data/dataset_clean.csv
data/rejected_records.csv
data/dataset_final_secuencial.csv
data/dataset_final_concurrente.csv
data/spam_detected_secuencial.csv
data/spam_detected_concurrente.csv
```

Estos archivos están excluidos del repositorio mediante `.gitignore`.

## Reproducibilidad

Para reproducir el flujo completo:

```bash
python scripts/generate_dataset.py
cd src
go run cleaner_pipeline.go
go run spam_sequential.go
go run spam_concurrent.go
go run benchmark.go
```

Luego revisar:

```txt
logs/
data/
```

## Alcance del proyecto

El proyecto tiene fines académicos y experimentales.

Incluye:

- Generación de datos sintéticos ampliados.
- Limpieza concurrente de registros.
- Detección heurística de spam.
- Comparación secuencial vs concurrente.
- Benchmark de rendimiento.
- Modelado formal de procesos concurrentes.

## Limitaciones

- La detección de spam se basa en reglas heurísticas, no en un modelo de Machine Learning.
- Los datos ampliados son sintéticos y derivados de un dataset base.
- Los resultados pueden variar según hardware y carga del sistema.
- La concurrencia puede mejorar tiempos, pero también introduce overhead.
- El speedup depende del volumen de datos, número de workers y costo de cada tarea.
- La verificación formal modela una abstracción del sistema, no todo el código Go completo.

## Posibles mejoras

- Parametrizar rutas y tamaño del dataset mediante flags CLI.
- Agregar tests unitarios para reglas de spam.
- Agregar métricas de precisión si se cuenta con etiquetas reales.
- Exportar resultados de benchmark en CSV.
- Permitir configurar número de workers.
- Mejorar el manejo de errores.
- Agregar contexto de cancelación con `context.Context`.
- Separar paquetes Go en módulos reutilizables.
- Incorporar profiling con `pprof`.
- Ampliar propiedades LTL en Promela.

## Estado del proyecto

Proyecto académico funcional, reorganizado y documentado para presentación en GitHub.

## Autores

- Omar Junior Acuña Villegas
- Rafael Tomás Chui Sánchez
- Axel Yamir Pariona Rojas
