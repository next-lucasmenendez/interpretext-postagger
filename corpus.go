// Load and parse local corpus to train HMM model.
package gopostagger

import (
	"os"
	"fmt"
	"errors"
	"regexp"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const CorpusDir = "./corpus"

// Auxiliar struct that contains tagged word data.
type Item struct {
	Order int
	Raw, Tag string
}

// Line struct parsed items list
type Record []Item

func (r Record) Len() int {
	return len(r)
}

func (r Record) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Record) Less(i, j int) bool {
	return r[i].Order < r[j].Order
}

// Auxiliar struct that contains corpus structured data.
type Corpus struct {
	Name string
	Tags []string	
	Records []Record
}

// Return a structured and parsed corpus based on a name.
func GetCorpus(name string) (Corpus, error) {
	var corpus Corpus = Corpus{Name: name}

	var err error
	var exists bool
	if exists, err = corpus.exists(); err != nil {
		return corpus, err
	} else if !exists {
		return corpus, errors.New("Corpus doesn't exists.")
	}

	if corpus.Records, err = corpus.getRecords(); err != nil {
		return corpus, err
	}
	corpus.Tags = corpus.parseTags()
	return corpus, nil
}

// Check if corpus exists on based path.
func (corpus Corpus) exists() (bool, error) {
	var err error
	var location string

	if location, err = filepath.Abs(CorpusDir); err != nil {
		return false, err
	}

	var avalible []os.FileInfo
	if avalible, err = ioutil.ReadDir(location); err != nil {
		return false, err
	}

	for _, c := range avalible {
		if c.IsDir() && c.Name() == corpus.Name {
			return true, nil
		}
	}

	return false, nil
}

// Call auxiliar functions to get all records parsed from corpus raw information.
func (corpus Corpus) getRecords() ([]Record, error) {
	var err error
	var lines []string
	if lines, err = corpus.getLines(); err != nil {
		return nil, err
	}

	var records []Record = parseLines(lines)
	return records, nil
}

// Read file locally and split in lines cleaning empty lines.
func (corpus Corpus) getLines() ([]string, error) {
	var err error
	var location string = fmt.Sprintf("%s/%s", CorpusDir, corpus.Name)

	var files []os.FileInfo
	if files, err = ioutil.ReadDir(location); err != nil {
		return nil, err
	}
	
	var lines []string
	for _, file := range files {
		if !file.IsDir() {
			var filename string = fmt.Sprintf("%s/%s", location, file.Name())

			var raw []byte
			if raw, err = ioutil.ReadFile(filename); err != nil {
				return nil, err
			}

			var rgxLinebreak *regexp.Regexp = regexp.MustCompile(`\n`)
			var parsedLines []string = rgxLinebreak.Split(string(raw), -1)
			lines = append(lines, parsedLines...)
		}
	}
	return lines, nil
}

// Split list of lines into words extracting each tags from raw corpus, and generate Item structs with their data.
func parseLines(lines []string) []Record {
	var wordRgx *regexp.Regexp = regexp.MustCompile(`\s`)
	var tagRgx *regexp.Regexp = regexp.MustCompile(`(.+)\/(.+)`)

	var records []Record
	for _, line := range lines {
		var record Record
		var words []string = wordRgx.Split(line, -1)

		for order, raw_word := range words {
			var groups []string = tagRgx.FindStringSubmatch(raw_word)
			if len(groups) > 0 {
				var word, tag string = groups[1], groups[2]
				var item Item = Item{Raw: strings.ToLower(word), Tag: tag, Order: order}

				record = append(record, item)
			}
		}
		if len(record) > 0 {
			records = append(records, record)
		}
	}

	return records
}

// GetCorpus all posible tags of a corpus. Iterate over all corpus to prevent repeated tags.
func (corpus Corpus) parseTags() []string {
	var tags []string
	for _, record := range corpus.Records {
		for _, item := range record {
			var included bool = false
			for _, tag := range tags {
				if tag == item.Tag {
					included = true
					break
				}
			}
			if !included {
				tags = append(tags, item.Tag)
			}
		}
	}
	return tags
}
