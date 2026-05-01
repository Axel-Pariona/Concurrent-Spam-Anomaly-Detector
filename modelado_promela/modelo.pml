mtype = { DATA, END };

/* Canales para pipeline de procesamiento con buffer de 100 */
chan ch_reader_norm = [100] of { mtype };
chan ch_norm_val   = [100] of { mtype };
chan ch_val_dedup  = [100] of { mtype };
chan ch_dedup_class = [100] of { mtype };

/* Contadores de registros */
int processed = 0;
int rejected = 0;
int spam_detected = 0;

/* Buffer pequeño para 1M de registros (escala reducida para verificación) */
#define TOTAL_RECORDS 1000000
#define SCALE_FACTOR 100  /* Reducido para simulación más rápida */

proctype Reader() {
    int i = 0;
    int max_records = TOTAL_RECORDS / SCALE_FACTOR;
    
    do
    :: i < max_records ->
        ch_reader_norm!DATA;
        i++
    :: else ->
        ch_reader_norm!END;
        break
    od
}

proctype Normalizer() {
    mtype msg;
    int invalid_count = 0;
    
    do
    :: ch_reader_norm?msg ->
        if
        :: msg == DATA ->
            /* Simular 2% IPs inválidas + 2% timestamps inválidos + 1% campos vacíos = 5% rechazo */
            if
            :: (invalid_count % 20 == 0) -> /* ~5% rechazo */
                atomic {
                    rejected++
                }
            :: else ->
                ch_norm_val!DATA
            fi;
            invalid_count++
        :: msg == END ->
            ch_norm_val!END;
            break
        fi
    od
}

proctype Validator() {
    mtype msg;
    int spam_count = 0;
    
    do
    :: ch_norm_val?msg ->
        if
        :: msg == DATA ->
            /* Simular 2% spam + 2% duplicados + 8% ruido = ~12% variantes detectables */
            if
            :: (spam_count % 10 == 0) -> /* ~10% detección de spam/clusters */
                atomic {
                    spam_detected++
                };
                ch_val_dedup!DATA
            :: else ->
                ch_val_dedup!DATA
            fi;
            spam_count++
        :: msg == END ->
            ch_val_dedup!END;
            break
        fi
    od
}

proctype Deduplicator() {
    mtype msg;
    int dedup_count = 0;
    
    do
    :: ch_val_dedup?msg ->
        if
        :: msg == DATA ->
            /* Simular 2% duplicados que se eliminan */
            if
            :: (dedup_count % 50 == 0) -> /* ~2% duplicados descartados */
                skip  /* registro duplicado descartado */
            :: else ->
                atomic {
                    processed++
                }
            fi;
            dedup_count++
        :: msg == END ->
            printf("=== RESULTADOS FINALES ===\n");
            printf("Registros procesados: %d\n", processed);
            printf("Registros rechazados: %d\n", rejected);
            printf("Spam detectado: %d\n", spam_detected);
            printf("Total esperado: %d\n", (TOTAL_RECORDS / SCALE_FACTOR));
            break
        fi
    od
}

init {
    run Reader();
    run Normalizer();
    run Validator();
    run Deduplicator();
}
