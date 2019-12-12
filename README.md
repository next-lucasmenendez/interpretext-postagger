[![GoDoc](https://godoc.org/github.com/next-lucasmenendez/interpretext-postagger?status.svg)](https://godoc.org/github.com/next-lucasmenendez/interpretext-postagger)
[![Report](https://goreportcard.com/badge/github.com/next-lucasmenendez/interpretext-postagger)](https://goreportcard.com/report/github.com/next-lucasmenendez/interpretext-postagger)

# Gopostager
HMM applied to Part-Of-Speech Tagging in Go. Implementation of [*Part-of-Speech Tagging with Hidden Markov Models - Graham Neubig*](http://www.phontron.com/slides/nlp-programming-en-04-hmm.pdf)
    
## Installation
```
go get github.com/next-lucasmenendez/interpretext-postagger
```


## Tested corpus

 Name | Language | Size | Link corpus
----- | ----- | ------ | ----
Brown | en | 11.6 Mb | [Link](https://github.com/next-lucasmenendez/interpretext-postagger/tree/master/en)
AnCora | es | 0.54 Mb | [Link](https://github.com/next-lucasmenendez/interpretext-postagger/tree/master/es)

## Examples

### Tag sentence
```go
    package main

    import (
        tokenizer "github.com/next-lucasmenendez/interpretext-tokenizer"
        postagger "github.com/next-lucasmenendez/interpretext-postagger"
        "fmt"
    )

    func main() {
        var s string = "El mundo del tatuaje es la forma de representación artística más expresiva que puede existir para un artista, puesto que su obra permanece inalterable de por vida."

        if m, e := postagger.LoadModel("./models/es"); e != nil {
            fmt.Println(e)
        } else {
            var tagger *postagger.Tagger = postagger.NewTagger(m)
            var tokens []string = tokenizer.Words(s)
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
        postagger "github.com/next-lucasmenendez/interpretext-postagger"
        "fmt"
    )

    func main() {
        if m, e := postagger.Train("./es"); e != nil {
            fmt.Println(e)
        } else if e = m.Store("./models/es"); e != nil {
            fmt.Println(e)
        } else {
            fmt.Println("Trained!")
        }
    }
```