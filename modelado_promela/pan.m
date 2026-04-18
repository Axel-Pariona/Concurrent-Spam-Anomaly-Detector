#define rand	pan_rand
#define pthread_equal(a,b)	((a)==(b))
#if defined(HAS_CODE) && defined(VERBOSE)
	#ifdef BFS_PAR
		bfs_printf("Pr: %d Tr: %d\n", II, t->forw);
	#else
		cpu_printf("Pr: %d Tr: %d\n", II, t->forw);
	#endif
#endif
	switch (t->forw) {
	default: Uerror("bad forward move");
	case 0:	/* if without executable clauses */
		continue;
	case 1: /* generic 'goto' or 'skip' */
		IfNotBlocked
		_m = 3; goto P999;
	case 2: /* generic 'else' */
		IfNotBlocked
		if (trpt->o_pm&1) continue;
		_m = 3; goto P999;

		 /* PROC :init: */
	case 3: // STATE 1 - modelo.pml:65 - [(run Reader())] (0:0:0 - 1)
		IfNotBlocked
		reached[4][1] = 1;
		if (!(addproc(II, 1, 0)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 4: // STATE 2 - modelo.pml:66 - [(run Normalizer())] (0:0:0 - 1)
		IfNotBlocked
		reached[4][2] = 1;
		if (!(addproc(II, 1, 1)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 5: // STATE 3 - modelo.pml:67 - [(run Validator())] (0:0:0 - 1)
		IfNotBlocked
		reached[4][3] = 1;
		if (!(addproc(II, 1, 2)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 6: // STATE 4 - modelo.pml:68 - [(run Deduplicator())] (0:0:0 - 1)
		IfNotBlocked
		reached[4][4] = 1;
		if (!(addproc(II, 1, 3)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 7: // STATE 5 - modelo.pml:69 - [-end-] (0:0:0 - 1)
		IfNotBlocked
		reached[4][5] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */

		 /* PROC Deduplicator */
	case 8: // STATE 1 - modelo.pml:52 - [ch_val_dedup?msg] (0:0:1 - 1)
		reached[3][1] = 1;
		if (q_len(now.ch_val_dedup) == 0) continue;

		XX=1;
		(trpt+1)->bup.oval = ((P3 *)_this)->msg;
		;
		((P3 *)_this)->msg = qrecv(now.ch_val_dedup, XX-1, 0, 1);
#ifdef VAR_RANGES
		logval("Deduplicator:msg", ((P3 *)_this)->msg);
#endif
		;
		
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[32];
			sprintf(simvals, "%d?", now.ch_val_dedup);
		sprintf(simtmp, "%d", ((P3 *)_this)->msg); strcat(simvals, simtmp);		}
#endif
		;
		_m = 4; goto P999; /* 0 */
	case 9: // STATE 2 - modelo.pml:54 - [((msg==DATA))] (0:0:1 - 1)
		IfNotBlocked
		reached[3][2] = 1;
		if (!((((P3 *)_this)->msg==2)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P3 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P3 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 10: // STATE 3 - modelo.pml:56 - [processed = (processed+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[3][3] = 1;
		(trpt+1)->bup.oval = processed;
		processed = (processed+1);
#ifdef VAR_RANGES
		logval("processed", processed);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 11: // STATE 5 - modelo.pml:58 - [((msg==END))] (0:0:1 - 1)
		IfNotBlocked
		reached[3][5] = 1;
		if (!((((P3 *)_this)->msg==1)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P3 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P3 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 12: // STATE 12 - modelo.pml:62 - [-end-] (0:0:0 - 3)
		IfNotBlocked
		reached[3][12] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */

		 /* PROC Validator */
	case 13: // STATE 1 - modelo.pml:38 - [ch_norm_val?msg] (0:0:1 - 1)
		reached[2][1] = 1;
		if (q_len(now.ch_norm_val) == 0) continue;

		XX=1;
		(trpt+1)->bup.oval = ((P2 *)_this)->msg;
		;
		((P2 *)_this)->msg = qrecv(now.ch_norm_val, XX-1, 0, 1);
#ifdef VAR_RANGES
		logval("Validator:msg", ((P2 *)_this)->msg);
#endif
		;
		
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[32];
			sprintf(simvals, "%d?", now.ch_norm_val);
		sprintf(simtmp, "%d", ((P2 *)_this)->msg); strcat(simvals, simtmp);		}
#endif
		;
		_m = 4; goto P999; /* 0 */
	case 14: // STATE 2 - modelo.pml:40 - [((msg==DATA))] (0:0:1 - 1)
		IfNotBlocked
		reached[2][2] = 1;
		if (!((((P2 *)_this)->msg==2)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P2 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P2 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 15: // STATE 3 - modelo.pml:41 - [ch_val_dedup!DATA] (0:0:0 - 1)
		IfNotBlocked
		reached[2][3] = 1;
		if (q_full(now.ch_val_dedup))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_val_dedup);
		sprintf(simtmp, "%d", 2); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_val_dedup, 0, 2, 1);
		_m = 2; goto P999; /* 0 */
	case 16: // STATE 4 - modelo.pml:42 - [((msg==END))] (0:0:1 - 1)
		IfNotBlocked
		reached[2][4] = 1;
		if (!((((P2 *)_this)->msg==1)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P2 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P2 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 17: // STATE 5 - modelo.pml:43 - [ch_val_dedup!END] (0:0:0 - 1)
		IfNotBlocked
		reached[2][5] = 1;
		if (q_full(now.ch_val_dedup))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_val_dedup);
		sprintf(simtmp, "%d", 1); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_val_dedup, 0, 1, 1);
		_m = 2; goto P999; /* 0 */
	case 18: // STATE 12 - modelo.pml:47 - [-end-] (0:0:0 - 3)
		IfNotBlocked
		reached[2][12] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */

		 /* PROC Normalizer */
	case 19: // STATE 1 - modelo.pml:24 - [ch_reader_norm?msg] (0:0:1 - 1)
		reached[1][1] = 1;
		if (q_len(now.ch_reader_norm) == 0) continue;

		XX=1;
		(trpt+1)->bup.oval = ((P1 *)_this)->msg;
		;
		((P1 *)_this)->msg = qrecv(now.ch_reader_norm, XX-1, 0, 1);
#ifdef VAR_RANGES
		logval("Normalizer:msg", ((P1 *)_this)->msg);
#endif
		;
		
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[32];
			sprintf(simvals, "%d?", now.ch_reader_norm);
		sprintf(simtmp, "%d", ((P1 *)_this)->msg); strcat(simvals, simtmp);		}
#endif
		;
		_m = 4; goto P999; /* 0 */
	case 20: // STATE 2 - modelo.pml:26 - [((msg==DATA))] (0:0:1 - 1)
		IfNotBlocked
		reached[1][2] = 1;
		if (!((((P1 *)_this)->msg==2)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P1 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P1 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 21: // STATE 3 - modelo.pml:27 - [ch_norm_val!DATA] (0:0:0 - 1)
		IfNotBlocked
		reached[1][3] = 1;
		if (q_full(now.ch_norm_val))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_norm_val);
		sprintf(simtmp, "%d", 2); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_norm_val, 0, 2, 1);
		_m = 2; goto P999; /* 0 */
	case 22: // STATE 4 - modelo.pml:28 - [((msg==END))] (0:0:1 - 1)
		IfNotBlocked
		reached[1][4] = 1;
		if (!((((P1 *)_this)->msg==1)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P1 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P1 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 23: // STATE 5 - modelo.pml:29 - [ch_norm_val!END] (0:0:0 - 1)
		IfNotBlocked
		reached[1][5] = 1;
		if (q_full(now.ch_norm_val))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_norm_val);
		sprintf(simtmp, "%d", 1); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_norm_val, 0, 1, 1);
		_m = 2; goto P999; /* 0 */
	case 24: // STATE 12 - modelo.pml:33 - [-end-] (0:0:0 - 3)
		IfNotBlocked
		reached[1][12] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */

		 /* PROC Reader */
	case 25: // STATE 1 - modelo.pml:12 - [((i<10))] (0:0:0 - 1)
		IfNotBlocked
		reached[0][1] = 1;
		if (!((((P0 *)_this)->i<10)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 26: // STATE 2 - modelo.pml:13 - [ch_reader_norm!DATA] (0:0:0 - 1)
		IfNotBlocked
		reached[0][2] = 1;
		if (q_full(now.ch_reader_norm))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_reader_norm);
		sprintf(simtmp, "%d", 2); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_reader_norm, 0, 2, 1);
		_m = 2; goto P999; /* 0 */
	case 27: // STATE 3 - modelo.pml:14 - [i = (i+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[0][3] = 1;
		(trpt+1)->bup.oval = ((P0 *)_this)->i;
		((P0 *)_this)->i = (((P0 *)_this)->i+1);
#ifdef VAR_RANGES
		logval("Reader:i", ((P0 *)_this)->i);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 28: // STATE 5 - modelo.pml:16 - [ch_reader_norm!END] (0:0:0 - 1)
		IfNotBlocked
		reached[0][5] = 1;
		if (q_full(now.ch_reader_norm))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_reader_norm);
		sprintf(simtmp, "%d", 1); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_reader_norm, 0, 1, 1);
		_m = 2; goto P999; /* 0 */
	case 29: // STATE 10 - modelo.pml:19 - [-end-] (0:0:0 - 3)
		IfNotBlocked
		reached[0][10] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */
	case  _T5:	/* np_ */
		if (!((!(trpt->o_pm&4) && !(trpt->tau&128))))
			continue;
		/* else fall through */
	case  _T2:	/* true */
		_m = 3; goto P999;
#undef rand
	}

