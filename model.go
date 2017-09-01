// Contains all structs and functions to load, save and train a model based on a corpus.
package gopostagger

import (
	"os"
	"fmt"
	"sort"
	"bufio"
	"errors"
	"regexp"
	"strconv"
	"io/ioutil"
	"path/filepath"
)

const (
	ModelsPath string = "./models"
	Start string = "<s>"
	End string = "</s>"
)

// Auxiliar struct that contains tag or word relations with their weights.
type Link struct {
	current, previous string // word, tag (emission) - tag, tag (transition)
	occurrences, weight float64 
}

// List of links to sort it.
type Links []*Link

func (l Links) Len() int {
	return len(l)
}

func (l Links) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l Links) Less(i, j int) bool {
	return l[i].weight < l[j].weight
}

// Auxiliar function to return link between two nodes of transition if already exists or create one and return it.
func getLink(links Links, current, previous string) (*Link, bool)  {
	if len(links) > 0 {
		for _, link := range links {
			if link.current == current && link.previous == previous {
				return link, true
			}
		}
	}
	return &Link{current: current, previous: previous, occurrences: 1}, false
}

// Contains corpus and model name, available tags and transitions and emissions tables.
type Model struct {
	Name string
	Tags []string
	Transitions, Emissions Links
}

// Check if model exists locally and return it.
func Load(name string) (*Model, error) {
	var model *Model = &Model{Name: name}
	if exists, err := model.exists(); err != nil {
		return model, err
	} else if !exists {
		return model, errors.New("Model not found.")
	}

	model.loadLocal(name)
	return model, nil
}

// Check if model is already trained and saved locally to return it.
// If model doesn't exists yet, it'll be trained based on corpus provided and saved locally.
func (model *Model) Train(corpusName string) error {
	if exists, err := model.exists(); err != nil {
		return err
	} else if exists {
		model.loadLocal(corpusName)
		return nil
	}

	if c, err := GetCorpus(corpusName); err != nil {
		return err
	} else {
		model.trainCorpus(c)
		if err := model.saveLocal(); err != nil {
			return err
		}
	}
	return nil
}

// Calculate word possibilities based on previous tag, with transmission and emission costs using model provided.
func (model *Model) GetProbs(currentWord, prevTag string) map[string]float64 {
	var transitions Links
	for _, transition := range model.Transitions {
		if transition.previous == prevTag {
			transitions = append(transitions, transition)
		}
	}

	var emissions Links
	for _, emission := range model.Emissions {
		if emission.current == currentWord {
			emissions = append(emissions, emission)
		}
	}

	var probs map[string]float64 = make(map[string]float64, len(transitions))
	for _, emission := range emissions {
		var score float64 = emission.weight
		for _, transition := range transitions {
			if emission.current == transition.previous {
				score += transition.weight 
			}
		}
		probs[emission.previous] = score
	}

	return probs
}

// Check if model already exists
func (model *Model) exists() (bool, error) {
	var err error
	var location string

	if location, err = filepath.Abs(ModelsPath); err != nil {
		return false, err
	}

	var avalible []os.FileInfo
	if avalible, err = ioutil.ReadDir(location); err != nil {
		return false, err
	}

	for _, item := range avalible {
		if item.IsDir() && item.Name() == model.Name {
			return true, nil
		}
	}
	return false, nil
}

