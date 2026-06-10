# Informe Técnico — Análisis de GAPs con IA

**Proyecto:** Sistema concurrente para la detección de spam y anomalías en registros de reclamos  
**Curso:** Programación Concurrente y Distribuida  
**Fecha:** 2026-05-10  
**Modelo de IA utilizado:** GPT-5.5 Thinking  

---

## 1. Objetivo del análisis

El objetivo de este informe es analizar técnicamente el código fuente del proyecto, identificando fortalezas y GAPs relacionados con calidad de código, seguridad, manejo de errores, patrones de concurrencia, rendimiento, uso de memoria, escalabilidad, mantenibilidad y verificación formal.

El análisis se realizó sobre un sistema implementado principalmente en Go, complementado con un script Python para generación de datos y un modelo Promela para verificación formal mediante SPIN.

---

## 2. Prompt utilizado

Actúa como revisor técnico de un proyecto del curso Programación Concurrente y Distribuida.

Analiza el código fuente de un sistema concurrente para detección de spam y anomalías en registros de reclamos. El proyecto está implementado principalmente en Go y contiene:

1. Un script Python para generar y ampliar un dataset experimental.
2. Un pipeline concurrente de limpieza en Go usando goroutines, channels, sync.WaitGroup y sync.Mutex.
3. Una versión secuencial del detector de spam.
4. Una versión concurrente del detector de spam usando worker pool.
5. Un benchmark para medir tiempo de ejecución, memoria y speedup.
6. Un modelo Promela para verificar sincronización, exclusión mutua y ausencia de deadlocks mediante SPIN.

Evalúa el código considerando calidad de código, seguridad, manejo de errores, patrones de concurrencia, posibles condiciones de carrera, uso de memoria, rendimiento, escalabilidad, mantenibilidad, coherencia entre implementación e informe y oportunidades de mejora.

Devuelve el resultado en formato Markdown con resumen general, fortalezas técnicas, GAPs encontrados, impacto de cada GAP, recomendaciones y conclusión técnica.

---

## 3. Archivos analizados

| Archivo | Descripción |
|---|---|
| `ampliacion_de_datos.py` | Genera y amplía el dataset experimental hasta 1,700,000 registros. |
| `main.go` | Implementa el pipeline concurrente de limpieza de datos. |
| `secuencial.go` | Implementa la versión secuencial del detector de spam. |
| `concurrente.go` | Implementa la versión concurrente del detector de spam mediante worker pool. |
| `benchmark.go` | Ejecuta pruebas de rendimiento, memoria y cálculo de speedup. |
| `modelo.pml` | Modela en Promela el pipeline concurrente para verificación con SPIN. |

---

## 4. Resumen general del análisis

El proyecto presenta una solución funcional y coherente con el objetivo del trabajo. Se implementó un flujo completo compuesto por generación de datos, limpieza concurrente, clasificación secuencial, clasificación concurrente, benchmarking y modelado formal en Promela.

La solución utiliza correctamente conceptos del curso, tales como `goroutines`, `channels`, `sync.WaitGroup`, `sync.Mutex`, pipeline, worker pool, speedup, media recortada y verificación formal con SPIN.

El sistema logra procesar un dataset ampliado de 1,700,000 registros, de los cuales 1,489,190 registros fueron considerados válidos después de la limpieza. La versión concurrente del detector obtuvo un speedup máximo de 2.15x con 8 workers, lo que demuestra una mejora significativa frente a la versión secuencial.

Sin embargo, el análisis también permitió identificar algunos GAPs técnicos relacionados con manejo de errores, uso de memoria, duplicación de código, parametrización del número de workers, medición del benchmark, evaluación de precisión y alcance del modelo Promela.

---

## 5. Fortalezas técnicas identificadas

### 5.1 Flujo completo del sistema

El proyecto no se limita únicamente a clasificar registros, sino que implementa un flujo completo:

- Generación y ampliación del dataset.
- Limpieza concurrente de datos.
- Clasificación secuencial.
- Clasificación concurrente.
- Medición de rendimiento.
- Modelado formal en Promela.
- Verificación con SPIN.

Esto demuestra una solución integral y alineada con los entregables solicitados.

---

### 5.2 Uso adecuado de programación concurrente en Go

El código utiliza correctamente elementos propios de Go para programación concurrente:

- `goroutines`
- `channels`
- `sync.WaitGroup`
- `sync.Mutex`

En el archivo `main.go`, el proceso de limpieza está organizado como un pipeline concurrente:

