package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
)

const BASE_URL string = "https://ciphersuite.info/cs/%s"

func main() {

	inputFileName := flag.String("f", "ciphers.txt", "Text file containing the list of ciphers suites to check.")

	outputFileName := flag.String("o", "output-ciphersuites.csv", "Output file containing the list of ciphers suites results.")

	reader, err := os.Open(*inputFileName)
	if err != nil {
		fmt.Printf("[!] Cannot read file %q: %s\n", *inputFileName, err)
		return
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)

	file, err := os.Create(*outputFileName)

	if err != nil {
		fmt.Printf("[!] Cannot create file %q: %s\n", *outputFileName, err)
		return
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Cipher Suite", "Strength"})

	c := colly.NewCollector()

	c.OnHTML(".mb-4", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText(".break-all"),
			e.ChildText(".badge"),
		})
	})

	for scanner.Scan() {
		cipher := scanner.Text()
		cipherUrl := fmt.Sprintf(BASE_URL, cipher)
		c.Visit(cipherUrl)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[!] Error reading input file %q: %s, ", *inputFileName, err)
	}

	fmt.Printf("[+] Check finished, check file %q for results\n", *outputFileName)

}
