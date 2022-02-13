package models

//[minIntervalValue, maxIntervalValue, amountOfValues, densityOnInterval, uninterruptedProbability]
type Interval [5]float64

const (
	MinIdx = iota
	MaxIdx
	ValuesAmountIdx
	DensityIdx
	UninterruptedProbIdx
)

const (
	P0 = 1
)

func CountIntervalsDensitySumBeforeIndex(idx int, intervals []Interval) (sum float64) {
	for _, v := range intervals[:idx] {
		sum += v[DensityIdx]
	}
	return
}

func CountDeltha(gamma float64, idx int, intervals []Interval) float64 {
	var (
		probFirst float64
		probLast  float64
	)
	if idx == 0 {
		probFirst = P0
	} else {
		probFirst = intervals[idx-1][UninterruptedProbIdx]
	}
	probLast = intervals[idx][UninterruptedProbIdx]
	return (probLast - gamma) / (probLast - probFirst)
}
