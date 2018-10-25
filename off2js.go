package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func main() {
	inPtr := flag.String("i", "foo.off", "the desired input file name")
	outPtr := flag.String("o", "foo.js", "the desired output file name")
	modelPtr := flag.String("m", "FooModel", "the desired model name")

	flag.Parse()

	fmt.Println("in: ", *inPtr)
	fmt.Println("out: ", *outPtr)
	fmt.Println("model: ", *modelPtr)

	input, err := os.Open(*inPtr)
	check(err)
	defer input.Close()

	scanner := bufio.NewScanner(input)
	scanner.Scan()

	if strings.ToUpper(scanner.Text()) == "OFF" {
		scanner.Scan()
		lineCounts := strings.Fields(scanner.Text())
		numVertices := toInt(lineCounts[0])
		numFaces := toInt(lineCounts[1])

		i := 0
		for scanner.Scan() {

			i++
			if i > numVertices {
				break
			}
		}

		i = 0
		for scanner.Scan() {
			i++
			if i > numFaces {
				break
			}
		}
	}
}
