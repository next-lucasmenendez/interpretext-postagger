package gopostagger

import "testing"

func TestTaggerAndTag(t *testing.T) {
	if ok := downloadCorpus(t); !ok {
		t.Error("Expected true, got false")
	} else if m, e := Train(example_corpus); e != nil {
		t.Fatalf("Expected nil, got %s", e.Error())
	} else if e := m.Store(example_model); e != nil {
		t.Errorf("Expected nil, got %s", e.Error())
	} else {
		defer deleteCorpus(t)

		if m, e := LoadModel(example_model); e != nil {
			t.Errorf("Expected nil, got %s", e.Error())
		} else {
			var (
				tagger   *Tagger    = NewTagger(m)
				sentence []string   = []string{"El", "mundo", "está", "girando", "sobre", "sí", "mismo", "."}
				expected [][]string = [][]string{
					{"El", "DET"},
					{"mundo", "NOUN"},
					{"está", "AUX"},
					{"girando", "VERB?"},
					{"sobre", "ADP"},
					{"sí", "ADV"},
					{"mismo", "ADJ"},
					{".", "PUNCT"},
				}
			)

			var res [][]string = tagger.Tag(sentence)
			for i, e := range expected {
				var r []string = res[i]

				if e[0] != r[0] || e[1] != r[1] {
					t.Errorf("Expected %q, got %q", e, r)
				}
			}
		}
	}
}
