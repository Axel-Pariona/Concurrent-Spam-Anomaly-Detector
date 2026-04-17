mtype = { DATA, END };

chan ch_reader_norm = [5] of { mtype };
chan ch_norm_val   = [5] of { mtype };
chan ch_val_dedup  = [5] of { mtype };

int processed = 0;

proctype Reader() {
    int i = 0;
    do
    :: i < 10 ->
        ch_reader_norm!DATA;
        i++
    :: else ->
        ch_reader_norm!END;
        break
    od
}

proctype Normalizer() {
    mtype msg;
    do
    :: ch_reader_norm?msg ->
        if
        :: msg == DATA ->
            ch_norm_val!DATA
        :: msg == END ->
            ch_norm_val!END;
            break
        fi
    od
}

proctype Validator() {
    mtype msg;
    do
    :: ch_norm_val?msg ->
        if
        :: msg == DATA ->
            ch_val_dedup!DATA
        :: msg == END ->
            ch_val_dedup!END;
            break
        fi
    od
}

proctype Deduplicator() {
    mtype msg;
    do
    :: ch_val_dedup?msg ->
        if
        :: msg == DATA ->
            atomic {
                processed++
            }
        :: msg == END ->
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
