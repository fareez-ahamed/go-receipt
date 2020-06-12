package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"os"
	"os/exec"
)

const layoutPath = "layout/layout.html"
const tmpPath = "tmp"
const pdfPath = "pdf"
const fileIDFieldName = "FirstName"

func main() {
	filename := getCsvFilename()
	if filename == "" {
		return
	}

	headers, data := getCsvData(filename)

	fmt.Printf("Total number of records read : %d", len(data))

	dataMap := make(map[string]string)

	for _, record := range data {
		for i, field := range record {
			dataMap[headers[i]] = field
		}
		fileID := dataMap[fileIDFieldName]
		filename := generateHTML(dataMap, fileID)
		generatePdf(filename, fileID)
	}
}

func getCsvFilename() string {
	fmt.Println(os.Args)
	if len(os.Args) < 2 {
		fmt.Println("Please provide input file in command line argument")
		return ""
	}
	return os.Args[1]
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

func generateHTML(data map[string]string, fileID string) string {
	tmpl, err := template.ParseFiles(layoutPath)
	if err != nil {
		return ""
	}

	filename := tmpPath + fileID + ".html"

	outputFile, err := os.Create(filename)
	if err != nil {
		return ""
	}
	tmpl.Execute(outputFile, data)

	return filename
}

func generatePdf(htmlFileName string, fileID string) {
	outfile := fmt.Sprintf("%s/%s.pdf", pdfPath, fileID)
	cmd := exec.Command("wkhtmltopdf", htmlFileName, outfile)
	cmd.Run()
}
