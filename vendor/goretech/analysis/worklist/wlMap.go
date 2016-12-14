package worklist

import (
	"goretech/analysis/ssabuilder"
	"strconv"
)

/* TODO:
 * WlList sollte besser ein Interface definieren/implementieren.
 * Dann können wir hinterher evtl. besser auf go-Routinen umstellen. [stolz]
 */

// WlList represents a simple implementation of a worklist.
// An element can only be added once and the list is a FiFo list.
type WlList struct {
	*wlList
}

type wlList struct {
	wlMap map[*ContextCallSite]bool
	order []*ContextCallSite
}

var maxElemsIn int

// NewWlList returns a new worklist
func NewWlList() *WlList {
	m := make(map[*ContextCallSite]bool)
	order := new([]*ContextCallSite)
	l := &WlList{&wlList{m, *order}}
	return l
}

// GetFirstCCS returns the first contextCallsite.
func (l *WlList) getFirstCCS() *ContextCallSite {
	// Iterate through order until the element is in the map
	for i := 0; i < len(l.order); i++ {
		if ok := l.wlMap[l.order[i]]; ok {
			return l.order[i]
		}
	}
	return nil
}

func (l *WlList) String() string {
	s := "Wl: "
	for i := 0; i < len(l.order); i++ {
		if ok := l.wlMap[l.order[i]]; ok {
			s += "[ " + strconv.Itoa(l.order[i].GetID()) + " : " + l.order[i].Node().String() + "]"
		}
	}
	return s
}

// RemoveFirstCCS removes the first ContextCallSite from the list
func (l *WlList) removeFirstCCS() *ContextCallSite {
	var c (*ContextCallSite) = nil
	for i := 0; i < len(l.order); i++ {
		if ok := l.wlMap[l.order[i]]; ok {
			c = l.order[i]
			delete(l.wlMap, l.order[i])
			l.order = l.order[i:]
			return c
		}
	}
	return c
}

// RemoveFirst returns the first contextcallsite of the worklit and removes it from the list.
func (l *WlList) RemoveFirst() *ContextCallSite {
	return l.removeFirstCCS()
}

// Empty returns true if no element is remaining in the worklist.
func (l *WlList) Empty() bool {
	return len(l.wlMap) == 0
}

// Add adds a new contextCallsite to the list
// Does not update the position of a contextcallsite which is already in the list,
func (l *WlList) Add(c *ContextCallSite) {
	// do nothing if c is already in the list and return
	// increment the highest number of l and add c to lo
	_, ok := l.wlMap[c]
	if !ok {
		l.order = append(l.order, c)
		l.wlMap[c] = true
	}
	if maxElemsIn < len(l.wlMap) {
		maxElemsIn = len(l.wlMap)
	}
}

// Len returns the lenogth of the worklist
func (l *WlList) Len() int {
	return len(l.wlMap)
}

// AddSucc adds all sccessors of n to the worklist
func (l *WlList) AddSucc(n *ContextCallSite) {
	succs := getSuccessors(n.Node())
	for _, s := range succs {
		c := getContext(n, s)
		if c == nil {
			c = NewContextCallSite(n.Context(), s)
			c.SetIn(n.GetOut())
		}
		worklist.Add(c)
	}
	// For node send: Add all calls which uses the channel als successor
	send, ok := n.Node().(*ssabuilder.Send)
	if ok {
		for _, s := range send.GetCalls() {
			// should only one context, because every node should exists once within each context.
			c := getContextChannel(n, s)
			if c == nil {
				c = NewContextCallSite(n.Context(), s)
				// Overapproximate and set the out lattice of the sending node
				c.SetIn(n.GetOut())
			} else {
				newValue := n.GetOut().GetVal(send.Send.Chan)
				c.GetIn().SetVal(send.Send.Chan, newValue)
			}
			worklist.Add(c)
		}
	}
}
