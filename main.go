package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var layoutPath = flag.String("template", "layout/layout.html", "HTML template used to generate pdf document")
var pdfPath = flag.String("out", "pdf", "Output folder")
var csvFilePath = flag.String("in", "input.csv", "Input CSV file path")

const fileIDFieldName = "FirstName"

func main() {
	flag.Parse()

	headers, data := getCsvData(*csvFilePath)

	fmt.Printf("Total number of records read : %d\n", len(data))

	dataMap := make(map[string]string)

	tempDir, err := ioutil.TempDir("", "go_receipt*")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to create a temporary folder: %v", err))
	}
	for _, record := range data {
		for i, field := range record {
			dataMap[headers[i]] = field
		}
		fileID := dataMap[fileIDFieldName]
		filename := generateHTML(dataMap, tempDir, fileID)
		generatePdf(filename, fileID)
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to delete temporary folder: %v", err))
	}
}

func getCsvData(filename string) ([]string, [][]string) {
	var data [][]string
	var headers []string

	csvfile, err := os.Open(filename)
	defer csvfile.Close()

	if err != nil {
		return headers, data
	}

	r := csv.NewReader(csvfile)
	r.TrimLeadingSpace = true

	data, _ = r.ReadAll()

	if len(data) > 0 {
		headers = data[0]
	}
	return headers, data[1:]
}

func generateHTML(data map[string]string, tempDir, fileID string) string {
	tmpl, err := template.ParseFiles(*layoutPath)
	if err != nil {
		return ""
	}

	filename := tempDir + "/" + fileID + ".html"

	outputFile, err := os.Create(filename)
	defer outputFile.Close()

	if err != nil {
		return ""
	}
	tmpl.Execute(outputFile, data)

	return filename
}

func generatePdf(htmlFileName string, fileID string) {
	outfile := fmt.Sprintf("%s/%s.pdf", *pdfPath, fileID)
	cmd := exec.Command("wkhtmltopdf", htmlFileName, outfile)
	cmd.Run()
}
