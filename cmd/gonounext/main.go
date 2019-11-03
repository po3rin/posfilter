package main

import (
	"fmt"
	"log"

	"github.com/po3rin/posfilter"
)

func main() {
	var filter posfilter.PosFilter
	words, err := filter.Do("東京都へ行く")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(words)
}
