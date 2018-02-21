# Gopostager
HMM applied to Part-Of-Speech Tagging in Go.
    
## Installation
```
go get github.com/lucasmenendez/gopostagger
```


## Tested corpus

 Name | Language | Size | Link corpus
----- | ----- | ------ | ----
Brown | en | 11.6 Mb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/brown)
AnCora | es | 0.54 Mb | [Link](https://github.com/lucasmenendez/gopostagger/tree/master/ancora)

## Examples

### Tag sentence
```go
    package main

    import (
        "github.com/lucasmenendez/gotokenizer"
        "github.com/lucasmenendez/gopostagger"
        "fmt"
    )

    func main() {
        var s string = "El mundo del tatuaje es la forma de representación artística más expresiva que puede existir para un artista, puesto que su obra permanece inalterable de por vida."

        if m, e := gopostagger.LoadModel("./models/ancora"); e != nil {
            fmt.Println(e)
        } else {
            var tagger *gopostagger.Tagger = gopostagger.NewTagger(m)
            var tokens []string = gotokenizer.Words(s)
            var tagged [][]string = tagger.Tag(tokens)

            for _, token := range tagged {
                fmt.Printf("%q ", token)
            }
        }
    }
```

### Train corpus
IMPORTANT: All datasets must have the following format: `raw_word/tag_propossed`

```go
    package main

    import (
        "github.com/lucasmenendez/gopostagger"
        "fmt"
    )

    func main() {
        if m, e := gopostagger.Train("./corpus/ancora"); e != nil {
            fmt.Println(e)
        } else if e = m.Store("./models/ancora"); e != nil {
            fmt.Println(e)
        } else {
            fmt.Println("Trained!")
        }
    }
```
	
## Credits
- POS tagging Brown Corpus [Link](https://en.wikipedia.org/wiki/Brown_Corpus)
- Part-of-Speech Tagging with Hidden Markov Models - Graham Neubig [Link](http://www.phontron.com/slides/nlp-programming-en-04-hmm.pdf)