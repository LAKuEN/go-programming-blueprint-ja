package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

var transforms = []string{
	otherWord + "App",
	"get " + otherWord,
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherWord, sc.Text(), 1))
	}
}
