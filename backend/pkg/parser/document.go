package parser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"
)

// ParsePDF extracts text from a PDF file
func ParsePDF(reader io.ReaderAt, size int64) (string, error) {
	pdfReader, err := pdf.NewReader(reader, size)
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var text strings.Builder
	numPages := pdfReader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		pageText, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		text.WriteString(pageText)
		text.WriteString("\n")
	}

	result := text.String()
	if len(result) == 0 {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return strings.TrimSpace(result), nil
}

// ParseDOCX extracts text from a DOCX file
func ParseDOCX(reader io.Reader) (string, error) {
	// Read all data into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read data: %w", err)
	}

	// Create ReaderAt from data
	readerAt := bytes.NewReader(data)

	doc, err := docx.ReadDocxFromMemory(readerAt, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX: %w", err)
	}

	var text strings.Builder
	for _, item := range doc.Editable().GetContent() {
		text.WriteString(string(item))
		text.WriteString("\n")
	}

	result := text.String()
	if len(result) == 0 {
		return "", fmt.Errorf("no text content found in DOCX")
	}

	return strings.TrimSpace(result), nil
}

// ParsePlainText extracts text from plain text files
func ParsePlainText(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %w", err)
	}

	text := string(data)
	if len(text) == 0 {
		return "", fmt.Errorf("file is empty")
	}

	return strings.TrimSpace(text), nil
}

// ParseExcel extracts text from Excel files (xlsx, xls)
func ParseExcel(reader io.Reader) (string, error) {
	// Read all data into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read Excel file: %w", err)
	}

	// Open Excel file from memory
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to parse Excel file: %w", err)
	}
	defer f.Close()

	var text strings.Builder

	// Get list of all sheets
	sheets := f.GetSheetList()

	for _, sheetName := range sheets {
		// Add sheet name as header
		text.WriteString(fmt.Sprintf("\n=== Sheet: %s ===\n", sheetName))

		// Get all rows in sheet
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}

		// Process each row
		for _, row := range rows {
			if len(row) == 0 {
				continue
			}

			// Join cells with tab separator
			rowText := strings.Join(row, "\t")
			text.WriteString(rowText)
			text.WriteString("\n")
		}
	}

	result := text.String()
	if len(result) == 0 {
		return "", fmt.Errorf("no content found in Excel file")
	}

	return strings.TrimSpace(result), nil
}

// ParseCSV extracts text from CSV files
func ParseCSV(reader io.Reader) (string, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := csvReader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return "", fmt.Errorf("CSV file is empty")
	}

	var text strings.Builder

	for _, record := range records {
		if len(record) == 0 {
			continue
		}

		// Join fields with tab separator
		rowText := strings.Join(record, "\t")
		text.WriteString(rowText)
		text.WriteString("\n")
	}

	result := text.String()
	if len(result) == 0 {
		return "", fmt.Errorf("no content found in CSV")
	}

	return strings.TrimSpace(result), nil
}

// DocumentMetadata contains metadata extracted from documents
type DocumentMetadata struct {
	PageCount  int
	WordCount  int
	FileSize   int64
	SheetCount int // For Excel files
	FileType   string
}

// ExtractMetadata extracts metadata from document content
func ExtractMetadata(content string, fileType string, fileSize int64) DocumentMetadata {
	metadata := DocumentMetadata{
		FileType: fileType,
		FileSize: fileSize,
	}

	// Count words
	metadata.WordCount = len(strings.Fields(content))

	// Estimate page count (assuming ~500 words per page)
	if metadata.WordCount > 0 {
		metadata.PageCount = (metadata.WordCount / 500) + 1
	}

	// For Excel, count sheets (simple estimation from sheet headers)
	if fileType == "xlsx" || fileType == "xls" {
		metadata.SheetCount = strings.Count(content, "=== Sheet:")
	}

	return metadata
}
