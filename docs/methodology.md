# Methodology

Este documento describe la metodología aplicada en el proyecto **Concurrent Spam Anomaly Detector**.

El proyecto implementa un flujo completo para generar, limpiar, procesar y clasificar registros de reclamos usando versiones secuencial y concurrente en Go.

## 1. Planteamiento del problema

El problema consiste en procesar un gran volumen de registros de reclamos para identificar posibles patrones de spam o anomalías.

Los registros pueden contener:

- Texto repetido.
- Mensajes demasiado cortos.
- Exceso de mayúsculas.
- Exceso de símbolos.
- Direcciones IP con alta frecuencia.
- Texto artificial.
- Baja diversidad léxica.
- Campos incompletos o inválidos.

El reto principal es procesar estos datos de forma eficiente usando programación concurrente.

## 2. Generación del dataset

El proceso inicia con un dataset base:

```txt
data/dataset_base.csv
```

A partir de este archivo, el script:

```txt
scripts/generate_dataset.py
```

genera un dataset ampliado:

```txt
data/dataset_1M_raw.csv
```

El dataset ampliado incluye registros legítimos, variaciones, duplicados, ruido, campos vacíos, timestamps inválidos, IPs inválidas y patrones artificiales de spam.

Esta etapa permite simular un volumen grande de datos para evaluar el comportamiento de las versiones secuencial y concurrente.

## 3. Limpieza de datos

La limpieza se implementa en:

```txt
src/cleaner_pipeline.go
```

Este archivo usa un pipeline concurrente que divide el procesamiento en etapas.

### Etapas del pipeline

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

### Reader

Lee registros desde el CSV de entrada y los envía al pipeline.

### Normalizer

Limpia texto, elimina espacios innecesarios y normaliza campos.

### Validator

Descarta registros inválidos, por ejemplo:

- Texto vacío.
- Texto muy corto.
- Texto con demasiados símbolos.
- IP inválida.
- Timestamp inválido.

### Deduplicator

Detecta y elimina duplicados.

### Output

Guarda los registros válidos y rechazados en archivos CSV.

## 4. Comunicación concurrente

Las etapas del pipeline se comunican mediante channels de Go.

Esto permite que diferentes partes del proceso trabajen simultáneamente sobre distintos registros.

El flujo concurrente reduce el tiempo de espera entre lectura, validación y escritura cuando el volumen de datos es alto.

## 5. Métricas de limpieza

Durante la limpieza se registran métricas como:

- Total de registros leídos.
- Total de registros procesados.
- Total de registros descartados.
- Textos vacíos.
- Textos limpiados.
- Tiempo total del proceso.

El resumen se guarda en:

```txt
logs/cleaning_summary.txt
```

## 6. Detección de spam

La detección se implementa en dos versiones:

```txt
src/spam_sequential.go
src/spam_concurrent.go
```

Ambas versiones usan reglas heurísticas para asignar un puntaje a cada registro.

## 7. Reglas heurísticas

El puntaje de spam se calcula evaluando señales como:

### Texto muy repetido

Si un mismo texto aparece muchas veces, se considera una señal de posible spam.

### IP abusiva

Si una IP aparece en demasiados registros, se considera una señal sospechosa.

### Exceso de mayúsculas

Los mensajes con proporción alta de mayúsculas pueden indicar texto artificial o insistente.

### Baja diversidad léxica

Un texto con pocas palabras únicas en relación con su longitud puede ser repetitivo.

### Palabras asociadas a spam

El sistema busca palabras como:

```txt
URGENTE
ESTAFA
DINERO
GRATIS
PREMIO
CLICK
FRAUDE
OFERTA
GANA
```

### Texto demasiado corto

Mensajes extremadamente cortos pueden no contener suficiente información válida.

### Repetición excesiva

Palabras repetidas muchas veces dentro del mismo texto aumentan el puntaje.

### Exceso de símbolos

Mensajes con muchos símbolos pueden representar ruido o spam.

### Ausencia de conectores legítimos

La ausencia de conectores comunes puede indicar texto poco natural, dependiendo de la longitud del mensaje.

## 8. Versión secuencial

