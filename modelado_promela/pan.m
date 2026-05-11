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

		 /* CLAIM NO_NEGATIVOS */
	case 3: // STATE 1 - _spin_nvr.tmp:3 - [(!((((read_count>=0)&&(rejected_count>=0))&&(validated_count>=0))))] (6:0:0 - 1)
		
#if defined(VERI) && !defined(NP)
#if NCLAIMS>1
		{	static int reported1 = 0;
			if (verbose && !reported1)
			{	int nn = (int) ((Pclaim *)pptr(0))->_n;
				printf("depth %ld: Claim %s (%d), state %d (line %d)\n",
					depth, procname[spin_c_typ[nn]], nn, (int) ((Pclaim *)pptr(0))->_p, src_claim[ (int) ((Pclaim *)pptr(0))->_p ]);
				reported1 = 1;
				fflush(stdout);
		}	}
#else
		{	static int reported1 = 0;
			if (verbose && !reported1)
			{	printf("depth %d: Claim, state %d (line %d)\n",
					(int) depth, (int) ((Pclaim *)pptr(0))->_p, src_claim[ (int) ((Pclaim *)pptr(0))->_p ]);
				reported1 = 1;
				fflush(stdout);
		}	}
#endif
#endif
		reached[4][1] = 1;
		if (!( !((((now.read_count>=0)&&(now.rejected_count>=0))&&(now.validated_count>=0)))))
			continue;
		/* merge: assert(!(!((((read_count>=0)&&(rejected_count>=0))&&(validated_count>=0)))))(0, 2, 6) */
		reached[4][2] = 1;
		spin_assert( !( !((((now.read_count>=0)&&(now.rejected_count>=0))&&(now.validated_count>=0)))), " !( !((((read_count>=0)&&(rejected_count>=0))&&(validated_count>=0))))", II, tt, t);
		/* merge: .(goto)(0, 7, 6) */
		reached[4][7] = 1;
		;
		_m = 3; goto P999; /* 2 */
	case 4: // STATE 10 - _spin_nvr.tmp:8 - [-end-] (0:0:0 - 1)
		
#if defined(VERI) && !defined(NP)
#if NCLAIMS>1
		{	static int reported10 = 0;
			if (verbose && !reported10)
			{	int nn = (int) ((Pclaim *)pptr(0))->_n;
				printf("depth %ld: Claim %s (%d), state %d (line %d)\n",
					depth, procname[spin_c_typ[nn]], nn, (int) ((Pclaim *)pptr(0))->_p, src_claim[ (int) ((Pclaim *)pptr(0))->_p ]);
				reported10 = 1;
				fflush(stdout);
		}	}
#else
		{	static int reported10 = 0;
			if (verbose && !reported10)
			{	printf("depth %d: Claim, state %d (line %d)\n",
					(int) depth, (int) ((Pclaim *)pptr(0))->_p, src_claim[ (int) ((Pclaim *)pptr(0))->_p ]);
				reported10 = 1;
				fflush(stdout);
		}	}
#endif
#endif
		reached[4][10] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */

		 /* PROC :init: */
	case 5: // STATE 1 - modelo.pml:141 - [i = 0] (0:0:1 - 1)
		IfNotBlocked
		reached[3][1] = 1;
		(trpt+1)->bup.oval = ((P3 *)_this)->i;
		((P3 *)_this)->i = 0;
