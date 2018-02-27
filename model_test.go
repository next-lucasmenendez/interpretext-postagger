package gopostagger

import (
	"net/http"
	"os"
	"testing"
	"fmt"
	"io"
)

const example_corpus string = "./example_corpus"
const example_model string = "./example_model"

func downloadCorpus(t *testing.T) bool {
	if fd, e := os.Create(example_corpus); e != nil {
		t.Errorf("Expected nil, got Error:'%s'", e)
		return false
	} else {
		defer fd.Close()

		if r, e := http.Get("https://raw.githubusercontent.com/lucasmenendez/gopostagger/master/ancora"); e != nil {
			t.Errorf("Expected nil, got Error:'%s'", e)
			return false
		} else {
			defer r.Body.Close()

			if _, e := io.Copy(fd, r.Body); e != nil {
				t.Errorf("Expected nil, got Error:'%s'", e)
				return false
			} else {
				return true
			}
		}
	}

	return false
}

func deleteCorpus(t *testing.T) bool {
	if e := os.Remove(example_corpus); e != nil {
		t.Errorf("Expected nil, got Error:'%s'", e)
		return false
	}

	return true
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
	} else if m, e := Train(example_corpus); e != nil {
		t.Fatalf("Expected nil, got %s", e.Error())
	} else if e := m.Store(example_model); e != nil {
		t.Errorf("Expected nil, got %s", e.Error())
	}
	return
}

