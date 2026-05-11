	switch (t->back) {
	default: Uerror("bad return move");
	case  0: goto R999; /* nothing to undo */

		 /* CLAIM NO_NEGATIVOS */
;
		
	case 3: // STATE 1
		goto R999;

	case 4: // STATE 10
		;
		p_restor(II);
		;
		;
		goto R999;

		 /* PROC :init: */

	case 5: // STATE 1
		;
		((P3 *)_this)->i = trpt->bup.oval;
		;
		goto R999;

	case 6: // STATE 2
		;
		;
		delproc(0, now._nr_pr-1);
		;
		goto R999;
;
		;
		
	case 8: // STATE 4
		;
		;
		delproc(0, now._nr_pr-1);
		;
		goto R999;

	case 9: // STATE 5
		;
		((P3 *)_this)->i = trpt->bup.oval;
		;
		goto R999;

	case 10: // STATE 11
		;
		((P3 *)_this)->i = trpt->bup.oval;
		;
		goto R999;
;
		;
		
	case 12: // STATE 13
		;
		;
		delproc(0, now._nr_pr-1);
		;
		goto R999;

	case 13: // STATE 14
		;
		((P3 *)_this)->i = trpt->bup.oval;
		;
		goto R999;

	case 14: // STATE 20
		;
		p_restor(II);
		;
		;
		goto R999;

		 /* PROC Validator */

	case 15: // STATE 1
		;
		XX = 1;
		unrecv(now.ch_clean, XX-1, 0, ((P2 *)_this)->msg, 1);
		((P2 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 16: // STATE 2
		;
	/* 0 */	((P2 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 17: // STATE 3
		;
		now.validated_count = trpt->bup.oval;
		;
		goto R999;

	case 18: // STATE 5
		;
	/* 0 */	((P2 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 19: // STATE 6
		;
		now.val_done = trpt->bup.oval;
		;
		goto R999;

	case 20: // STATE 14
		;
		p_restor(II);
		;
		;
		goto R999;

		 /* PROC Normalizer */

	case 21: // STATE 1
		;
		((P1 *)_this)->c = trpt->bup.oval;
		;
		goto R999;

	case 22: // STATE 2
		;
		XX = 1;
		unrecv(now.ch_raw, XX-1, 0, ((P1 *)_this)->msg, 1);
		((P1 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 23: // STATE 3
		;
	/* 0 */	((P1 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;
;
		;
		
	case 25: // STATE 13
		;
		((P1 *)_this)->c = trpt->bup.ovals[1];
		now.rejected_count = trpt->bup.ovals[0];
		;
		ungrab_ints(trpt->bup.ovals, 2);
		goto R999;

	case 26: // STATE 8
		;
		now.normalized_count = trpt->bup.oval;
		;
		goto R999;

	case 27: // STATE 10
		;
		_m = unsend(now.ch_clean);
		;
		goto R999;

	case 28: // STATE 13
		;
		((P1 *)_this)->c = trpt->bup.oval;
		;
		goto R999;

	case 29: // STATE 14
		;
	/* 0 */	((P1 *)_this)->msg = trpt->bup.oval;
		;
		;
		goto R999;

	case 30: // STATE 15
		;
		now.norm_done = trpt->bup.oval;
		;
		goto R999;
;
		;
		
	case 32: // STATE 18
		;
		((P1 *)_this)->c = trpt->bup.oval;
		;
		goto R999;
;
		;
		
	case 34: // STATE 20
		;
		_m = unsend(now.ch_clean);
		;
		goto R999;

	case 35: // STATE 21
		;
		((P1 *)_this)->c = trpt->bup.oval;
		;
		goto R999;

	case 36: // STATE 37
		;
		p_restor(II);
		;
		;
		goto R999;

		 /* PROC Reader */

	case 37: // STATE 1
		;
		((P0 *)_this)->i = trpt->bup.oval;
		;
		goto R999;
;
		;
		
	case 39: // STATE 3
		;
		_m = unsend(now.ch_raw);
		;
		goto R999;

	case 40: // STATE 6
		;
		((P0 *)_this)->i = trpt->bup.ovals[1];
		now.read_count = trpt->bup.ovals[0];
		;
		ungrab_ints(trpt->bup.ovals, 2);
		goto R999;

	case 41: // STATE 8
		;
		((P0 *)_this)->i = trpt->bup.oval;
		;
		goto R999;
;
		;
		
	case 43: // STATE 10
		;
		_m = unsend(now.ch_raw);
		;
		goto R999;

	case 44: // STATE 11
		;
		((P0 *)_this)->i = trpt->bup.oval;
		;
		goto R999;

	case 45: // STATE 21
		;
		p_restor(II);
		;
		;
		goto R999;
	}

