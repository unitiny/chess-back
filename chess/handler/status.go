package handler

import "Chess/chess/lib"

type ChessStatus struct {
	lib.BasicStatus
}

func hasSameStatus(s1, s2 *ChessStatus, pos int) bool {
	offset := 1 << pos
	return s1.Has(offset) && s2.Has(offset)
}

func isSameCamp(s1, s2 *ChessStatus) bool {
	return s1.Has(1) == s2.Has(1)
}

func isSameCamp2(s *ChessStatus, camp int) bool {
	return s.IsSame(camp, 0)
}

func isAlive(s *ChessStatus) bool {
	return s.Has(ALIVE)
}
