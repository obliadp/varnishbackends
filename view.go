package main

import (
	"encoding/json"

	tm "github.com/buger/goterm"
)

func printRawJSON(backends backendSlice) {

	var outstring string

	for _, v := range backends {

		w, _ := json.Marshal(v)

		if v.State == "Went" {
			outstring = tm.Color(string(w), tm.YELLOW)
		} else {
			if v.Happy {
				outstring = tm.Color(string(w), tm.GREEN)
			} else {
				outstring = tm.Color(string(w), tm.RED)
			}
		}
		tm.Println(outstring)
	}
}

func printTerse(backends backendSlice) {

	var outstring, prev string

	for _, v := range backends {
		if prev != v.Director {
			tm.Printf("\n%25s: ", v.Director)
		}
		if v.State == "Went" {
			outstring = "[" + tm.Color(string(v.Name), tm.YELLOW) + "]\t"
		} else {
			if v.Happy {
				outstring = "[" + tm.Color(string(v.Name), tm.GREEN) + "]\t"
			} else {
				outstring = "[" + tm.Color(string(v.Name), tm.RED) + "]\t"
			}
		}
		// For some reason I don't know my X position here. Why?
		tm.Print(outstring)
		prev = v.Director
	}
}
