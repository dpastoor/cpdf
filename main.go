/*
 * Merges PDF files
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dpastoor/goutils"
	flag "github.com/ogier/pflag"
	"github.com/spf13/afero"
	"github.com/unidoc/unidoc/pdf"
)

// VERSION is the version
const VERSION = "1.0.0"

var (
	ver   bool
	debug bool
)

func main() {
	flag.BoolVarP(&ver, "version", "v", false, "print version")
	flag.BoolVarP(&debug, "debug", "d", false, "print debug information such as the files that will be combined")
	flag.Parse()
	if ver {
		fmt.Print(VERSION)
		os.Exit(0)
	}
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
	if debug {
		fmt.Println("merging files:")
		fmt.Println(pdfFiles)
	}
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
	outputFilename := filepath.Base(outputPath)
	for _, inputPath := range inputPaths {
		if inputPath == outputFilename {
			continue
		}
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
