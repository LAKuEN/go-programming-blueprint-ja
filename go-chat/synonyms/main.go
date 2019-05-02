package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/LAKuEN/go-programming-blueprint-ja/go-chat/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := thesaurus.NewBigHuge(apiKey)

	s := bufio.NewScanner(os.Stdin)
	if !s.Scan() {
		panic("cannot read text from stdin")
	}
	word := s.Text()

	syns, err := thesaurus.Synonyms(word)
	if err != nil {
		panic(fmt.Errorf("cannot find synonyms of %s\nreason: %v", word, err))
	}

	for _, syn := range syns {
		fmt.Println(syn)
	}
}
