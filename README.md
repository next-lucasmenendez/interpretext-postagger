# Gopostager
HMM applied to Part-Of-Speech Tagging in Go.
    
## Installation
```
go get github.com/lucasmenendez/gopostagger
```

## Pretrained corpus

Language | Alias | Corpus | Size | Link to model | Download corpus 
-------- | ----- | ------ | ---- | ------------- | --------------- 
English | en | Brown and Wikipedia | 12.1 Mb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/models/en)  | [Link](https://drive.google.com/file/d/0B6YI1HgpnJOjcmRhaDc3MjhiSlk)
Spanish | es | AnCora, Wikipedia and Public domain books  | 3.54 Mb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/models/es) | [Link](https://drive.google.com/open?id=0B6YI1HgpnJOjZzJFWWNEamFubUU)

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
        model := "es" //Set here model alias
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