#ifdef VAR_RANGES
		logval(":init::i", ((P3 *)_this)->i);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 6: // STATE 2 - modelo.pml:143 - [(run Reader())] (0:0:0 - 1)
		IfNotBlocked
		reached[3][2] = 1;
		if (!(addproc(II, 1, 0)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 7: // STATE 3 - modelo.pml:146 - [((i<2))] (0:0:0 - 1)
		IfNotBlocked
		reached[3][3] = 1;
		if (!((((P3 *)_this)->i<2)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 8: // STATE 4 - modelo.pml:147 - [(run Normalizer())] (0:0:0 - 1)
		IfNotBlocked
		reached[3][4] = 1;
		if (!(addproc(II, 1, 1)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 9: // STATE 5 - modelo.pml:148 - [i = (i+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[3][5] = 1;
		(trpt+1)->bup.oval = ((P3 *)_this)->i;
		((P3 *)_this)->i = (((P3 *)_this)->i+1);
#ifdef VAR_RANGES
		logval(":init::i", ((P3 *)_this)->i);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 10: // STATE 11 - modelo.pml:153 - [i = 0] (0:17:1 - 3)
		IfNotBlocked
		reached[3][11] = 1;
		(trpt+1)->bup.oval = ((P3 *)_this)->i;
		((P3 *)_this)->i = 0;
#ifdef VAR_RANGES
		logval(":init::i", ((P3 *)_this)->i);
#endif
		;
		/* merge: .(goto)(0, 18, 17) */
		reached[3][18] = 1;
		;
		_m = 3; goto P999; /* 1 */
	case 11: // STATE 12 - modelo.pml:156 - [((i<2))] (0:0:0 - 1)
		IfNotBlocked
		reached[3][12] = 1;
		if (!((((P3 *)_this)->i<2)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 12: // STATE 13 - modelo.pml:157 - [(run Validator())] (0:0:0 - 1)
		IfNotBlocked
		reached[3][13] = 1;
		if (!(addproc(II, 1, 2)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 13: // STATE 14 - modelo.pml:158 - [i = (i+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[3][14] = 1;
		(trpt+1)->bup.oval = ((P3 *)_this)->i;
		((P3 *)_this)->i = (((P3 *)_this)->i+1);
#ifdef VAR_RANGES
		logval(":init::i", ((P3 *)_this)->i);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 14: // STATE 20 - modelo.pml:162 - [-end-] (0:0:0 - 3)
		IfNotBlocked
		reached[3][20] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */

		 /* PROC Validator */
	case 15: // STATE 1 - modelo.pml:116 - [ch_clean?msg] (0:0:1 - 1)
		reached[2][1] = 1;
		if (q_len(now.ch_clean) == 0) continue;

		XX=1;
		(trpt+1)->bup.oval = ((P2 *)_this)->msg;
		;
		((P2 *)_this)->msg = qrecv(now.ch_clean, XX-1, 0, 1);
#ifdef VAR_RANGES
		logval("Validator:msg", ((P2 *)_this)->msg);
#endif
		;
		
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[32];
			sprintf(simvals, "%d?", now.ch_clean);
		sprintf(simtmp, "%d", ((P2 *)_this)->msg); strcat(simvals, simtmp);		}
#endif
		;
		_m = 4; goto P999; /* 0 */
	case 16: // STATE 2 - modelo.pml:119 - [((msg==DATA))] (0:0:1 - 1)
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
	case 17: // STATE 3 - modelo.pml:122 - [validated_count = (validated_count+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[2][3] = 1;
		(trpt+1)->bup.oval = now.validated_count;
		now.validated_count = (now.validated_count+1);
#ifdef VAR_RANGES
		logval("validated_count", now.validated_count);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 18: // STATE 5 - modelo.pml:126 - [((msg==END))] (0:0:1 - 1)
		IfNotBlocked
		reached[2][5] = 1;
		if (!((((P2 *)_this)->msg==1)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P2 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P2 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 19: // STATE 6 - modelo.pml:129 - [val_done = (val_done+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[2][6] = 1;
		(trpt+1)->bup.oval = now.val_done;
		now.val_done = (now.val_done+1);
#ifdef VAR_RANGES
		logval("val_done", now.val_done);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 20: // STATE 14 - modelo.pml:135 - [-end-] (0:0:0 - 3)
		IfNotBlocked
		reached[2][14] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */

		 /* PROC Normalizer */
	case 21: // STATE 1 - modelo.pml:56 - [c = 0] (0:0:1 - 1)
		IfNotBlocked
		reached[1][1] = 1;
		(trpt+1)->bup.oval = ((P1 *)_this)->c;
		((P1 *)_this)->c = 0;
#ifdef VAR_RANGES
		logval("Normalizer:c", ((P1 *)_this)->c);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 22: // STATE 2 - modelo.pml:59 - [ch_raw?msg] (0:0:1 - 1)
		reached[1][2] = 1;
		if (q_len(now.ch_raw) == 0) continue;

		XX=1;
		(trpt+1)->bup.oval = ((P1 *)_this)->msg;
		;
		((P1 *)_this)->msg = qrecv(now.ch_raw, XX-1, 0, 1);
#ifdef VAR_RANGES
		logval("Normalizer:msg", ((P1 *)_this)->msg);
#endif
		;
		
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[32];
			sprintf(simvals, "%d?", now.ch_raw);
		sprintf(simtmp, "%d", ((P1 *)_this)->msg); strcat(simvals, simtmp);		}
#endif
		;
		_m = 4; goto P999; /* 0 */
	case 23: // STATE 3 - modelo.pml:62 - [((msg==DATA))] (0:0:1 - 1)
		IfNotBlocked
		reached[1][3] = 1;
		if (!((((P1 *)_this)->msg==2)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P1 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P1 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 24: // STATE 4 - modelo.pml:65 - [(((c%8)==0))] (0:0:0 - 1)
		IfNotBlocked
		reached[1][4] = 1;
		if (!(((((P1 *)_this)->c%8)==0)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 25: // STATE 5 - modelo.pml:68 - [rejected_count = (rejected_count+1)] (0:34:2 - 1)
		IfNotBlocked
		reached[1][5] = 1;
		(trpt+1)->bup.ovals = grab_ints(2);
		(trpt+1)->bup.ovals[0] = now.rejected_count;
		now.rejected_count = (now.rejected_count+1);
#ifdef VAR_RANGES
		logval("rejected_count", now.rejected_count);
#endif
		;
		/* merge: .(goto)(34, 12, 34) */
		reached[1][12] = 1;
		;
		/* merge: c = (c+1)(34, 13, 34) */
		reached[1][13] = 1;
		(trpt+1)->bup.ovals[1] = ((P1 *)_this)->c;
		((P1 *)_this)->c = (((P1 *)_this)->c+1);
#ifdef VAR_RANGES
		logval("Normalizer:c", ((P1 *)_this)->c);
#endif
		;
		/* merge: .(goto)(0, 33, 34) */
		reached[1][33] = 1;
		;
		/* merge: .(goto)(0, 35, 34) */
		reached[1][35] = 1;
		;
		_m = 3; goto P999; /* 4 */
	case 26: // STATE 8 - modelo.pml:74 - [normalized_count = (normalized_count+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[1][8] = 1;
		(trpt+1)->bup.oval = now.normalized_count;
		now.normalized_count = (now.normalized_count+1);
#ifdef VAR_RANGES
		logval("normalized_count", now.normalized_count);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 27: // STATE 10 - modelo.pml:78 - [ch_clean!DATA] (0:0:0 - 1)
		IfNotBlocked
		reached[1][10] = 1;
		if (q_full(now.ch_clean))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_clean);
		sprintf(simtmp, "%d", 2); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_clean, 0, 2, 1);
		_m = 2; goto P999; /* 0 */
	case 28: // STATE 13 - modelo.pml:81 - [c = (c+1)] (0:34:1 - 3)
		IfNotBlocked
		reached[1][13] = 1;
		(trpt+1)->bup.oval = ((P1 *)_this)->c;
		((P1 *)_this)->c = (((P1 *)_this)->c+1);
#ifdef VAR_RANGES
		logval("Normalizer:c", ((P1 *)_this)->c);
#endif
		;
		/* merge: .(goto)(0, 33, 34) */
		reached[1][33] = 1;
		;
		/* merge: .(goto)(0, 35, 34) */
		reached[1][35] = 1;
		;
		_m = 3; goto P999; /* 2 */
	case 29: // STATE 14 - modelo.pml:83 - [((msg==END))] (0:0:1 - 1)
		IfNotBlocked
		reached[1][14] = 1;
		if (!((((P1 *)_this)->msg==1)))
			continue;
		if (TstOnly) return 1; /* TT */
		/* dead 1: msg */  (trpt+1)->bup.oval = ((P1 *)_this)->msg;
#ifdef HAS_CODE
		if (!readtrail)
#endif
			((P1 *)_this)->msg = 0;
		_m = 3; goto P999; /* 0 */
	case 30: // STATE 15 - modelo.pml:86 - [norm_done = (norm_done+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[1][15] = 1;
		(trpt+1)->bup.oval = now.norm_done;
		now.norm_done = (now.norm_done+1);
#ifdef VAR_RANGES
		logval("norm_done", now.norm_done);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 31: // STATE 17 - modelo.pml:90 - [((norm_done==2))] (0:0:0 - 1)
		IfNotBlocked
		reached[1][17] = 1;
		if (!((now.norm_done==2)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 32: // STATE 18 - modelo.pml:92 - [c = 0] (0:0:1 - 1)
		IfNotBlocked
		reached[1][18] = 1;
		(trpt+1)->bup.oval = ((P1 *)_this)->c;
		((P1 *)_this)->c = 0;
#ifdef VAR_RANGES
		logval("Normalizer:c", ((P1 *)_this)->c);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 33: // STATE 19 - modelo.pml:95 - [((c<2))] (0:0:0 - 1)
		IfNotBlocked
		reached[1][19] = 1;
		if (!((((P1 *)_this)->c<2)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 34: // STATE 20 - modelo.pml:96 - [ch_clean!END] (0:0:0 - 1)
		IfNotBlocked
		reached[1][20] = 1;
		if (q_full(now.ch_clean))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_clean);
		sprintf(simtmp, "%d", 1); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_clean, 0, 1, 1);
		_m = 2; goto P999; /* 0 */
	case 35: // STATE 21 - modelo.pml:97 - [c = (c+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[1][21] = 1;
		(trpt+1)->bup.oval = ((P1 *)_this)->c;
		((P1 *)_this)->c = (((P1 *)_this)->c+1);
#ifdef VAR_RANGES
		logval("Normalizer:c", ((P1 *)_this)->c);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 36: // STATE 37 - modelo.pml:109 - [-end-] (0:0:0 - 7)
		IfNotBlocked
		reached[1][37] = 1;
		if (!delproc(1, II)) continue;
		_m = 3; goto P999; /* 0 */

		 /* PROC Reader */
	case 37: // STATE 1 - modelo.pml:22 - [i = 0] (0:0:1 - 1)
		IfNotBlocked
		reached[0][1] = 1;
		(trpt+1)->bup.oval = ((P0 *)_this)->i;
		((P0 *)_this)->i = 0;
#ifdef VAR_RANGES
		logval("Reader:i", ((P0 *)_this)->i);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 38: // STATE 2 - modelo.pml:25 - [((i<6))] (0:0:0 - 1)
		IfNotBlocked
		reached[0][2] = 1;
		if (!((((P0 *)_this)->i<6)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 39: // STATE 3 - modelo.pml:27 - [ch_raw!DATA] (0:0:0 - 1)
		IfNotBlocked
		reached[0][3] = 1;
		if (q_full(now.ch_raw))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_raw);
		sprintf(simtmp, "%d", 2); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_raw, 0, 2, 1);
		_m = 2; goto P999; /* 0 */
	case 40: // STATE 4 - modelo.pml:30 - [read_count = (read_count+1)] (0:18:2 - 1)
		IfNotBlocked
		reached[0][4] = 1;
		(trpt+1)->bup.ovals = grab_ints(2);
		(trpt+1)->bup.ovals[0] = now.read_count;
		now.read_count = (now.read_count+1);
#ifdef VAR_RANGES
		logval("read_count", now.read_count);
#endif
		;
		/* merge: i = (i+1)(18, 6, 18) */
		reached[0][6] = 1;
		(trpt+1)->bup.ovals[1] = ((P0 *)_this)->i;
		((P0 *)_this)->i = (((P0 *)_this)->i+1);
#ifdef VAR_RANGES
		logval("Reader:i", ((P0 *)_this)->i);
#endif
		;
		/* merge: .(goto)(0, 19, 18) */
		reached[0][19] = 1;
		;
		_m = 3; goto P999; /* 2 */
	case 41: // STATE 8 - modelo.pml:37 - [i = 0] (0:0:1 - 1)
		IfNotBlocked
		reached[0][8] = 1;
		(trpt+1)->bup.oval = ((P0 *)_this)->i;
		((P0 *)_this)->i = 0;
#ifdef VAR_RANGES
		logval("Reader:i", ((P0 *)_this)->i);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 42: // STATE 9 - modelo.pml:40 - [((i<2))] (0:0:0 - 1)
		IfNotBlocked
		reached[0][9] = 1;
		if (!((((P0 *)_this)->i<2)))
			continue;
		_m = 3; goto P999; /* 0 */
	case 43: // STATE 10 - modelo.pml:41 - [ch_raw!END] (0:0:0 - 1)
		IfNotBlocked
		reached[0][10] = 1;
		if (q_full(now.ch_raw))
			continue;
#ifdef HAS_CODE
		if (readtrail && gui) {
			char simtmp[64];
			sprintf(simvals, "%d!", now.ch_raw);
		sprintf(simtmp, "%d", 1); strcat(simvals, simtmp);		}
#endif
		
		qsend(now.ch_raw, 0, 1, 1);
		_m = 2; goto P999; /* 0 */
	case 44: // STATE 11 - modelo.pml:42 - [i = (i+1)] (0:0:1 - 1)
		IfNotBlocked
		reached[0][11] = 1;
		(trpt+1)->bup.oval = ((P0 *)_this)->i;
		((P0 *)_this)->i = (((P0 *)_this)->i+1);
#ifdef VAR_RANGES
		logval("Reader:i", ((P0 *)_this)->i);
#endif
		;
		_m = 3; goto P999; /* 0 */
	case 45: // STATE 21 - modelo.pml:49 - [-end-] (0:0:0 - 5)
		IfNotBlocked
		reached[0][21] = 1;
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

