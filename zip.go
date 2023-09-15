package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func zipUpLabels(outputDirectory, zipName string) error {
	// if zip doesn't end it .zip, add it, otherwise leave it alone
	if !strings.HasSuffix(zipName, ".zip") {
		zipName += ".zip"
	}

	zipFile, err := os.Create(filepath.Join(outputDirectory, zipName))
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(outputDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(outputDirectory, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		}

		entry, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(entry, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}
