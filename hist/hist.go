package hist

import (
	"fmt"
	"strings"
)

type IntHistogram struct {
	points    []int
	bins      []float64
	counts    []int
	maxcounts int
	min       int
	max       int
	nbins     int
}

func NewIntHistogram(nbins int) *IntHistogram {
	return &IntHistogram{
		make([]int, 0, 10000),
		make([]float64, nbins),
		make([]int, nbins),
		-1,
		-1,
		-1,
		nbins,
	}
}

func (ih *IntHistogram) AddPoint(p int) {
	if len(ih.points) == 0 || p < ih.min {
		ih.min = p
	}
	if len(ih.points) == 0 || p > ih.max {
		ih.max = p
	}
	ih.points = append(ih.points, p)
}

func (ih *IntHistogram) Draw(width int) string {
	ih.updateBins()

	var sb strings.Builder

	maxnamelen := 0
	for b := range ih.bins {
		l := len(fmt.Sprintf("%.2f", ih.bins[b]))
		if l > maxnamelen {
			maxnamelen = l
		}
	}

	for b := range ih.bins {
		name := fmt.Sprintf("%.2f", ih.bins[b])
		if len(name) < maxnamelen {
			for i := 0; i < maxnamelen-len(name); i++ {
				sb.WriteString(" ")
			}
		}
		sb.WriteString(fmt.Sprintf("%s ", name))
		w := float64(ih.counts[b]) * float64(width) / float64(ih.maxcounts)
		for i := 0; i < int(w); i++ {
			sb.WriteString("*")
		}
		sb.WriteString(fmt.Sprintf(" %d", ih.counts[b]))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (ih *IntHistogram) updateBins() {
	for i, p := range ih.points {
		bin := int(float64(ih.nbins-1) * float64(p-ih.min) / float64(ih.max-ih.min))
		ih.counts[bin]++
		if i == 0 || ih.counts[bin] > ih.maxcounts {
			ih.maxcounts = ih.counts[bin]
		}
	}
	for b := range ih.counts {
		left := float64(b)*float64(ih.max-ih.min)/float64(ih.nbins-1) + float64(ih.min)
		right := float64(b+1)*float64(ih.max-ih.min)/float64(ih.nbins-1) + float64(ih.min)
		binmid := ((left + right) / 2.0)
		ih.bins[b] = binmid
	}
}