func TestLoadTransitions(t *testing.T) {
	var (
		tp string = fmt.Sprintf("%s/transitions", example_model)
	 	m *Model = &Model{}
	)

	if err := m.loadTransitions("fail"); err == nil {
		t.Error("Expected error, got nil")
		return
	} else if err = m.loadTransitions(tp); err != nil {
		t.Errorf("Expected nil, got '%s'", err.Error())
		return
	} else {
		var expected links = links{
			{ previous: "<s>", current: "DET", weight: 0.3690260133091349 },
			{ previous: "DET", current: "NOUN", weight: 0.7237326536391957 },
			{ previous: "NOUN", current: "PUNCT", weight: 0.21732791835884618 },
			{ previous: "PUNCT", current: "ADP", weight: 0.11075697211155379 },
			{ previous: "ADP", current: "ADJ", weight: 0.020774647887323944 },
			{ previous: "ADJ", current: "NOUN", weight: 0.215375918598078 },
			{ previous: "NOUN", current: "ADP", weight: 0.3826929084661043 },
			{ previous: "ADP", current: "SCONJ", weight: 0.01784037558685446 },
			{ previous: "SCONJ", current: "VERB", weight: 0.15371024734982333 },
			{ previous: "VERB", current: "AUX", weight: 0.008492569002123142 },
			{ previous: "AUX", current: "NUM", weight: 0.013201320132013201 },
			{ previous: "NUM", current: "NOUN", weight: 0.5497142857142857 },
			{ previous: "ADP", current: "NOUN", weight: 0.26103286384976526 },
			{ previous: "ADP", current: "VERB", weight: 0.059741784037558684 },
			{ previous: "VERB", current: "ADP", weight: 0.2745930644019816 },
			{ previous: "NOUN", current: "DET", weight: 0.01437050921587004 },
			{ previous: "NOUN", current: "CCONJ", weight: 0.04956784338227637 },
			{ previous: "CCONJ", current: "ADJ", weight: 0.08459422283356259 },
			{ previous: "ADJ", current: "ADP", weight: 0.2798191068400226 },
			{ previous: "ADP", current: "DET", weight: 0.3926056338028169 },
			{ previous: "NOUN", current: "ADJ", weight: 0.17005102572112882 },
			{ previous: "ADJ", current: "PUNCT", weight: 0.24788015828151497 },
			{ previous: "PUNCT", current: "VERB", weight: 0.07474103585657371 },
			{ previous: "DET", current: "PRON", weight: 0.028603794958935145 },
			{ previous: "PRON", current: "DET", weight: 0.040522875816993466 },
			{ previous: "DET", current: "PROPN", weight: 0.09119229679977343 },
			{ previous: "PROPN", current: "PROPN", weight: 0.22966978807294233 },
			{ previous: "PROPN", current: "ADP", weight: 0.1429275505174963 },
			{ previous: "ADP", current: "PROPN", weight: 0.15246478873239436 },
			{ previous: "PROPN", current: "PUNCT", weight: 0.36717594874322323 },
			{ previous: "PUNCT", current: "PROPN", weight: 0.06262948207171315 },
			{ previous: "VERB", current: "VERB", weight: 0.025005897617362584 },
			{ previous: "VERB", current: "ADV", weight: 0.08185892899268696 },
			{ previous: "ADV", current: "DET", weight: 0.08452380952380953 },
			{ previous: "PUNCT", current: "ADJ", weight: 0.03219123505976096 },
			{ previous: "VERB", current: "PRON", weight: 0.052370842179759375 },
			{ previous: "DET", current: "NUM", weight: 0.0264797507788162 },
			{ previous: "NUM", current: "ADP", weight: 0.20114285714285715 },
			{ previous: "</s>", current: "PUNCT", weight: 0.9915305505142166 },
			{ previous: "<s>", current: "PROPN", weight: 0.11252268602540835 },
			{ previous: "PROPN", current: "PRON", weight: 0.025874815179891573 },
			{ previous: "PRON", current: "VERB", weight: 0.4522875816993464 },
			{ previous: "NOUN", current: "NOUN", weight: 0.009892741851504738 },
			{ previous: "DET", current: "ADV", weight: 0.0053809119229679975 },
			{ previous: "ADV", current: "ADJ", weight: 0.13452380952380952 },
			{ previous: "ADP", current: "PUNCT", weight: 0.004694835680751174 },
			{ previous: "VERB", current: "DET", weight: 0.2800188723755603 },
			{ previous: "CCONJ", current: "VERB", weight: 0.12448418156808803 },
			{ previous: "DET", current: "PUNCT", weight: 0.01047861795525347 },
			{ previous: "PUNCT", current: "NOUN", weight: 0.04653386454183267 },
			{ previous: "<s>", current: "ADP", weight: 0.16394434361766486 },
			{ previous: "ADP", current: "ADV", weight: 0.01431924882629108 },
			{ previous: "ADV", current: "VERB", weight: 0.15357142857142858 },
			{ previous: "<s>", current: "PRON", weight: 0.028433151845130067 },
			{ previous: "PRON", current: "AUX", weight: 0.14640522875816994 },
			{ previous: "AUX", current: "DET", weight: 0.12409240924092409 },
			{ previous: "DET", current: "ADJ", weight: 0.09274992919852733 },
			{ previous: "NOUN", current: "PRON", weight: 0.047068624388212014 },
			{ previous: "PUNCT", current: "CCONJ", weight: 0.05147410358565737 },
			{ previous: "CCONJ", current: "SCONJ", weight: 0.039889958734525444 },
			{ previous: "SCONJ", current: "PRON", weight: 0.10600706713780919 },
			{ previous: "VERB", current: "ADJ", weight: 0.024298183533852324 },
			{ previous: "ADJ", current: "CCONJ", weight: 0.06698699830412662 },
			{ previous: "CCONJ", current: "PRON", weight: 0.07771664374140302 },
			{ previous: "AUX", current: "ADP", weight: 0.0858085808580858 },
			{ previous: "NOUN", current: "AUX", weight: 0.022076434447568467 },
			{ previous: "AUX", current: "VERB", weight: 0.41254125412541254 },
			{ previous: "ADV", current: "ADV", weight: 0.05654761904761905 },
			{ previous: "PUNCT", current: "SCONJ", weight: 0.03059760956175299 },
			{ previous: "SCONJ", current: "DET", weight: 0.3083038869257951 },
			{ previous: "PUNCT", current: "PUNCT", weight: 0.10342629482071714 },
			{ previous: "CCONJ", current: "PROPN", weight: 0.10866574965612105 },
			{ previous: "PUNCT", current: "PRON", weight: 0.06406374501992032 },
			{ previous: "AUX", current: "PRON", weight: 0.02574257425742574 },
			{ previous: "PRON", current: "ADV", weight: 0.05403050108932462 },
			{ previous: "DET", current: "DET", weight: 0.013027470971396206 },
			{ previous: "VERB", current: "SCONJ", weight: 0.07289455060155697 },
			{ previous: "ADJ", current: "ADJ", weight: 0.030243075183719617 },
			{ previous: "NUM", current: "CCONJ", weight: 0.02857142857142857 },
			{ previous: "CCONJ", current: "DET", weight: 0.19944979367262725 },
			{ previous: "VERB", current: "NOUN", weight: 0.051663128096249115 },
			{ previous: "CCONJ", current: "NOUN", weight: 0.11485557083906466 },
			{ previous: "NOUN", current: "ADV", weight: 0.02082682495053629 },
			{ previous: "ADV", current: "ADP", weight: 0.18571428571428572 },
			{ previous: "SCONJ", current: "PROPN", weight: 0.06007067137809187 },
			{ previous: "PROPN", current: "VERB", weight: 0.08624938393297191 },
			{ previous: "<s>", current: "SCONJ", weight: 0.024198427102238355 },
			{ previous: "PROPN", current: "ADV", weight: 0.014785608674223755 },
			{ previous: "ADV", current: "PRON", weight: 0.0744047619047619 },
			{ previous: "ADP", current: "AUX", weight: 0.013615023474178404 },
			{ previous: "AUX", current: "SCONJ", weight: 0.044884488448844885 },
			{ previous: "VERB", current: "PUNCT", weight: 0.07926397735314933 },
			{ previous: "PUNCT", current: "ADV", weight: 0.03856573705179283 },
			{ previous: "ADV", current: "PUNCT", weight: 0.10595238095238095 },
			{ previous: "PROPN", current: "AUX", weight: 0.030556924593395762 },
			{ previous: "ADP", current: "ADP", weight: 0.0028169014084507044 },
			{ previous: "ADJ", current: "PROPN", weight: 0.02798191068400226 },
			{ previous: "ADJ", current: "ADV", weight: 0.017241379310344827 },
			{ previous: "ADJ", current: "DET", weight: 0.014697569248162803 },
			{ previous: "CCONJ", current: "ADP", weight: 0.09903713892709766 },
			{ previous: "PRON", current: "PRON", weight: 0.08714596949891068 },
			{ previous: "SCONJ", current: "AUX", weight: 0.09363957597173145 },
			{ previous: "<s>", current: "ADV", weight: 0.06352087114337568 },
			{ previous: "PUNCT", current: "DET", weight: 0.08286852589641434 },
			{ previous: "CCONJ", current: "AUX", weight: 0.024759284731774415 },
			{ previous: "AUX", current: "AUX", weight: 0.058085808580858087 },
			{ previous: "AUX", current: "NOUN", weight: 0.04026402640264026 },
			{ previous: "NOUN", current: "PROPN", weight: 0.016140789336665625 },
			{ previous: "CCONJ", current: "ADV", weight: 0.07359009628610728 },
			{ previous: "ADV", current: "SCONJ", weight: 0.0494047619047619 },
			{ previous: "<s>", current: "PUNCT", weight: 0.0780399274047187 },
			{ previous: "AUX", current: "ADJ", weight: 0.10957095709570958 },
			{ previous: "SCONJ", current: "ADV", weight: 0.07773851590106007 },
			{ previous: "PRON", current: "ADP", weight: 0.09803921568627451 },
			{ previous: "ADJ", current: "PRON", weight: 0.02798191068400226 },
			{ previous: "VERB", current: "PROPN", weight: 0.021467327199811276 },
			{ previous: "<s>", current: "VERB", weight: 0.03932244404113733 },
			{ previous: "PROPN", current: "CCONJ", weight: 0.06037456875308034 },
			{ previous: "ADV", current: "NOUN", weight: 0.02738095238095238 },
			{ previous: "</s>", current: "ADJ", weight: 0.0012099213551119178 },
			{ previous: "NOUN", current: "VERB", weight: 0.03446839529313756 },
			{ previous: "<s>", current: "ADJ", weight: 0.010889292196007259 },
			{ previous: "<s>", current: "CCONJ", weight: 0.03327283726557774 },
			{ previous: "ADJ", current: "VERB", weight: 0.0322215941209723 },
			{ previous: "SCONJ", current: "PUNCT", weight: 0.05653710247349823 },
			{ previous: "DET", current: "SCONJ", weight: 0.0004248088360237893 },
			{ previous: "ADJ", current: "AUX", weight: 0.017241379310344827 },
			{ previous: "ADJ", current: "SCONJ", weight: 0.018089315997738834 },
			{ previous: "ADV", current: "AUX", weight: 0.1005952380952381 },
			{ previous: "PUNCT", current: "AUX", weight: 0.028525896414342628 },
			{ previous: "SCONJ", current: "SCONJ", weight: 0.015901060070671377 },
			{ previous: "PROPN", current: "ADJ", weight: 0.015278462296697881 },
			{ previous: "</s>", current: "PROPN", weight: 0.004839685420447671 },
			{ previous: "<s>", current: "AUX", weight: 0.02722323049001815 },
			{ previous: "VERB", current: "NUM", weight: 0.017692852087756547 },
			{ previous: "<s>", current: "NOUN", weight: 0.02722323049001815 },
			{ previous: "AUX", current: "ADV", weight: 0.052805280528052806 },
			{ previous: "PRON", current: "PUNCT", weight: 0.06710239651416122 },
			{ previous: "PRON", current: "CCONJ", weight: 0.007407407407407408 },
			{ previous: "CCONJ", current: "PUNCT", weight: 0.02544704264099037 },
			{ previous: "ADP", current: "PRON", weight: 0.01936619718309859 },
			{ previous: "SCONJ", current: "NOUN", weight: 0.05035335689045936 },
			{ previous: "SCONJ", current: "ADP", weight: 0.06802120141342756 },
			{ previous: "AUX", current: "PUNCT", weight: 0.02706270627062706 },
			{ previous: "NUM", current: "PUNCT", weight: 0.12571428571428572 },
			{ previous: "PRON", current: "ADJ", weight: 0.011328976034858388 },
			{ previous: "ADP", current: "NUM", weight: 0.038849765258215964 },
			{ previous: "NUM", current: "ADV", weight: 0.004571428571428572 },
			{ previous: "NUM", current: "PROPN", weight: 0.009142857142857144 },
			{ previous: "ADV", current: "PROPN", weight: 0.005952380952380952 },
			{ previous: "NOUN", current: "PART", weight: 0.0002082682495053629 },
			{ previous: "PART", current: "ADV", weight: 0.4 },
			{ previous: "PRON", current: "NUM", weight: 0.014814814814814815 },
			{ previous: "PRON", current: "SCONJ", weight: 0.004357298474945534 },
			{ previous: "ADV", current: "CCONJ", weight: 0.01488095238095238 },
			{ previous: "NOUN", current: "SCONJ", weight: 0.01072581484952619 },
			{ previous: "ADJ", current: "NUM", weight: 0.003674392312040701 },
			{ previous: "NOUN", current: "NUM", weight: 0.004061230865354577 },
			{ previous: "PRON", current: "PROPN", weight: 0.00392156862745098 },
			{ previous: "PUNCT", current: "INTJ", weight: 0.00015936254980079682 },
			{ previous: "INTJ", current: "PUNCT", weight: 1 },
			{ previous: "VERB", current: "CCONJ", weight: 0.009907997169143666 },
			{ previous: "PRON", current: "NOUN", weight: 0.011764705882352941 },
			{ previous: "DET", current: "ADP", weight: 0.004389691305579156 },
			{ previous: "AUX", current: "CCONJ", weight: 0.0019801980198019802 },
			{ previous: "PROPN", current: "DET", weight: 0.01232134056185313 },
			{ previous: "PUNCT", current: "PART", weight: 0.0009561752988047809 },
			{ previous: "PROPN", current: "SCONJ", weight: 0.007639231148348941 },
			{ previous: "<s>", current: "NUM", weight: 0.019963702359346643 },
			{ previous: "NUM", current: "DET", weight: 0.005714285714285714 },
			{ previous: "<s>", current: "SYM", weight: 0.0006049606775559589 },
			{ previous: "SYM", current: "PUNCT", weight: 0.18181818181818182 },
			{ previous: "PUNCT", current: "SYM", weight: 0.00015936254980079682 },
			{ previous: "PUNCT", current: "NUM", weight: 0.011155378486055778 },
			{ previous: "CCONJ", current: "NUM", weight: 0.026822558459422285 },
			{ previous: "ADV", current: "NUM", weight: 0.006547619047619048 },
			{ previous: "AUX", current: "PROPN", weight: 0.0033003300330033004 },
			{ previous: "NUM", current: "ADJ", weight: 0.03428571428571429 },
			{ previous: "NUM", current: "AUX", weight: 0.006857142857142857 },
			{ previous: "NUM", current: "NUM", weight: 0.010285714285714285 },
			{ previous: "PROPN", current: "NOUN", weight: 0.0014785608674223755 },
			{ previous: "DET", current: "VERB", weight: 0.0002832058906825262 },
			{ previous: "PROPN", current: "NUM", weight: 0.0034499753573188764 },
			{ previous: "NUM", current: "VERB", weight: 0.017142857142857144 },
			{ previous: "PART", current: "NOUN", weight: 0.4666666666666667 },
			{ previous: "NUM", current: "SCONJ", weight: 0.001142857142857143 },
			{ previous: "SCONJ", current: "ADJ", weight: 0.007950530035335688 },
			{ previous: "<s>", current: "PART", weight: 0.0018148820326678765 },
			{ previous: "SCONJ", current: "CCONJ", weight: 0.0017667844522968198 },
			{ previous: "VERB", current: "PART", weight: 0.00023590469450342062 },
			{ previous: "NUM", current: "PRON", weight: 0.004571428571428572 },
			{ previous: "</s>", current: "NUM", weight: 0.0006049606775559589 },
			{ previous: "</s>", current: "NOUN", weight: 0.0018148820326678765 },
			{ previous: "DET", current: "AUX", weight: 0.0001416029453412631 },
			{ previous: "ADP", current: "PART", weight: 0.00011737089201877934 },
			{ previous: "PART", current: "PROPN", weight: 0.06666666666666667 },
			{ previous: "DET", current: "SYM", weight: 0.002973661852166525 },
			{ previous: "SYM", current: "ADP", weight: 0.5681818181818182 },
			{ previous: "ADP", current: "SYM", weight: 0.0017605633802816902 },
			{ previous: "SYM", current: "ADV", weight: 0.022727272727272728 },
			{ previous: "SYM", current: "ADJ", weight: 0.022727272727272728 },
			{ previous: "AUX", current: "SYM", weight: 0.0006600660066006601 },
			{ previous: "SYM", current: "NOUN", weight: 0.06818181818181818 },
			{ previous: "NOUN", current: "SYM", weight: 0.0002082682495053629 },
			{ previous: "SYM", current: "CCONJ", weight: 0.06818181818181818 },
			{ previous: "PRON", current: "SYM", weight: 0.00043572984749455336 },
			{ previous: "DET", current: "CCONJ", weight: 0.0001416029453412631 },
			{ previous: "PROPN", current: "PART", weight: 0.0002464268112370626 },
			{ previous: "SYM", current: "VERB", weight: 0.045454545454545456 },
			{ previous: "VERB", current: "SYM", weight: 0.00023590469450342062 },
			{ previous: "PRON", current: "PART", weight: 0.00043572984749455336 },
			{ previous: "CCONJ", current: "SYM", weight: 0.000687757909215956 },
			{ previous: "SYM", current: "DET", weight: 0.022727272727272728 },
			{ previous: "PART", current: "ADP", weight: 0.06666666666666667 },
		}

		for _, tm := range m.transitions {
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

func TestLoadEmissions(t *testing.T) {}

func TestLoadModel(t *testing.T) {}

func TestProbs(t *testing.T) {}

func TestScore(t *testing.T) {}
