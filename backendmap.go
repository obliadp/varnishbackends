package main

import (
	"time"
)

type backendSlice []*logLine

type customSort struct {
	l    backendSlice
	less func(x, y *logLine) bool
}

// return length of backendSlice slice, to satifsy sort.Sort
func (x customSort) Len() int { return len(x.l) }

// swap satisfies sort.Sort interface
func (x customSort) Swap(i, j int) { x.l[i], x.l[j] = x.l[j], x.l[i] }

// compare two elements of backendSlice slice, to satifsy sort.Sort
func (x customSort) Less(i, j int) bool { return x.less(x.l[i], x.l[j]) }

// insert or update existing *logLine in backendSlice
func (b backendSlice) upsert(l *logLine) backendSlice {
	insert := true

	for k, v := range b {
		if v.Name == l.Name {
			// Same backend, replace
			b = append(b, l)
			copy(b[k:], b[k+1:]) // shift
			b[len(b)-1] = nil    // remove reference
			b = b[:len(b)-1]     // reslice
			insert = false
		}
	}
	if insert == true {
		b = append(b, l)
	}
	return b
}

// takes a backendSlice and returns a new one minus expired *logLines
func (b backendSlice) pruneKeys() backendSlice {

	var r backendSlice

	// create a new slice from all valid lines
	for _, v := range b {

		// Now - timestamp should not exceed PruneAfterSeconds
		delta := time.Now().Sub(v.Timestamp)
		if delta.Seconds() < PruneAfterSeconds {
			r = append(r, v)
		}
	}

	return r
}
