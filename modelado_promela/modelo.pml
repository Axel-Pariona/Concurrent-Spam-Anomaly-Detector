mtype = { DATA, END };

#define TOTAL_RECORDS 6
#define N_NORMALIZERS 2
#define N_VALIDATORS  2

chan ch_raw   = [20] of { mtype };
chan ch_clean = [20] of { mtype };

int read_count = 0;
int normalized_count = 0;
int rejected_count = 0;
int validated_count = 0;

int norm_done = 0;
int val_done = 0;

proctype Reader() {

    int i;

    i = 0;

    do
    :: i < TOTAL_RECORDS ->

        ch_raw!DATA;

        atomic {
            read_count = read_count + 1
        };

        i = i + 1

    :: else ->

        i = 0;

        do
        :: i < N_NORMALIZERS ->
            ch_raw!END;
            i = i + 1
        :: else ->
            break
        od;

        break
    od
}

proctype Normalizer() {

    mtype msg;
    int c;

    c = 0;

    do
    :: ch_raw?msg ->

        if
        :: msg == DATA ->

            if
            :: (c % 8 == 0) ->

                atomic {
                    rejected_count = rejected_count + 1
                }

            :: else ->

                atomic {
                    normalized_count =
                    normalized_count + 1
                };

                ch_clean!DATA
            fi;

            c = c + 1

        :: msg == END ->

            atomic {
                norm_done = norm_done + 1
            };

            if
            :: norm_done == N_NORMALIZERS ->

                c = 0;

                do
                :: c < N_VALIDATORS ->
                    ch_clean!END;
                    c = c + 1
                :: else ->
                    break
                od

            :: else ->
                skip
            fi;

            break
        fi
    od
}

proctype Validator() {

    mtype msg;

    do
    :: ch_clean?msg ->

        if
        :: msg == DATA ->

            atomic {
                validated_count =
                validated_count + 1
            }

        :: msg == END ->

            atomic {
                val_done = val_done + 1
            };

            break
        fi
    od
}

init {

    int i;

    i = 0;

    run Reader();

    do
    :: i < N_NORMALIZERS ->
        run Normalizer();
        i = i + 1
    :: else ->
        break
    od;

    i = 0;

    do
    :: i < N_VALIDATORS ->
        run Validator();
        i = i + 1
    :: else ->
        break
    od
}

ltl NO_NEGATIVOS {
    [] (
        read_count >= 0 &&
        rejected_count >= 0 &&
        validated_count >= 0
    )
}