```go
rawChan := make(chan Record, 1000)
cleanChan := make(chan Record, 1000)
validChan := make(chan Record, 1000)
rejectChan := make(chan RejectedRecord, 1000)

Este diseño permite que la lectura, normalización, validación y escritura se ejecuten como etapas conectadas mediante canales.

## 5.3 Implementación de pipeline concurrente

El proceso de limpieza sigue una arquitectura clara:

Reader → Normalizer Workers → Validator Workers → Writer / Rejected Writer

Esta estructura es adecuada para el procesamiento de datos en flujo, ya que cada etapa cumple una responsabilidad específica.

Además, se utilizan `WaitGroup` para cerrar correctamente los canales una vez que los workers han terminado:

```go
go func() {
	wgNorm.Wait()
	close(cleanChan)
}()
```

y:

```go
go func() {
	wgVal.Wait()
	close(validChan)
	close(rejectChan)
}()
```

Esto reduce el riesgo de cerrar canales antes de tiempo.

## 5.4 Uso de worker pool en la versión concurrente

En `concurrente.go`, la clasificación utiliza un patrón Worker Pool. Los registros se envían a `inputChan`, son procesados por varios workers y luego los resultados se envían a `outputChan`.

```go
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
```

Este patrón es adecuado porque la clasificación de cada registro puede realizarse de forma independiente una vez calculadas las frecuencias globales.

## 5.5 Comparación justa entre versión secuencial y concurrente

Las versiones `secuencial.go` y `concurrente.go` aplican las mismas reglas heurísticas de clasificación. Esto permite que la comparación de rendimiento sea válida, ya que ambas versiones trabajan sobre el mismo dataset limpio y con la misma lógica de detección.

El criterio de clasificación utilizado es:

```go
isSpam := score >= 6
```

Esto permite una comparación directa entre ambos enfoques.

## 5.6 Uso de media recortada en el benchmark

El archivo `benchmark.go` ejecuta varias pruebas y calcula una media recortada eliminando los valores extremos:

```go
for i := 1; i < len(values)-1; i++ {
	sum += values[i]
}
```

Esto es positivo porque reduce el impacto de ejecuciones atípicas y permite obtener resultados más estables.

## 5.7 Verificación formal con Promela y SPIN

El archivo `modelo.pml` representa una abstracción del pipeline concurrente mediante procesos Reader, Normalizer y Validator.

El modelo usa canales FIFO:

```promela
chan ch_raw   = [20] of { mtype };
chan ch_clean = [20] of { mtype };
```

También utiliza mensajes de control:

```promela
mtype = { DATA, END };
```

La propiedad LTL definida permite verificar que los contadores principales no tomen valores negativos:

```promela
ltl NO_NEGATIVOS {
	[] (
		read_count >= 0 &&
		rejected_count >= 0 &&
		validated_count >= 0
	)
}
```

La verificación con SPIN reportó `errors: 0`, lo cual evidencia que no se encontraron violaciones en el modelo evaluado.

Además, se realizó una segunda ejecución sin la propiedad LTL para verificar estados finales inválidos, obteniendo también `errors: 0`, lo que permite sustentar la ausencia de deadlocks dentro del modelo abstracto.

## 6. GAPs técnicos identificados

### GAP 1: Uso de `go run` dentro del benchmark

**Descripción**

En `benchmark.go`, las pruebas se ejecutan usando:

```go
duration, output := runProgram(
	"go",
	"run",
	"secuencial.go",
)
```

y:

```go
duration, output := runProgram(
	"go",
	"run",
	"concurrente.go",
)
```

El uso de `go run` puede incluir tiempo de compilación dentro de la medición total. Esto puede afectar la precisión del benchmark, ya que el objetivo principal es medir el tiempo de ejecución del programa, no el tiempo de compilación.

**Impacto**

El speedup calculado puede verse influenciado por factores externos al procesamiento del dataset. Aunque se aplica el mismo método a ambas versiones, la medición sería más precisa si se ejecutaran binarios previamente compilados.

**Recomendación**

Compilar previamente ambos programas:

```bash
go build -o secuencial secuencial.go
go build -o concurrente concurrente.go
```

Luego ejecutar en el benchmark:

```go
runProgram("./secuencial")
runProgram("./concurrente")
```

De esta forma, el tiempo medido correspondería principalmente a la ejecución real del detector.

### GAP 2: Número de workers configurado manualmente

**Descripción**

En `concurrente.go`, el número de workers está definido manualmente:

```go
numWorkers := 4 //runtime.NumCPU()
```

Para probar con 4, 8 y 16 workers, se debe modificar el código manualmente antes de cada ejecución.

**Impacto**

Reduce la reproducibilidad de las pruebas y dificulta automatizar experimentos.

**Recomendación**

Permitir configurar el número de workers mediante argumentos de consola:

```go
numWorkers := runtime.NumCPU()
if len(os.Args) > 1 {
	n, err := strconv.Atoi(os.Args[1])
	if err == nil && n > 0 {
		numWorkers = n
	}
}
```

Y ejecutar:

```bash
go run concurrente.go 4
go run concurrente.go 8
go run concurrente.go 16
```

### GAP 3: Carga completa del dataset en memoria

**Descripción**

Tanto la versión secuencial como la concurrente cargan todos los registros limpios en memoria:

```go
var records []Record
records = append(records, Record{ /* ... */ })
```

**Impacto**

Funciona para ~1.5M registros, pero limita la escalabilidad si el dataset crece a decenas de millones.

**Recomendación**

Procesamiento por bloques/chunks o dos pasadas: primera para frecuencias, segunda para clasificación sin cargar todo en memoria.

### GAP 4: Conteo de frecuencias secuencial

**Descripción**

La fase de lectura y cálculo de frecuencias se mantiene secuencial:

```go
textFreq[text]++
ipFreq[ip]++
```

**Impacto**

Limita el speedup máximo al mantener una sección secuencial significativa.

**Recomendación**

Paralelizar el conteo con mapas locales por worker y luego reducir/mergearlos.

### GAP 5: Duplicación de lógica entre `secuencial.go` y `concurrente.go`

**Descripción**

Funciones compartidas (p. ej. `containsSpamWords`, `capsRatio`, `spamScore`) están duplicadas en ambos archivos.

**Impacto**

Aumenta la probabilidad de inconsistencias y dificulta mantenimiento.

**Recomendación**

Extraer funciones comunes a un archivo `scoring.go` y reutilizarlas.

### GAP 6: Manejo incompleto de errores en creación y escritura de archivos

**Descripción**

Errores ignorados con `_` en `os.Create` y escrituras CSV sin comprobar errores.

**Recomendación**

Validar errores explícitamente y comprobar `Flush()`:

```go
cleanFile, err := os.Create("../dataset/dataset_final_secuencial.csv")
if err != nil { panic(err) }
if err := cleanWriter.Write(newRow); err != nil { panic(err) }
cleanWriter.Flush()
if err := cleanWriter.Error(); err != nil { panic(err) }
```

### GAP 7: Uso de `panic` como manejo general de errores

**Descripción**

Se usa `panic(err)` en varias partes.

**Recomendación**

Usar logging claro y manejo controlado de errores (o helper `checkError`) en lugar de `panic` para mayor robustez.

### GAP 8: Rutas relativas fijas

**Descripción**

Rutas como `"../dataset/dataset_clean.csv"` hacen depender la ejecución del directorio actual.

**Recomendación**

Permitir rutas por argumentos o configuraciones.

### GAP 9: Métrica de memoria capturada solo al final

**Descripción**

Se mide `runtime.ReadMemStats` al final; no se reporta pico de memoria ni otras métricas relevantes.

**Recomendación**

Reportar `HeapAlloc`, `TotalAlloc`, `NumGC`, etc., para caracterizar mejor uso de memoria.

### GAP 10: No se calculan métricas de precisión del detector

**Descripción**

No se calculan precision/recall/F1 pese a que el dataset sintético contiene etiqueta `is_synthetic_spam`.

**Recomendación**

Preservar la etiqueta y calcular TP/FP/FN/TN para generar precision, recall y F1.

### GAP 11: Pérdida de columnas útiles en la limpieza

**Descripción**

>`dataset_clean.csv` solo conserva columnas básicas y descarta campos útiles como `is_synthetic_spam`.

**Recomendación**

Extender `Record` para conservar columnas de evaluación y escribirlas en el dataset limpio.

### GAP 12: `CLUSTER_SPAM_PERCENT` no usado

**Descripción**

Constante definida en `ampliacion_de_datos.py` pero no aplicada.

**Recomendación**

Eliminar o implementar la lógica de clusters de spam (IP, timestamps, usuarios).

### GAP 13: Simulación incompleta de duplicados

**Descripción**

La generación marca `tipo_duplicado` pero no reutiliza textos previos para duplicados exactos.

**Recomendación**

Mantener buffer de textos previos y reutilizarlos cuando se genere un duplicado exacto.

### GAP 14: Detección de símbolos poco clara

**Descripción**

`demasiadosSimbolos` usa una cadena con puntos repetidos.

**Recomendación**

Definir una constante clara o usar `unicode.IsPunct`.

### GAP 15: Umbrales heurísticos rígidos

**Descripción**

Los umbrales (`>50`, `>100`, `>0.75`, `<0.45`) son estáticos.

**Recomendación**

Hacerlos configurables mediante constantes o argumentos.

### GAP 16: Modelo Promela no modela la clasificación

**Descripción**

El modelo representa limpieza pero no la fase de classification/worker pool.

**Recomendación**

Agregar proceso `Classifier` al modelo Promela para ampliar la verificación.

### GAP 17: Propiedad LTL demasiado básica

**Descripción**

La LTL `NO_NEGATIVOS` solo verifica no negatividad de contadores.

**Recomendación**

Agregar aserciones más fuertes, p. ej. `assert(read_count == rejected_count + validated_count)`.

### GAP 18: Ausencia de pruebas unitarias

**Descripción**

No hay tests para funciones críticas.

**Recomendación**

Agregar tests en Go para `capsRatio`, `lexicalDiversity`, `spamScore`, etc.

### GAP 19: Benchmark no automatiza escenarios de 4/8/16 workers

**Recomendación**

Modificar `benchmark.go` para lanzar automáticamente `concurrente` con 4, 8 y 16 workers y recoger resultados.

### GAP 20: Falta de documentación de ejecución

**Recomendación**

Agregar instrucciones reproducibles en README para generación, limpieza, clasificación, benchmark y verificación con SPIN.

## 7. Análisis de seguridad

El proyecto usa datos sintéticos; riesgo de exposición bajo. Recomendaciones:

- Validar y permitir configuración de rutas.
- Verificar errores de I/O.
- Documentar que los datos son sintéticos.
- Evitar subir CSVs pesados al repositorio.
- Usar `net.ParseIP` (ya presente) para validar IPs.

## 8. Análisis de concurrencia

### 8.1 Limpieza concurrente (archivo: `main.go`)
Patrón: Pipeline (Reader → Normalizers → Validators → Writers). Mecanismos: goroutines, channels, WaitGroup, Mutex. Evaluación: diseño correcto; cierre de canales realizado tras finalizar workers.

### 8.2 Clasificación concurrente (archivo: `concurrente.go`)
Patrón: Worker Pool. Evaluación: adecuado; posible mejora paralelizar conteo de frecuencias.

## 9. Análisis de rendimiento

Resultados observados (resumen):

| Escenario | Secuencial | Concurrente | Speedup |
|---:|---:|---:|---:|
| 4 workers | 24.84 s | 14.70 s | 1.69x |
| 8 workers | 22.77 s | 10.57 s | 2.15x |
| 16 workers | 24.62 s | 11.73 s | 2.10x |

El mejor resultado se obtuvo con 8 workers (speedup 2.15x). La saturación y overhead explican por qué 16 workers no mejora respecto a 8.

## 10. Análisis de uso de memoria

Mediciones observadas (resumen):

| Escenario | Secuencial (MB) | Concurrente (MB) |
|---:|---:|---:|
| 4 workers | 1073.95 | 1003.52 |
| 8 workers | 1076.05 | 987.81 |
| 16 workers | 1096.57 | 1253.69 |

Causa principal: `var records []Record` mantiene el dataset completo en memoria.

## 11. Recomendaciones priorizadas

Prioridad alta:

- Compilar antes de medir benchmark.
- Configurar workers por argumento.
- Validar errores de archivos y escrituras CSV.
- Conservar `is_synthetic_spam` en el dataset limpio para evaluación.

Prioridad media:

- Refactorizar funciones duplicadas.
- Procesar por chunks o streaming.
- Paralelizar conteo de frecuencias.
- Agregar pruebas unitarias.

Prioridad baja:

- Ampliar modelo Promela con `Classifier`.
- Documentar comandos en README.

## 12. Conclusión técnica

El proyecto es funcional y demuestra un uso adecuado de patrones de concurrencia y verificación formal parcial. Identificados GAPs (benchmark, manejo de errores, memoria, duplicación y métricas de calidad) deben abordarse para convertir esta entrega académica en una solución reproducible y mantenible.

---

Archivo generado automáticamente por el análisis y convertido a Markdown.

# Informe Técnico — Análisis de GAPs con IA

**Proyecto:** Sistema concurrente para la detección de spam y anomalías en registros de reclamos  
**Curso:** Programación Concurrente y Distribuida  
**Fecha:** 2026-05-10  
**Modelo de IA utilizado:** GPT-5.5 Thinking  

---

## 1. Objetivo del análisis

El objetivo de este informe es analizar técnicamente el código fuente del proyecto, identificando fortalezas y GAPs relacionados con calidad de código, seguridad, manejo de errores, patrones de concurrencia, rendimiento, uso de memoria, escalabilidad, mantenibilidad y verificación formal.

El análisis se realizó sobre un sistema implementado principalmente en Go, complementado con un script Python para generación de datos y un modelo Promela para verificación formal mediante SPIN.

---

## 2. Prompt utilizado

Actúa como revisor técnico de un proyecto del curso Programación Concurrente y Distribuida.

Analiza el código fuente de un sistema concurrente para detección de spam y anomalías en registros de reclamos. El proyecto está implementado principalmente en Go y contiene:

1. Un script Python para generar y ampliar un dataset experimental.
2. Un pipeline concurrente de limpieza en Go usando goroutines, channels, sync.WaitGroup y sync.Mutex.
3. Una versión secuencial del detector de spam.
4. Una versión concurrente del detector de spam usando worker pool.
5. Un benchmark para medir tiempo de ejecución, memoria y speedup.
6. Un modelo Promela para verificar sincronización, exclusión mutua y ausencia de deadlocks mediante SPIN.

Evalúa el código considerando calidad de código, seguridad, manejo de errores, patrones de concurrencia, posibles condiciones de carrera, uso de memoria, rendimiento, escalabilidad, mantenibilidad, coherencia entre implementación e informe y oportunidades de mejora.

Devuelve el resultado en formato Markdown con resumen general, fortalezas técnicas, GAPs encontrados, impacto de cada GAP, recomendaciones y conclusión técnica.

---

## 3. Archivos analizados

| Archivo | Descripción |
|---|---|
| `ampliacion_de_datos.py` | Genera y amplía el dataset experimental hasta 1,700,000 registros. |
| `main.go` | Implementa el pipeline concurrente de limpieza de datos. |
| `secuencial.go` | Implementa la versión secuencial del detector de spam. |
| `concurrente.go` | Implementa la versión concurrente del detector de spam mediante worker pool. |
| `benchmark.go` | Ejecuta pruebas de rendimiento, memoria y cálculo de speedup. |
| `modelo.pml` | Modela en Promela el pipeline concurrente para verificación con SPIN. |

---

## 4. Resumen general del análisis

El proyecto presenta una solución funcional y coherente con el objetivo del trabajo. Se implementó un flujo completo compuesto por generación de datos, limpieza concurrente, clasificación secuencial, clasificación concurrente, benchmarking y modelado formal en Promela.

La solución utiliza correctamente conceptos del curso, tales como `goroutines`, `channels`, `sync.WaitGroup`, `sync.Mutex`, pipeline, worker pool, speedup, media recortada y verificación formal con SPIN.

El sistema logra procesar un dataset ampliado de 1,700,000 registros, de los cuales 1,489,190 registros fueron considerados válidos después de la limpieza. La versión concurrente del detector obtuvo un speedup máximo de 2.15x con 8 workers, lo que demuestra una mejora significativa frente a la versión secuencial.

Sin embargo, el análisis también permitió identificar algunos GAPs técnicos relacionados con manejo de errores, uso de memoria, duplicación de código, parametrización del número de workers, medición del benchmark, evaluación de precisión y alcance del modelo Promela.

---

## 5. Fortalezas técnicas identificadas

### 5.1 Flujo completo del sistema

El proyecto no se limita únicamente a clasificar registros, sino que implementa un flujo completo:

- Generación y ampliación del dataset.
- Limpieza concurrente de datos.
- Clasificación secuencial.
- Clasificación concurrente.
- Medición de rendimiento.
- Modelado formal en Promela.
- Verificación con SPIN.

Esto demuestra una solución integral y alineada con los entregables solicitados.

---

### 5.2 Uso adecuado de programación concurrente en Go

El código utiliza correctamente elementos propios de Go para programación concurrente:

- `goroutines`
- `channels`
- `sync.WaitGroup`
- `sync.Mutex`

En el archivo `main.go`, el proceso de limpieza está organizado como un pipeline concurrente:

```go
rawChan := make(chan Record, 1000)
cleanChan := make(chan Record, 1000)
validChan := make(chan Record, 1000)
rejectChan := make(chan RejectedRecord, 1000)