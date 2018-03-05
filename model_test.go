package gopostagger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

const exampleCorpus string = "./example_corpus"
const exampleModel string = "./example_model"

func downloadCorpus(t *testing.T) bool {
	if fd, e := os.Create(exampleCorpus); e != nil {
		t.Errorf("Expected nil, got Error:'%s'", e)
	} else {
		defer fd.Close()

		if r, e := http.Get("https://raw.githubusercontent.com/lucasmenendez/gopostagger/master/ancora"); e != nil {
			t.Errorf("Expected nil, got Error:'%s'", e)
		} else {
			defer r.Body.Close()

			if _, e := io.Copy(fd, r.Body); e != nil {
				t.Errorf("Expected nil, got Error:'%s'", e)
			} else {
				return true
			}
		}
	}

	return false
}

func deleteCorpus(t *testing.T) {
	if e := os.Remove(exampleCorpus); e != nil {
		t.Errorf("Expected nil, got Error:'%s'", e)
	} else if e = os.RemoveAll(exampleModel); e != nil {
		t.Errorf("Expected nil, got Error:'%s'", e)
	}
}

func TestGetLink(t *testing.T) {
	var ls links = []*link{
		{"AA", "BB", 0, 0},
	}

	if _, exists := ls.getLink("AA", "CC"); exists {
		t.Error("Expected false, got true")
		return
	}

	if r, exists := ls.getLink("AA", "BB"); !exists {
		t.Error("Expected true, got false")
		return
	} else if e := ls[0]; r.current != e.current || r.previous != e.previous {
		t.Error("Expected true, got false")
		return
	}
}

func TestTrainAndStore(t *testing.T) {
	if _, e := Train("fail"); e == nil {
		t.Error("Expected error, got nil")
		return
	}

	if ok := downloadCorpus(t); !ok {
		t.Error("Expected true, got false")
	} else if m, e := Train(exampleCorpus); e != nil {
		t.Fatalf("Expected nil, got %s", e.Error())
	} else if e := m.Store(exampleModel); e != nil {
		t.Errorf("Expected nil, got %s", e.Error())
	}
}

func TestLoadTransitions(t *testing.T) {
	var (
		tp string = fmt.Sprintf("%s/transitions", exampleModel)
		m  *Model = &Model{}
	)

	if err := m.loadTransitions("fail"); err == nil {
		t.Error("Expected error, got nil")
		return
	} else if err = m.loadTransitions(tp); err != nil {
		t.Errorf("Expected nil, got '%s'", err.Error())
		return
	} else {
		var expected links = links{
			{previous: "<s>", current: "DET", weight: 0.3690260133091349},
			{previous: "DET", current: "NOUN", weight: 0.7237326536391957},
			{previous: "NOUN", current: "PUNCT", weight: 0.21732791835884618},
			{previous: "PUNCT", current: "ADP", weight: 0.11075697211155379},
			{previous: "ADP", current: "ADJ", weight: 0.020774647887323944},
			{previous: "ADJ", current: "NOUN", weight: 0.215375918598078},
			{previous: "NOUN", current: "ADP", weight: 0.3826929084661043},
			{previous: "ADP", current: "SCONJ", weight: 0.01784037558685446},
			{previous: "SCONJ", current: "VERB", weight: 0.15371024734982333},
			{previous: "VERB", current: "AUX", weight: 0.008492569002123142},
		}

		for _, tm := range m.transitions[:10] {
			var in bool = false
			for _, te := range expected {
				in = in || (tm.previous == te.previous && tm.current == te.current && tm.weight == te.weight)
			}

			if !in {
				t.Error("Expected true, got false")
				return
			}
		}
	}
}

