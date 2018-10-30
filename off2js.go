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

	vertices := []float64{0, 0, 0}

	fmt.Println("in: ", *inPtr)
	fmt.Println("out: ", *outPtr)
	fmt.Println("model: ", *modelPtr)

	input, err := os.Open(*inPtr)
	check(err)
	defer input.Close()

	scanner := bufio.NewScanner(input)
	scanner.Scan()

	// If this is an OFF file
	if strings.ToUpper(scanner.Text()) == "OFF" {
		var numVertices int
		var numFaces int
		gotCounts := false
		verticesRead := false
		i := 0
		for scanner.Scan() {
			// Do stuff if the line is not a comment
			if !(string(scanner.Text()[0]) == "#"){
				if !gotCounts {
					nums := strings.Fields(scanner.Text())
					numVertices = toInt(nums[0])
					numFaces = toInt(nums[1])
					gotCounts = true
					fmt.Println(numVertices, numFaces)
				} else if numVertices != 0 && numFaces != 0 {
					if !verticesRead {
						verts := strings.Fields(scanner.Text())
						for k := 0; k < 3; k++{
							vert, err := strconv.ParseFloat(verts[k], 64)
							check(err)
							vertices[k] = vert
						}
						check(err)
						fmt.Println(vertices)
						i++
						if i == numVertices { verticesRead = true }
					} else {
						fmt.Println("Face: ", scanner.Text())
					}
				} else {
					fmt.Println("Cannot define a model")
					os.Exit(0)
				}
			}
		}
	} else {
		fmt.Println("No OFF header found")
		os.Exit(0)
	}
}
