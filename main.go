package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// Initialize a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord")

	// Process each line from standard input as a separate domain
	for scanner.Scan() {
		domain := scanner.Text()
		checkDomain(domain)
	}

	// Check for errors that might have occurred during the scan
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from input: %v", err)
	}
}

// checkDomain performs DNS lookups for the specified domain and outputs relevant DNS record information
func checkDomain(domain string) {
	// Print DNS record information for the domain
	printDNSRecords(domain)
}

// printDNSRecords queries DNS records for the given domain and prints them
func printDNSRecords(domain string) {
	hasMX, hasSPF, hasDMARC := false, false, false
	spfRecord, dmarcRecord := "", ""

	// Check for MX records (used for email exchange servers)
	if mxRecords, err := net.LookupMX(domain); err == nil && len(mxRecords) > 0 {
		hasMX = true // MX records found
	}

	// Check for TXT records (used for SPF and DMARC records)
	txtRecords, _ := net.LookupTXT(domain)
	for _, record := range txtRecords {
		switch {
		case strings.HasPrefix(record, "v=spf1"):
			hasSPF, spfRecord = true, record // SPF record found
		case strings.HasPrefix(record, "v=DMARC1"):
			hasDMARC, dmarcRecord = true, record // DMARC record found
		}
	}

	// Check specifically for DMARC TXT records at the _dmarc subdomain
	dmarcRecords, _ := net.LookupTXT("_dmarc." + domain)
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC, dmarcRecord = true, record // DMARC record found
		}
	}

	// Output the results in a CSV format
	fmt.Printf("%s, %t, %t, %s, %t, %s\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
