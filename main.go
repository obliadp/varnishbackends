package main

import (
	"fmt"
	"sort"
	"time"

	tm "github.com/buger/goterm"
	"github.com/phenomenes/vago"
)

var (
	backends          backendSlice
	PruneAfterSeconds = float64(60)
)

func main() {

	// Open the default Varnish Shared Memory file
	c := vago.Config{}
	v, err := vago.Open(&c)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse log lines forever in a go-routine. Upsert to backends struct and prune expired backends.
	go func() {
		for {
			v.Log("Backend_health", vago.RAW, vago.COPT_TAIL|vago.COPT_BATCH, func(vxid uint32, tag, _type, data string) int {
				l := parseLogLine(data)
				backends = backends.upsert(l)
				backends = backends.pruneKeys()
				return 0
			})
			v.Close()
		}
	}()

	for {
		// By moving cursor to top-left position we ensure that console output
		// will be overwritten each time, instead of adding new.
		tm.Clear()
		tm.MoveCursor(1, 1)

		sort.Sort(customSort{backends, func(x, y *logLine) bool {
			if x.Director != y.Director {
				return x.Director < y.Director
			}
			if x.Name != y.Name {
				return x.Name < y.Name
			}
			if x.Backend != y.Backend {
				return x.Backend < y.Backend
			}
			return false
		}})

		//		printRawJSON(backends)
		printTerse(backends)

		tm.Flush() // Call it every time at the end of rendering

		time.Sleep(time.Second)
	}
}
