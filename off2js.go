package main

import (
	"bufio"
	"encoding/json"
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

func getTriangles(verts [][3]float64, faces [][]string, BC *[][3]int) [][3]float64{
	var triangles [][3]float64
	for i := 0; i < len(faces); i++ {
		numPoints := toInt(faces[i][0])
		point0 := toInt(faces[i][1])
		for j := 1; j <= numPoints - 2; j++ {
			point1 := toInt(faces[i][j+1])
			point2 := toInt(faces[i][j+2])
			triangles = append(triangles, verts[point0], verts[point1], verts[point2])
			*BC = append(*BC, [...]int{1, 0, 0}, [...]int{0, 1, 0}, [...]int{0, 0, 1})
		}
	}
	return triangles
}

func main() {
	inPtr := flag.String("i", "foo.off", "the desired input file name")
	outPtr := flag.String("o", "foo.js", "the desired output file name")
	modelPtr := flag.String("modelName", "CubeModel", "the desired model name")

	flag.Parse()

	vertices := [...]float64{0, 0, 0}
	var vertexArray [][3]float64

	var faces [][]string
	var BC [][3]int

	fmt.Println("in: ", *inPtr)
	fmt.Println("out: ", *outPtr)
	fmt.Println("model: ", *modelPtr)

	input, err := os.Open(*inPtr)
	check(err)
	defer input.Close()

	scanner := bufio.NewScanner(input)
	scanner.Scan()

	if strings.ToUpper(scanner.Text()) == "OFF" {
		var numVertices int
		var numFaces int
		gotCounts := false
		verticesRead := false
		i := 0
		for scanner.Scan() {
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
						for k := 0; k < 3; k++ {
							vert, err := strconv.ParseFloat(verts[k], 64)
							check(err)
							vertices[k] = vert
						}
						vertexArray = append(vertexArray, vertices)
						i++
						if i == numVertices { verticesRead = true }
					} else {
						faces = append(faces, strings.Fields(scanner.Text()))
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
	something := getTriangles(vertexArray, faces, &BC)

	tri, _ := json.Marshal(something)
	bary, _ := json.Marshal(BC)

	funcString := []string{"function ", *modelPtr, "() {\n\tthis.triangles = ", string(tri), ";\n\tthis.BC = ", string(bary), "; };"}

	outFile, err := os.Create(*outPtr)
	check(err)
	defer outFile.Close()

	outFile.WriteString(strings.Join(funcString, ""))
}
