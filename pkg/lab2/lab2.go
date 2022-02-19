package lab2

import (
	"fmt"
	"math"
	"nks/configs"
	"nks/pkg/models"
	"strings"
)

type Lab2 struct {
	configs.Lab2Config
	paths        [][]int //indexes of vertexes
	fullTable    [][]bool
	workingTable []models.TableAndProbSystem
	pSystem      float64
}

func NewLab2(c configs.Lab2Config) Lab2 {
	return Lab2{Lab2Config: c}
}

func (l *Lab2) Start() {
	//1. Визначимо усі можливі шляхи, якими можна пройти від початку до кінця
	//схеми системи, за допомогою алгоритму пошуку у глибину
	l.findPaths()
	//2. Визначимо усі можливі стани системи, коли вона знаходиться у працездатному стані.
	l.workingState()
	l.countPSystem()
}

func (l *Lab2) findPaths() {
	for i, v := range l.Vertexes {
		if v.HasStart {
			l.findAllPathsFromVertex(i, l.Vertexes[i], []int{})
		}
	}
}

func (l *Lab2) findAllPathsFromVertex(vertexIdx int, vertex configs.Vertex, currPath []int) {
	currPath = append(currPath, vertexIdx)
	if vertex.HasExit {
		l.paths = append(l.paths, currPath)
		return
	}

	for i, v := range vertex.Edges {
		if v {
			newCurrPath := make([]int, len(currPath))
			copy(newCurrPath, currPath)
			l.findAllPathsFromVertex(i, l.Vertexes[i], newCurrPath)
		}
	}
}

func (l *Lab2) workingState() {
	var (
		lines = int(math.Pow(2, float64(len(l.Vertexes))))
	)

	l.initWorkingVariants(lines)
	l.prefillTableAndSetOnlyWorking(lines)
}

func (l *Lab2) initWorkingVariants(lines int) {
	l.fullTable = make([][]bool, lines)
	rows := make([]bool, len(l.Vertexes)*lines)
	for i := range l.fullTable {
		l.fullTable[i] = rows[i*len(l.Vertexes) : (i+1)*len(l.Vertexes)]
	}
}

func (l *Lab2) prefillTableAndSetOnlyWorking(lines int) {
	aligningLine := fmt.Sprintf("%%0%db", len(l.Vertexes))

	for i := 0; i < lines; i++ {
		for j, v := range fmt.Sprintf(aligningLine, i) {
			if v == '0' {
				l.fullTable[i][j] = false
			} else {
				l.fullTable[i][j] = true
			}
		}
		if l.comparePathPatterns(l.fullTable[i]) {
			l.workingTable = append(l.workingTable, models.TableAndProbSystem{
				Table: l.fullTable[i],
				Prob:  l.countProbInLine(l.fullTable[i]),
			})
		}
	}
}

func (l Lab2) comparePathPatterns(line []bool) bool {
OuterLoop:
	for _, path := range l.paths {
	InnerLoop:
		for _, vertexNum := range path {
			if line[vertexNum] {
				continue InnerLoop
			}
			continue OuterLoop
		}
		return true
	}
	return false
}

func (l Lab2) countProbInLine(line []bool) (prob float64) {
	if len(line) == 0 {
		return
	}
	prob = 1
	for i, v := range line {
		if v {
			prob *= l.Vertexes[i].UninterruptedProb
		} else {
			prob *= (1 - l.Vertexes[i].UninterruptedProb)
		}
	}
	return
}

func (l *Lab2) countPSystem() {
	for _, v := range l.workingTable {
		l.pSystem += v.Prob
	}
}

func (l Lab2) String() string {
	var b strings.Builder

	b.WriteString("1. Визначимо усі можливі шляхи, якими можна пройти від початку до кінця схеми системи, за допомогою алгоритму пошуку у глибину\n")
	b.WriteString(fmt.Sprintf("Всього шляхів: %d\n", len(l.paths)))
	for _, path := range l.paths {
		var pathString []string

		for _, v := range path {
			pathString = append(pathString, fmt.Sprintf("E%d", v+1))
		}
		b.WriteString(strings.Join(pathString, "->"))
		b.WriteString("\n")
	}

	b.WriteString("2. Визначимо усі можливі стани системи, коли вона знаходиться у працездатному стані.\n")
	b.WriteString(fmt.Sprintf("Всього варіантів: %d\n", len(l.workingTable)))

	var vertexes []string
	b.WriteString("|")
	for i := range l.Vertexes {
		vertexes = append(vertexes, fmt.Sprintf("% 4s% 2s", fmt.Sprintf("E%d", i+1), ""))
	}
	b.WriteString(strings.Join(vertexes, "|"))
	b.WriteString(fmt.Sprintf("|% 8s\n", "Pstate"))

	for _, line := range l.workingTable {
		var lines []string
		b.WriteString("|")
		for _, v := range line.Table {
			if v { //it will beautifully work only for <10
				lines = append(lines, fmt.Sprintf("% 4s% 2s", "+", ""))
			} else {
				lines = append(lines, fmt.Sprintf("% 4s% 2s", "-", ""))
			}
		}
		b.WriteString(strings.Join(lines, "|"))
		b.WriteString(fmt.Sprintf("|% 8.6f\n", line.Prob))
	}

	b.WriteString(fmt.Sprintf("3. Psystem = %.6f\n", l.pSystem))

	return b.String()
}
