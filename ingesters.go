// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package junit

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// IngestDir will search the given directory for XML files and return a slice
// of all contained JUnit test suite definitions.
func IngestDir(directory string) ([]Suite, error) {
	var filenames []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Add all regular files that end with ".xml"
		if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".xml") {
			filenames = append(filenames, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return IngestFiles(filenames)
}

// IngestFiles will parse the given XML files and return a slice of all
// contained JUnit test suite definitions.
func IngestFiles(filenames []string) ([]Suite, error) {
	var all []Suite

	for _, filename := range filenames {
		suites, err := ingestFile(filename)
		if err != nil {
			return nil, err
		}
		all = append(all, suites...)
	}

	return all, nil
}

func ingestFile(filename string) (s []Suite, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()
	return Ingest(f)
}

// Ingest will parse the given XML data and return a slice of all contained
// JUnit test suite definitions.
func Ingest(r io.Reader) ([]Suite, error) {
	var (
		suiteChan = make(chan Suite)
		suites    []Suite
	)

	nodes, err := parse(r)
	if err != nil {
		return nil, err
	}

	go func() {
		findSuites(nodes, suiteChan)
		close(suiteChan)
	}()

	for suite := range suiteChan {
		suites = append(suites, suite)
	}

	return suites, nil
}
