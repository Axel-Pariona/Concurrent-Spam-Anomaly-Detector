# Detector de Spam y Anomalías (Concurrente)

Proyecto del curso **Programación Concurrente y Distribuida**  
Carrera de Ciencias de la Computación

---

## Descripción

Este proyecto implementa un sistema concurrente para la detección de anomalías y posibles patrones de spam en registros de reclamos, utilizando:

- Procesamiento concurrente en **Go**
- Técnicas de **preprocesamiento de datos**
- Pipeline tipo **productor-consumidor**
- Modelado formal en **Promela (SPIN)**

El sistema está diseñado para procesar grandes volúmenes de datos de forma eficiente, simulando escenarios reales con más de **1 millón de registros**.

---

## Arquitectura del sistema

El sistema sigue un modelo de pipeline concurrente:

- Reader → Normalizer → Validator → Deduplicator → Output


Cada etapa:
- Se ejecuta en paralelo (goroutines)
- Se comunica mediante canales
- Procesa datos de forma independiente

---

## Tecnologías utilizadas

- Go (Golang) → Concurrencia (goroutines, channels)
- CSV → Dataset de reclamos
- Promela → Modelado formal del sistema
- SPIN → Verificación de concurrencia

---

## Estructura del proyecto

- ampliacion_dataset/ → Generación de datos sintéticos
- codigo_fuente/ → Implementación en Go
- dataset/ → Dataset original y procesado
- modelado_promela/ → Modelo formal en Promela
- README.md → Documentación del proyecto


---

## Funcionalidades principales

- Procesamiento concurrente de datos  
- Normalización de registros  
- Validación de datos (IP, timestamp, campos)  
- Eliminación de duplicados  
- Generación de métricas del proceso

---

## Ejemplo de resultados

- Total Read: 1000000
- Final Clean: 937519
- Invalid IP Removed: 13531
- Duplicates Removed: 28850


---

## Modelado en Promela

Se modeló el sistema para verificar:

- Ausencia de deadlocks  
- Progreso del sistema (liveness)  
- Terminación correcta  
- Sincronización entre procesos  

---

## Integrantes

- Omar Junior Acuña Villegas  
- Axel Yamir Pariona Rojas  
- Rafael Tomas Chui Sánchez

---

## 📎 Repositorio

👉 https://github.com/Axel-Pariona/Detector-de-Spam---Programaci-n-concurrente-y-distribuida

---

## 📌 Notas

- Se utilizaron datos sintéticos para simular texto libre
- El enfoque principal está en concurrencia, no en ML avanzado