La versión secuencial procesa los registros uno por uno.

Archivo:

```txt
src/spam_sequential.go
```

Flujo:

```txt
Leer dataset_clean.csv
  ↓
Calcular frecuencias
  ↓
Evaluar cada registro
  ↓
Asignar score
  ↓
Guardar resultado final
  ↓
Guardar registros spam
```

Esta versión sirve como línea base para comparar contra la implementación concurrente.

## 9. Versión concurrente

La versión concurrente implementa un worker pool.

Archivo:

```txt
src/spam_concurrent.go
```

Flujo:

```txt
Leer dataset_clean.csv
  ↓
Calcular frecuencias compartidas
  ↓
Enviar registros a jobs channel
  ↓
Procesar registros en workers
  ↓
Recibir resultados
  ↓
Guardar archivos finales
```

El número de workers puede basarse en la cantidad de CPUs disponibles.

## 10. Benchmark

El benchmark se implementa en:

```txt
src/benchmark.go
```

Este programa ejecuta varias veces:

```txt
spam_sequential.go
spam_concurrent.go
```

Luego calcula:

- Tiempo promedio.
- Memoria promedio.
- Media recortada.
- Speedup.

La media recortada elimina valores extremos para obtener una comparación más estable.

## 11. Speedup

El speedup se calcula como:

```txt
speedup = tiempo_secuencial / tiempo_concurrente
```

Si el resultado es mayor que 1, la versión concurrente fue más rápida.

Si el resultado es cercano a 1 o menor, puede significar que el overhead de concurrencia fue similar o mayor que el beneficio obtenido.

## 12. Uso de memoria

Cada versión imprime la memoria utilizada durante la ejecución. El benchmark extrae este valor del output y calcula un promedio.

Esto permite comparar no solo tiempo, sino también consumo de recursos.

## 13. Verificación formal

El proyecto incluye modelos Promela en:

```txt
formal/
```

El objetivo del modelado formal es representar de forma abstracta el comportamiento del pipeline concurrente.

Elementos modelados:

- Proceso lector.
- Procesos normalizadores.
- Procesos validadores.
- Canales de comunicación.
- Mensajes de datos.
- Mensajes de finalización.
- Contadores compartidos.

## 14. Propiedades verificadas

El modelo permite razonar sobre propiedades como:

- Ausencia de deadlock.
- Correcta propagación de señales de finalización.
- Contadores no negativos.
- Finalización del flujo de procesos.
- Sincronización entre etapas.

El modelo no representa cada detalle del código Go, sino una abstracción del comportamiento concurrente principal.

## 15. Organización de salidas

Los archivos generados se organizan en:

```txt
data/
logs/
```

`data/` contiene CSVs generados.

`logs/` contiene resúmenes de generación y limpieza.

Los archivos grandes generados no se versionan en GitHub porque pueden reconstruirse.

## 16. Limitaciones metodológicas

El proyecto tiene algunas limitaciones:

- La detección se basa en reglas heurísticas.
- No se entrena un modelo de Machine Learning.
- Los datos ampliados son sintéticos.
- Los resultados dependen del hardware.
- La concurrencia introduce overhead.
- La verificación formal es una abstracción del sistema real.
- Las reglas de spam pueden generar falsos positivos o falsos negativos.

## 17. Posibles mejoras metodológicas

Se podrían implementar mejoras como:

- Medición automática de precisión, recall y F1 si se agregan etiquetas reales.
- Exportación de resultados del benchmark a CSV.
- Configuración del número de workers mediante flags.
- Parametrización del tamaño del dataset.
- Uso de `context.Context` para cancelación controlada.
- Tests unitarios para cada regla de spam.
- Profiling con `pprof`.
- Separación del código Go en paquetes.
- Ampliación de propiedades LTL en Promela.
- Comparación con procesamiento distribuido.

## Conclusión

La metodología del proyecto permite evaluar el impacto de la programación concurrente sobre una tarea de procesamiento masivo de registros. La comparación entre versión secuencial y concurrente, junto con el modelado Promela/SPIN, proporciona una base sólida para analizar rendimiento, sincronización y diseño de sistemas concurrentes.
