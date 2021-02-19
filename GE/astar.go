package GE

import (
	"container/heap"
	"errors"
	"math"
)

var (
	ErrPathNotFound = errors.New("PathNotFound")
)

func DistanceFrom(p1, p2 [2]int) float64 {
	dx := math.Abs(float64(p1[0] - p2[0]))
	dy := math.Abs(float64(p1[1] - p2[1]))
	return math.Sqrt(dx*dx + dy*dy)
}

// function for finding shortest path. it returns nodes have passed in shortest path.
// Returning slice order is latest visited node. Thus, goal is first element in slice.
func FindPath(start, goal [2]int, NeighborsOfPos func([2]int) [][2]int) ([][2]int, error) {
	startinfo := newnodeinfo(start, 0, DistanceFrom(start, goal))
	all := make(map[[2]int]*nodeinfo)
	all[start] = startinfo

	var pq priorityQueue
	current := startinfo

	for current.n != goal {

		current.closed = true

		for _, n := range NeighborsOfPos(current.n) {
			gs := current.gscore + DistanceFrom(current.n, n)
			ninfo := all[n]
			if ninfo == nil {
				ninfo = newnodeinfo(n, gs, gs+DistanceFrom(n, goal))
				all[n] = ninfo
				heap.Push(&pq, ninfo)
			} else {
				if ninfo.closed || gs >= ninfo.gscore {
					continue
				}
				ninfo.gscore = gs
				ninfo.fscore = gs + DistanceFrom(n, goal)
				heap.Fix(&pq, ninfo.idx)
			}
			ninfo.camefrom = current
		}

		if len(pq) == 0 {
			return nil, ErrPathNotFound
		}
		current = pq[0] // to avoid typecast
		heap.Pop(&pq)
	}

	var rtn [][2]int
	ninfo := current
	for {
		rtn = append(rtn, ninfo.n)
		ninfo = ninfo.camefrom
		if ninfo == nil {
			break
		}
	}
	return rtn, nil
}

type nodeinfo struct {
	idx int // slice idx

	n        [2]int
	camefrom *nodeinfo
	gscore   float64
	fscore   float64
	closed   bool
}

func newnodeinfo(n [2]int, g, f float64) *nodeinfo {
	return &nodeinfo{
		n:      n,
		gscore: g,
		fscore: f,
	}
}

type priorityQueue []*nodeinfo

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].fscore < pq[j].fscore
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idx = i
	pq[j].idx = j
}

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq *priorityQueue) Push(x interface{}) {
	e := x.(*nodeinfo)
	*pq = append(*pq, e)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	e := old[n-1]
	*pq = old[:n-1]
	return e
}
