package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type logLine struct {
	VCL                 string
	Name                string
	Director            string
	Backend             string
	State               string
	Healthy             bool
	IPv4                bool
	IPv6                bool
	TransmitSuccess     bool
	ReadResponseSuccess bool
	Happy               bool
	GoodPolls           int
	Threshold           int
	Window              int
	ResponseTime        float64
	ExponentialAverage  float64
	HTTPResponse        string
	Failure             bool
	Timestamp           time.Time
}

// deconstruct logLines
// this is a bit messy and needs some cleanup
func parseLogLine(s string) *logLine {

	failure := true
	var vcl, name, backend, director string

	matches := regexp.MustCompile("^(\\S+)\\.(\\S+)\\s+(\\w+)\\s+(\\w+)\\s+(4|-)(6|-)(x|-)(X|-)(r|-)(R|-)(H|-)\\s(\\d)\\s(\\d)\\s(\\d)\\s(\\S+)\\s(\\S+)\\s(.*)$").FindStringSubmatch(s)

	// Check for goto-backends, and handle them in a separate function
	// https://docs.varnish-software.com/varnish-cache-plus/vmods/goto/
	if strings.Contains(s, ".goto.0") {
		vcl, name, backend, director = parseGotoLogLine(s)
		failure = false
	} else {
		backend, director = matches[2], "none"

		// Try to pick some info out of matches[2]
		re, err := regexp.CompilePOSIX("^(.+)_([^_]+)$")
		if err != nil {
			fmt.Println(err)
		}
		sm := re.FindStringSubmatch(matches[2])

		// if our regex is successful, dress up backend and director values
		if len(sm) > 1 {
			backend, director = sm[2], sm[1]
			failure = false
		}
	}

	// vcl might be set if backend is goto-type
	if vcl == "" {
		vcl = matches[1]
	}

	// name might be set if backend is goto-type
	if name == "" {
		name = matches[2]
	}

	// an 'X' and no 'x' indicates socket transmit success
	trxSuccess := false
	if matches[7] == "-" && matches[8] == "X" {
		trxSuccess = true
	}

	// an 'R' and no 'r' indicates socket read success
	readSuccess := false
	if matches[9] == "-" && matches[10] == "R" {
		readSuccess = true
	}

	healthy := matches[4] == "healthy"
	ipv4 := matches[5] == "4"
	ipv6 := matches[6] == "6"
	happy := matches[11] == "H"
	goodPolls, _ := strconv.Atoi(matches[12])
	threshold, _ := strconv.Atoi(matches[13])
	window, _ := strconv.Atoi(matches[14])
	responseTime, _ := strconv.ParseFloat(matches[15], 64)
	exponentialAverage, _ := strconv.ParseFloat(matches[16], 64)

	m := logLine{
		VCL:                 vcl,
		Name:                name,
		Backend:             backend,
		Director:            director,
		State:               matches[3],
		Healthy:             healthy,
		IPv4:                ipv4,
		IPv6:                ipv6,
		TransmitSuccess:     trxSuccess,
		ReadResponseSuccess: readSuccess,
		Happy:               happy,
		GoodPolls:           goodPolls,
		Threshold:           threshold,
		Window:              window,
		ResponseTime:        responseTime,
		ExponentialAverage:  exponentialAverage,
		HTTPResponse:        matches[17],
		Timestamp:           time.Now(),
		Failure:             failure,
	}

	return &m
}

// Goto backends have radically different names
func parseGotoLogLine(s string) (string, string, string, string) {

	matches := regexp.MustCompile("^(\\S+)\\.goto\\.(\\S+)\\.\\((\\S+)\\)\\.\\((\\S+)(\\s|\\))").FindStringSubmatch(s)

	vcl := matches[1]
	name := matches[2]
	backend := matches[3]
	director := matches[4]

	if strings.Contains(director, "http") {
		m := regexp.MustCompile("^(?:https?)?(?:://)?(\\S+)").FindStringSubmatch(director)
		director = m[1]
	}

	if len(director) > 25 {
		director = director[:25]
	}

	return vcl, name, backend, director
}
