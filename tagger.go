// HMM applied to Part-Of-Speech Tagging in Go.
package gopostagger

import (
	"sort"
	"strings"
)

// token is a struct to contain token information including sentence order, raw
// content and proposed tag
type token struct {
	order    int
	raw string
	tag string
}

// sentences contains list of tokens pointers
type sentence []*token

func (s sentence) Len() int           { return len(s) }
func (s sentence) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sentence) Less(i, j int) bool { return s[i].order < s[j].order }

// Tagger struct is associated to a model
type Tagger struct {
	model *Model
}

// NewTagger function returns a tagger instance associated with model provided.
func NewTagger(m *Model) *Tagger {
	return &Tagger{m}
}

// Tag function proposes a tag for each tokens provided in based of tagger model.
func (t *Tagger) Tag(tokens []string) (tagged [][]string) {
	var s sentence
	for i, w := range tokens {
		s = append(s, &token{ order: i, raw: w })
	}

	sort.Sort(s)
	var c string = StartTag
	for _, tk := range s {
		var max float64
		var lt string = strings.ToLower(tk.raw)
		var ps map[string]float64 = t.model.probs(lt, c)
		if len(ps) > 0 {
			for tg, sc := range ps {
				if sc > max {
					c = tg
					max = sc
				}
			}
			tk.tag = c
		}

		tagged = append(tagged, []string{tk.raw, tk.tag})
	}
	return tagged
}
