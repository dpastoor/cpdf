/*
 * Merges PDF files
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dpastoor/goutils"
	"github.com/spf13/afero"
	"github.com/unidoc/unidoc/pdf"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Requires at least 3 arguments: output_path regex \n")
		fmt.Printf("Usage: go run pdf_merge.go output.pdf \n")
		os.Exit(1)
	}
	AppFs := afero.NewOsFs()
	dir := filepath.Dir(".")
	dirInfo, _ := afero.ReadDir(AppFs, dir)
	files := goutils.ListFiles(dirInfo)
	pdfFiles, err := goutils.ListMatchesByRegex(files, os.Args[2])
	if err != nil {
		fmt.Printf("error compiling regex: %s", err)
		os.Exit(1)
	}
	fmt.Println("merging files:")
	fmt.Println(pdfFiles)
	outputPath := os.Args[1]

	// Sanity check the input arguments.
	err = mergePdf(pdfFiles, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func mergePdf(inputPaths []string, outputPath string) error {
	pdfWriter := pdf.NewPdfWriter()

	for _, inputPath := range inputPaths {
		f, err := os.Open(inputPath)
		if err != nil {
			return err
		}

		defer f.Close()

		pdfReader, err := pdf.NewPdfReader(f)
		if err != nil {
			return err
		}

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