// Train model with corpus provided and generates transitions and emissions tables.
func (model *Model) trainCorpus(c Corpus) {
	fmt.Print("Training... ")

	var records []Record = c.Records
	var transitions, emissions Links

	var contextLength int = len(c.Tags) + 2
	var context map[string]float64 = make(map[string]float64, contextLength)

	//Calculate weights over all records
	for _, record := range records {
		var previous string = Start
		context[Start]++

		sort.Sort(record)
		for _, item := range record {
			if transition, exists := getLink(transitions, item.Tag, previous); exists {
				transition.occurrences++
			} else {
				transitions = append(transitions, transition)
			}

			if emission, exists := getLink(emissions, item.Raw, item.Tag); exists {
				emission.occurrences++
			} else {
				emissions = append(emissions, emission)
			}

			context[item.Tag]++
			previous = item.Tag
		}

		if transition, exists := getLink(transitions, previous, End); exists {
			transition.occurrences++
		} else {
			transitions = append(transitions, transition)
		}
		context[End]++
	}

	// Normalize weights
	for _, transition := range transitions {
		transition.weight = transition.occurrences / context[transition.previous] 
	}
	model.Transitions = transitions

	for _, emission := range emissions {
		emission.weight = emission.occurrences / context[emission.previous]
	}
	model.Emissions = emissions
	fmt.Println("Done!")
}

// Load model from local file system and parse transitions and emissions tables into data structures.
// All item was a Link, with current and previous state and their cost.
func (model *Model) loadLocal(corpusName string) error {
	var err error
	var transitionsPath string = fmt.Sprintf("%s/%s/transitions", ModelsPath, corpusName)
	var emissionsPath string = fmt.Sprintf("%s/%s/emissions", ModelsPath, corpusName)

	var tab_rgx *regexp.Regexp = regexp.MustCompile(`\t`)

	var transitions_fd *os.File
	if transitions_fd, err = os.Open(transitionsPath); err != nil {
		return err
	}
	defer transitions_fd.Close()

	var transitions Links
	var transitionScanner *bufio.Scanner = bufio.NewScanner(transitions_fd)
	for transitionScanner.Scan() {
		var line string = transitionScanner.Text()
		var data []string = tab_rgx.Split(line, -1)
		if len(data) == 3 {
			var weight float64
			if weight, err = strconv.ParseFloat(data[2], 64); err != nil {
				return err
			}

			var t *Link = &Link{previous: data[0], current: data[1], weight: weight}
			transitions = append(transitions, t)
		}
	}

	var emissions_fd *os.File
	if emissions_fd, err = os.Open(emissionsPath); err != nil {
		return err
	}
	defer emissions_fd.Close()

	var emissions Links
	var emissionScanner *bufio.Scanner = bufio.NewScanner(emissions_fd)
	for emissionScanner.Scan() {
		var line string = emissionScanner.Text()
		var data []string = tab_rgx.Split(line, -1)
		if len(data) == 3 {
			var weight float64
			if weight, err = strconv.ParseFloat(data[2], 64); err != nil {
				return err
			}

			var e *Link = &Link{current: data[1], previous: data[0], weight: weight}
			emissions = append(emissions, e)
		}
	}

	model.Transitions = transitions
	model.Emissions = emissions
	return nil
}

// Save trained model locally. Creates tabbed separated file with transitions and emissions and each weight.
func (model *Model) saveLocal() error {
	var err error
	var modelLocation string = fmt.Sprintf("%s/%s", ModelsPath, model.Name)
	if err = os.Mkdir(modelLocation, os.ModePerm); err != nil {
		return err
	}

	var transitionsPath string = fmt.Sprintf("%s/transitions", modelLocation)
	if fdTransitions, err := os.Create(transitionsPath); err == nil {
		defer fdTransitions.Close()

		for _, t := range model.Transitions {
			var line string = fmt.Sprintf("%s\t%s\t%g\n", t.previous, t.current, t.weight)
			if _, err = fdTransitions.WriteString(line); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	var emissionsPath string = fmt.Sprintf("%s/emissions", modelLocation)
	if fdEmissions, err := os.Create(emissionsPath); err == nil {
		defer fdEmissions.Close()

		for _, e := range model.Emissions {
			var line string = fmt.Sprintf("%s\t%s\t%g\n", e.previous, e.current, e.weight)
			if _, err = fdEmissions.WriteString(line); err != nil {
				return err
			}
		}
	} else {
		return err
	}
	return nil
}
