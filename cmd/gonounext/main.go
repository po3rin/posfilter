package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/po3rin/posfilter"
)

func main() {
	var filter posfilter.PosFilter

	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()

	words, err := filter.Do(text)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(words)
}
