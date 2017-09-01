# Gopostager
HMM applied to Part-Of-Speech Tagging in Go.
    
## Installation
```
go get github.com/lucasmenendez/gopostagger
```

## Pretrained corpus

Name | Alias to use | Language | Size | Link to model | Download corpus 
---- | ------------ | -------- | ---- | ------------- | --------------- 
Brown | brown | English | 12.3 Mb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/models/brown)  | [Link](https://drive.google.com/open?id=0B6YI1HgpnJOjTE5hbFhwVUhDR0k)
Cess | cess | Spanish | 471.1 Mb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/models/cess) | [Link](https://drive.google.com/file/d/0B6YI1HgpnJOjVTFDeVZ4ZjhXU28/view?usp=sharing)
AnCora | ancora | Spanish | 550 Kb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/models/ancora) | [Link](https://drive.google.com/file/d/0B6YI1HgpnJOjUzBEbnZjQW93YlU/view?usp=sharing)
Wikipedia | wikipedia-en | English | 1.6 Mb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/models/wikipedia-en) | [Link](https://drive.google.com/file/d/0B6YI1HgpnJOjTXM0d3V2aGJhMWM/view)
Wikipedia | wikipedia-es | Spanish | 2.7 Mb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/models/wikipedia-es) | [Link](https://drive.google.com/open?id=0B6YI1HgpnJOjel94WFZRYTNfc28)
Public domain books | books-es | Spanish | 500 Kb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/models/books-es) | [Link](https://drive.google.com/open?id=0B6YI1HgpnJOjakpESTZNc2RGU2M)

## Work with others corpus

Place all datasets files into a folder on `/corpus/<corpus_name>/`. 

IMPORTANT: All datasets must have the following format: `raw_word/tag_propossed` 


## Example

```go
package main

import (
	"fmt"
	g "github.com/lucasmenendez/gopostagger"
)

func main() {
	model := "wikipedia-es" //Set here model alias
	rawSentence := "En un lugar de la Mancha, de cuyo nombre no quiero acordarme, no ha mucho tiempo que vivía un hidalgo de los de lanza en astillero, adarga antigua, rocín flaco y galgo corredor."

	if s, err := g.TagSentence(rawSentence, model); err != nil {
		fmt.Println(err)
	} else {
	    for _, word := range s {
    		fmt.Printf("%s/%s ", word.Raw, word.Tag)
    	}
	}
}
```
	
## Credits
- POS tagging Brown Corpus [Link](https://en.wikipedia.org/wiki/Brown_Corpus)
- Part-of-Speech Tagging with Hidden Markov Models - Graham Neubig [Link](http://www.phontron.com/slides/nlp-programming-en-04-hmm.pdf)