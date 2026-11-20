package ui

const histLen = 120

func AppendHist(h []float64, v float64) []float64 {
	h = append(h, v)
	if len(h) > histLen {
		h = h[len(h)-histLen:]
	}
	return h
}

func LastOf(h []float64) float64 {
	if len(h) == 0 {
		return 0
	}
	return h[len(h)-1]
}
