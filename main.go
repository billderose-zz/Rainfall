package main

import (
	"fmt"
	"math"
	"runtime"
	"sort"
)

type Cell struct {
	row, col int
	value    float64
}

type Basin struct {
	sinks []Cell
}

type Basins []Basin

type Rainfall [][]Cell

// Wrapper for Math.max to avoid type conversions
func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

// Wrapper for Math.min to avoid type conversions
func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

// neighbors returns a slice of all cells adjacent to c
// in Rainfall r
func (c Cell) neighbors(r *Rainfall) []Cell {
	dim := len(*r)
	neighbors := make([]Cell, 0)
	for i := max(0, c.row-1); i <= min(c.row+1, dim-1); i++ {
		for j := max(0, c.col-1); j <= min(c.col+1, dim-1); j++ {
			if i != c.row || j != c.col {
				neighbors = append(neighbors, (*r)[i][j])
			}
		}
	}
	return neighbors
}

// isSink determines whether or not Cell c is a sink
func (rainfall Rainfall) isSink(c *Cell) bool {
	isSink := true
	dim := len(rainfall)
	for i := max(0, c.row-1); i <= min(c.row+1, dim-1); i++ {
		for j := max(0, c.col-1); j <= min(c.col+1, dim-1); j++ {
			isSink = isSink && c.value <= rainfall[i][j].value
		}
	}
	return isSink
}

// findSinks returns a channel over which the sinks
// of rainfall are sent over
func (rainfall Rainfall) findSinks() chan Cell {
	result := make(chan Cell, 0)
	go func() {
		for i := range rainfall {
			for j := range rainfall {
				if rainfall.isSink(&rainfall[i][j]) {
					result <- rainfall[i][j]
				}
			}
		}
		close(result)
	}()
	return result
}


// findBasins returns an object of type Basins (a slice of Basin)
// containing all the basins in rainfall
func (rainfall Rainfall) findBasins() Basins {
	sinks := rainfall.findSinks()
	basins := make(chan Basin, 10)
	size := 0
	for sink := range sinks {
		size++
		go func(s Cell) {
			b := Basin{sinks: make([]Cell, 0)}
			b.sinks = append(b.sinks, s)
			ns := s.neighbors(&rainfall)
			for _, neighbor := range ns {
				b.expand(s, neighbor, rainfall)
			}
			basins <- b
		}(sink)
	}
	resultBasins := make(Basins, 0)
	for i := 0; i < size; i++ {
		resultBasins = append(resultBasins, <-basins)
	}
	close(basins)
	return resultBasins
}

// expand recursively expands Basin b as far as possible
func (b *Basin) expand(sink, cand Cell, rainfall Rainfall) {
	neighbors := cand.neighbors(&rainfall)
	for _, neighbor := range neighbors {
		if neighbor.value < sink.value {
			// neighbor drains to another sink
			return
		}
	}
	b.sinks = append(b.sinks, cand)
	for _, neighbor := range neighbors {
		if neighbor.value > cand.value {
			b.expand(cand, neighbor, rainfall)
		}
	}
}

// Len required to sort Basins
func (b Basins) Len() int {
	return len(b)
}

// Swap required to sort Basins
func (b Basins) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Less required to sort Basins
func (b Basins) Less(i, j int) bool {
	return len(b[i].sinks) < len(b[j].sinks)
}

func main() {
	runtime.GOMAXPROCS(1)
	var (
		n int
		s float64
	)
	if _, err := fmt.Scan(&n); err != nil { // get dimension
		panic(err)
	}
	farm := make(Rainfall, n)
	for i := range farm { // read input
		farm[i] = make([]Cell, n)
		for j := range farm {
			if _, err := fmt.Scanf("%f", &s); err != nil {
				panic(err)
			}
			farm[i][j] = Cell{i, j, s}
		}
	}

	basins := farm.findBasins()
	sort.Sort(basins)
	sort.Sort(sort.Reverse(basins))
	for _, b := range basins {
		fmt.Printf("%d ", len(b.sinks))
	}
}
