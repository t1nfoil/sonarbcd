package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// written by ChatGPT. :)
func zipUpLabels(outputDirectory, zipName string) error {
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

		if !strings.HasSuffix(path, ".svg") {
			return nil
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