func TestLoadEmissions(t *testing.T) {
	var (
		ep string = fmt.Sprintf("%s/emissions", exampleModel)
		m  *Model = &Model{}
	)

	if err := m.loadEmissions("fail"); err == nil {
		t.Error("Expected error, got nil")
		return
	} else if err = m.loadEmissions(ep); err != nil {
		t.Errorf("Expected nil, got '%s'", err.Error())
		return
	} else {
		var expected links = links{
			{previous: "DET", current: "El", weight: 0.04290569243840272},
			{previous: "NOUN", current: "gobernante", weight: 0.0002082682495053629},
			{previous: "PUNCT", current: ",", weight: 0.45322709163346614},
			{previous: "ADP", current: "con", weight: 0.05},
			{previous: "ADJ", current: "ganada", weight: 0.0002826455624646693},
			{previous: "NOUN", current: "fama", weight: 0.00010413412475268145},
			{previous: "ADP", current: "desde", weight: 0.007863849765258215},
			{previous: "SCONJ", current: "que", weight: 0.6351590106007067},
			{previous: "VERB", current: "llegó", weight: 0.001887237556027365},
			{previous: "AUX", current: "hace", weight: 0.013861386138613862},
		}

		for _, em := range m.emissions[:10] {
			var in bool = false
			for _, ee := range expected {
				in = in || (em.previous == ee.previous && em.current == ee.current && em.weight == ee.weight)
			}

			if !in {
				t.Error("Expected true, got false")
				return
			}
		}
	}
}

func TestLoadModel(t *testing.T) {
	if _, err := LoadModel("dslknfdslkngsd"); err == nil {
		t.Error("Expected error, got nil")
	} else if m, err := LoadModel(exampleModel); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		var expt links = links{
			{previous: "<s>", current: "DET", weight: 0.3690260133091349},
			{previous: "DET", current: "NOUN", weight: 0.7237326536391957},
			{previous: "NOUN", current: "PUNCT", weight: 0.21732791835884618},
			{previous: "PUNCT", current: "ADP", weight: 0.11075697211155379},
			{previous: "ADP", current: "ADJ", weight: 0.020774647887323944},
			{previous: "ADJ", current: "NOUN", weight: 0.215375918598078},
			{previous: "NOUN", current: "ADP", weight: 0.3826929084661043},
			{previous: "ADP", current: "SCONJ", weight: 0.01784037558685446},
			{previous: "SCONJ", current: "VERB", weight: 0.15371024734982333},
			{previous: "VERB", current: "AUX", weight: 0.008492569002123142},
		}

		for _, tm := range m.transitions[:10] {
			var in bool = false
			for _, te := range expt {
				in = in || (tm.previous == te.previous && tm.current == te.current && tm.weight == te.weight)
			}

			if !in {
				t.Error("Expected true, got false")
				return
			}
		}

		var expe links = links{
			{previous: "DET", current: "El", weight: 0.04290569243840272},
			{previous: "NOUN", current: "gobernante", weight: 0.0002082682495053629},
			{previous: "PUNCT", current: ",", weight: 0.45322709163346614},
			{previous: "ADP", current: "con", weight: 0.05},
			{previous: "ADJ", current: "ganada", weight: 0.0002826455624646693},
			{previous: "NOUN", current: "fama", weight: 0.00010413412475268145},
			{previous: "ADP", current: "desde", weight: 0.007863849765258215},
			{previous: "SCONJ", current: "que", weight: 0.6351590106007067},
			{previous: "VERB", current: "llegó", weight: 0.001887237556027365},
			{previous: "AUX", current: "hace", weight: 0.013861386138613862},
		}

		for _, em := range m.emissions[:10] {
			var in bool = false
			for _, ee := range expe {
				in = in || (em.previous == ee.previous && em.current == ee.current && em.weight == ee.weight)
			}

			if !in {
				t.Error("Expected true, got false")
				return
			}
		}
	}
}

func TestProbs(t *testing.T) {
	defer deleteCorpus(t)

	if m, err := LoadModel(exampleModel); err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	} else {
		if ps, sg := m.probs("fdsfsf", ""); len(ps) > 0 {
			t.Errorf("Expecetd 0, got %d", len(ps))
		} else if sg != "<s>?" {
			t.Errorf("Expected \"<s>?\", got \"%s\"", sg)
		} else if ps, sg := m.probs("El", "<s>"); len(ps) == 0 {
			t.Error("Expecetd >0, got 0")
		} else if ps["DET"] != 0.04290569243840272 {
			t.Errorf("Expected \"DET\":0.04290569243840272, got \"%+v\"", ps)
		} else if sg != "" {
			t.Errorf("Expected \"\", got \"%s\"", sg)
		}
	}
}
