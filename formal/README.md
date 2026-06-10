# Formal Verification

Esta carpeta contiene los modelos Promela utilizados para analizar formalmente una abstracción del pipeline concurrente del proyecto **Concurrent Spam Anomaly Detector**.

## Archivos

```txt
modelo.pml
modelo_deadlock.pml
```

## Objetivo

El objetivo del modelado formal es representar el comportamiento principal del sistema concurrente usando Promela y verificar propiedades relacionadas con:

- Comunicación entre procesos.
- Uso de canales.
- Propagación de señales de finalización.
- Ausencia de deadlocks.
- Sincronización entre etapas.
- Valores válidos en contadores compartidos.

## Modelo principal

El archivo:

```txt
modelo.pml
```

representa una versión abstracta del pipeline concurrente.

Incluye:

- Un proceso `Reader`.
- Múltiples procesos `Normalizer`.
- Múltiples procesos `Validator`.
- Canales `ch_raw` y `ch_clean`.
- Mensajes `DATA` y `END`.
- Contadores globales.
- Propiedad LTL para validar que los contadores no sean negativos.

## Modelo de análisis de deadlock

El archivo:

```txt
modelo_deadlock.pml
```

se usa para explorar escenarios relacionados con sincronización y finalización del pipeline.

Permite analizar si la interacción entre procesos puede bloquearse en determinadas condiciones.

## Relación con el código Go

El modelo Promela no representa todo el código Go línea por línea.

En cambio, abstrae el comportamiento concurrente principal:

```txt
Reader
  ↓
Normalizer
  ↓
Validator
```

Esta abstracción permite verificar propiedades de concurrencia sin depender de detalles de implementación como lectura CSV, limpieza de texto o escritura de archivos.

## Requisitos

Para ejecutar los modelos se necesita SPIN.

Verificar instalación:

```bash
spin -V
```

También se requiere un compilador C si se desea generar el verificador `pan`.

En Linux o WSL normalmente se usa `gcc`.

En Windows se puede usar MinGW o ejecutar desde WSL.

## Simulación simple

Desde esta carpeta:

```bash
spin modelo.pml
```

También puedes ejecutar:

```bash
spin modelo_deadlock.pml
```

## Verificación completa

Generar el verificador:

```bash
spin -a modelo.pml
```

Compilar:

```bash
gcc -o pan pan.c
```

Ejecutar:

```bash
./pan
```

En Windows PowerShell con MinGW:

```powershell
spin -a modelo.pml
gcc -o pan.exe pan.c
.\pan.exe
```

## Verificación de propiedad LTL

El modelo principal incluye una propiedad LTL orientada a validar que los contadores se mantengan en valores no negativos.

Para verificar propiedades LTL, se puede generar el verificador con SPIN y ejecutar `pan` según el flujo estándar.

## Archivos generados por SPIN

SPIN puede generar archivos como:

```txt
pan
pan.c
pan.h
pan.m
pan.p
pan.t
_spin_nvr.tmp
*.trail
```

Estos archivos no deben subirse al repositorio porque son artefactos generados automáticamente.

El `.gitignore` ya los excluye.

## Interpretación de resultados

Si SPIN no reporta errores, significa que no encontró violaciones bajo el modelo y configuración analizados.

Si se genera un archivo `.trail`, significa que SPIN encontró una traza de error. Esa traza puede revisarse con:

```bash
spin -t modelo.pml
```

## Limitaciones

- El modelo es una abstracción del sistema real.
- No incluye detalles de lectura o escritura de archivos.
- No modela todas las reglas de detección de spam.
- No representa el benchmark.
- Su objetivo es analizar sincronización y comunicación, no lógica completa de negocio.

## Conclusión

Los modelos Promela permiten complementar la implementación en Go con una revisión formal del comportamiento concurrente. Esto fortalece el análisis académico del proyecto al incluir no solo implementación y benchmark, sino también verificación de propiedades de concurrencia.
