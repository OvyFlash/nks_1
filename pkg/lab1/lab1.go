package lab1

import (
	"fmt"
	"lab1/configs"
	"lab1/pkg/models"
	"sort"
	"strings"
)

type Lab1 struct {
	configs.Config `json:"-"`
	//result field
	Tcp       float64                         `json:"Tcp"`
	Ty        float64                         `json:"Ty"`
	Tt1       float64                         `json:"Tt1"`
	Lambda    float64                         `json:"lambda"`
	Intervals [intervalAmount]models.Interval `json:"-"`
	h         float64                         `json:"-"`
}

const (
	intervalAmount = 10
	smallNumber    = 1e-5
)

func NewLab1(c configs.Config) Lab1 {
	return Lab1{Config: c}
}

func (l *Lab1) Start() {
	// 1. Знайдемо середній наробіток до відмови Tср.
	l.findAverageOperatingTime()
	// 2. Відсортуємо вхідну вибірку наробітків до відмови
	sort.Slice(l.Data, func(i, j int) bool {
		return l.Data[i] < l.Data[j]
	})
	l.h = float64(l.Data[len(l.Data)-1]) / float64(intervalAmount)
	//3. Поділимо інтервал від 0 до максимального наробітку до відмови на 10 інтервалів
	l.splitByIntervals()
	//4. Для кожного інтервалу обрахуємо значення статистичної щільності розподілу ймовірності відмови
	l.calculateDensity()
	//5. Для кожного інтервалу обрахуємо значення ймовірності безвідмовної роботи
	//пристрою на час правої границі інтервалу
	l.calculateUninterruptedProbability()
	//6. Відсотковий наробіток на відмову
	l.calculateTy()
	//7. Ймовірність безвідмовної роботи на час time_1
	l.Tt1 = l.calculateTt(float64(l.Time1))
	//8. Інтенсивність відмов на час time_2
	l.calculateLambda()
}

func (l *Lab1) findAverageOperatingTime() {
	for _, v := range l.Data {
		l.Tcp += float64(v)
	}
	l.Tcp /= float64(len(l.Data))
}

func (l *Lab1) splitByIntervals() {
	for i := range l.Intervals {
		l.Intervals[i][models.MinIdx] = l.h * float64(i)
		l.Intervals[i][models.MaxIdx] = l.h * float64(i+1)
	}
}

func (l *Lab1) calculateDensity() {
	l.calculateAmountOfValuesInInterval()
	var (
		N = float64(len(l.Data))
	)
	for i := range l.Intervals {
		l.Intervals[i][models.DensityIdx] = l.Intervals[i][models.ValuesAmountIdx] / (N * l.h)
	}
}

func (l *Lab1) calculateAmountOfValuesInInterval() {
	var currIntervalIndex int
	for _, data := range l.Data {
		for float64(data) >= l.Intervals[currIntervalIndex][models.MaxIdx]+smallNumber && currIntervalIndex != len(l.Intervals)-1 {
			currIntervalIndex++
		}
		l.Intervals[currIntervalIndex][models.ValuesAmountIdx]++
	}
}

func (l *Lab1) calculateUninterruptedProbability() {
	for i := range l.Intervals {
		if i == len(l.Intervals)-1 {
			break
		}
		l.Intervals[i][models.UninterruptedProbIdx] = 1 - (l.h * models.CountIntervalsDensitySumBeforeIndex(i+1, l.Intervals[:]))
	}
}

func (l *Lab1) calculateTy() {
	for i, v := range l.Intervals {
		if l.Gamma > v[models.DensityIdx] {
			l.Ty = models.CountDeltha(l.Gamma, i, l.Intervals[:])
			return
		}
	}
}

func (l *Lab1) calculateTt(t float64) (tt float64){
	
	for _, v := range l.Intervals {
		if float64(t) >= v[models.MinIdx] {
			if v[models.MinIdx]+l.h < float64(t) {
				tt += v[models.DensityIdx] * l.h
				continue
			}
			tt += v[models.DensityIdx] * (float64(t) - v[models.MinIdx])
			break
		}
	}
	return 1 - tt
}

func (l *Lab1) calculateLambda() {
	var maxDensity float64
	for _, v := range l.Intervals {
		if float64(l.Time2) >= v[models.MinIdx] && v[models.MinIdx]+l.h > float64(l.Time2) {
			maxDensity = v[models.DensityIdx]
			break
		}
	}

	l.Lambda = maxDensity / l.calculateTt(float64(l.Time2))
}

func (l Lab1) String() string {
	var b strings.Builder
	if l.Tcp != 0 {
		b.WriteString(fmt.Sprintf("1. Середній наробіток до відмови Tcp: %.2f\n", l.Tcp))
	}
	b.WriteString(fmt.Sprintf("2. Відсортований масив: %v\n", l.Data))
	b.WriteString("3. Інтервали:\n")
	for i, v := range l.Intervals {
		b.WriteString(fmt.Sprintf("%3d-й інтервал від %8.2f  до %8.2f. Кількість елементів: %2.f\n",
			i+1,
			v[models.MinIdx],
			v[models.MaxIdx],
			v[models.ValuesAmountIdx]))
	}
	b.WriteString("4. Статистична щільність:\n")
	for i, v := range l.Intervals {
		b.WriteString(fmt.Sprintf("%3d-й інтервал від %8.2f  до %8.2f. f%2d = %.6f\n",
			i+1,
			v[models.MinIdx],
			v[models.MaxIdx],
			i+1,
			v[models.DensityIdx]))
	}
	b.WriteString("5. Ймовірність безвідмовної роботи:\n")
	for i, v := range l.Intervals {
		b.WriteString(fmt.Sprintf("%3d-й інтервал P(%7.2f) = %.6f\n",
			i+1,
			v[models.MaxIdx],
			v[models.UninterruptedProbIdx]))
	}
	b.WriteString(fmt.Sprintf("6. Статистичний відсотковий наробіток (Т_%.2f) = %.2f\n", l.Gamma, l.Ty))
	b.WriteString(fmt.Sprintf("7. Ймовірність безвідмовної роботи на час %d = %.5f\n", l.Time1, l.Tt1))
	b.WriteString(fmt.Sprintf("8. Інтенсивність відмов на час %d = %.5f\n", l.Time2, l.Lambda))

	return b.String()
}
