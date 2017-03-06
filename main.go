package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

const (
	// Change these constants if your contact base is not the USA
	countryCode               = "+1" // USA
	digitsInHomeCountryNumber = 10   // USA
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("First argument must be input CSV file from Google contacts.")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	w := csv.NewWriter(os.Stdout)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	// First header is the record
	isPhoneCol := make(map[int]bool, 0)
	for i := 0; ; i++ {
		line, err := r.Read()
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		if i == 0 { // Header line; only consider phone value fields
			for x := range line {
				field := strings.ToLower(line[x])
				if !strings.HasPrefix(field, "phone") ||
					!strings.Contains(field, "value") {
					continue
				}
				isPhoneCol[x] = true
			}
		} else { // Contact lines
			for x := range line {
				if !isPhoneCol[x] || line[x] == "" || strings.Contains(line[x], "+") {
					continue // Don't worry about blank, non-phone, or already country-coded lines
				}
				// Count digits
				var digits int
				for _, c := range line[x] {
					if c >= '0' && c <= '9' {
						digits++
					}
				}
				if digits != digitsInHomeCountryNumber {
					continue
				}
				line[x] = countryCode + " " + line[x] // Prefix the phone number as-is
			}
		}
		w.Write(line)
	}
	w.Flush()
}
