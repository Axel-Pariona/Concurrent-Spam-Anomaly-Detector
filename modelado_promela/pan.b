	switch (t->back) {
	default: Uerror("bad return move");
	case  0: goto R999; /* nothing to undo */

		 /* PROC :init: */

	case 3: // STATE 1
		;
		;
		delproc(0, now._nr_pr-1);
		;
		goto R999;

	case 4: // STATE 2
		;
		;
		delproc(0, now._nr_pr-1);
		;
		goto R999;

	case 5: // STATE 3
		;
		;
		delproc(0, now._nr_pr-1);
		;
		goto R999;

	case 6: // STATE 4
		;
		;
		delproc(0, now._nr_pr-1);
		;
		goto R999;

	case 7: // STATE 5
		;
		p_restor(II);
		;
		;
		goto R999;

		 /* PROC Deduplicator */

	case 8: // STATE 1
		;
		XX = 1;
		unrecv(now.ch_val_dedup, XX-1, 0, ((P3 *)_this)->msg, 1);
		((P3 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 9: // STATE 2
		;
	/* 0 */	((P3 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 10: // STATE 3
		;
		processed = trpt->bup.oval;
		;
		goto R999;

	case 11: // STATE 5
		;
	/* 0 */	((P3 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 12: // STATE 12
		;
		p_restor(II);
		;
		;
		goto R999;

		 /* PROC Validator */

	case 13: // STATE 1
		;
		XX = 1;
		unrecv(now.ch_norm_val, XX-1, 0, ((P2 *)_this)->msg, 1);
		((P2 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 14: // STATE 2
		;
	/* 0 */	((P2 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 15: // STATE 3
		;
		_m = unsend(now.ch_val_dedup);
		;
		goto R999;

	case 16: // STATE 4
		;
	/* 0 */	((P2 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 17: // STATE 5
		;
		_m = unsend(now.ch_val_dedup);
		;
		goto R999;

	case 18: // STATE 12
		;
		p_restor(II);
		;
		;
		goto R999;

		 /* PROC Normalizer */

	case 19: // STATE 1
		;
		XX = 1;
		unrecv(now.ch_reader_norm, XX-1, 0, ((P1 *)_this)->msg, 1);
		((P1 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 20: // STATE 2
		;
	/* 0 */	((P1 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 21: // STATE 3
		;
		_m = unsend(now.ch_norm_val);
		;
		goto R999;

	case 22: // STATE 4
		;
	/* 0 */	((P1 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 23: // STATE 5
		;
		_m = unsend(now.ch_norm_val);
		;
		goto R999;

	case 24: // STATE 12
		;
		p_restor(II);
		;
		;
		goto R999;

		 /* PROC Reader */
;
		;
		
	case 26: // STATE 2
		;
		_m = unsend(now.ch_reader_norm);
		;
		goto R999;

	case 27: // STATE 3
		;
		((P0 *)_this)->i = trpt->bup.oval;
		;
		goto R999;

	case 28: // STATE 5
		;
		_m = unsend(now.ch_reader_norm);
		;
		goto R999;

	case 29: // STATE 10
		;
		p_restor(II);
		;
		;
		goto R999;
	}

