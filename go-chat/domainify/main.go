package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

var re = regexp.MustCompile("^[a-z0-9-_]+")

func containInvalidCharacter(s []byte) bool {
	return len(re.Find(s)) != len(s)
}

var tlds = []string{"com", "net"}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		text := strings.ToLower(sc.Text())
		text = strings.ReplaceAll(text, " ", "-")

		if containInvalidCharacter([]byte(text)) {
			panic("specified text contains invalid character.")
		}

		tld := tlds[rand.Intn(2)]
		fmt.Println(text + "." + tld)
	}
}
