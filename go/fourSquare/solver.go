package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

type individual struct {
	a     [25]byte
	b     [25]byte
	c     [25]byte
	d     [25]byte
	score float64
}

type coords struct {
	row    int
	column int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func format(original []byte) (ciphertext string) {
	for i := 0; i < len(original); i++ {
		for j := 0; j < len(letters()); j++ {
			if original[i] == letters()[j] {
				ciphertext += string(original[i])
			}
		}
	}
	return ciphertext
}

// !const
var tg map[string]float64 = tetragrams() // used instead of initialiser function so the massive file is only read only

func tetragrams() map[string]float64 {
	tetragrams := map[string]float64{}
	text, err := os.ReadFile("../logTetragrams.txt")
	check(err)

	scanner := bufio.NewScanner(strings.NewReader(string(text)))

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ", ")
		tetragram := fields[0]
		frequency, err := strconv.ParseFloat(fields[1], 64)
		check(err)

		tetragrams[tetragram] = frequency
	}
	return tetragrams
}

func score(text []byte) (score float64) {
	textTGs := map[string]float64{}
	for i := 0; i < len(text)-4; i++ {
		slice := string(text[i : i+4])
		if _, ok := textTGs[slice]; ok {
			textTGs[slice]++ // texts containing more common tetragrams are more likely to be English
		} else {
			textTGs[slice] = 1.0
		}
	}

	for key, val := range textTGs {
		if _, ok := tg[key]; ok {
			score += val * tg[key]
		}
	}
	score /= float64((len(text) - 3)) // divide by number of tetragrams
	return score
}

func letters() [25]byte {
	return [25]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'} // remove whichever letter from here
}

func randomise() [25]byte { // used to randomise starting point
	key := letters()
	rand.Shuffle(len(key), func(i, j int) { key[i], key[j] = key[j], key[i] })
	return key
}

func coordinates(gridIndex int) coords {
	var row int = gridIndex / 5
	column := gridIndex - row*5
	return coords{row, column}
}

func gridIndex(c coords) int {
	return 5*c.row + c.column
}

func decipher(ciphertext []byte, key individual) (deciphered []byte) {
	for i := 0; i < len(ciphertext)-1; i += 2 {
		v := ciphertext[i : i+2]
		var i1, i2 int // indices in keyed grids
		for i := 0; i < len(key.b); i++ {
			if key.b[i] == v[0] {
				i1 = i
			}
		}
		for i := 0; i < len(key.b); i++ {
			if key.c[i] == v[1] {
				i2 = i
			}
		}

		coords1 := coordinates(i1)
		coords2 := coordinates(i2)

		di1 := gridIndex(coords{coords1.row, coords2.column}) // indices in non-keyed grids
		di2 := gridIndex(coords{coords2.row, coords1.column})

		deciphered = append(deciphered, key.a[di1])
		deciphered = append(deciphered, key.d[di2])

	}
	return deciphered
}

func mutateGrid(k [25]byte) [25]byte {
	n := k
	i1 := rand.Intn(25)
	i2 := rand.Intn(25)

	n[i1], n[i2] = n[i2], n[i1]

	return n
}

func mutateIndividual(i individual) individual {
	n := i
	x := rand.Intn(2) // random number deciding which grid to mutate
	if x == 0 {
		n.b = mutateGrid(i.b)
	} else {
		n.c = mutateGrid(i.c)
	}
	return n
}

func geneticAlgorithm(maxGens, genSize, k int, ciphertext string) string {
	population := make([]individual, genSize)
	for i := range population {
		n := individual{letters(), randomise(), randomise(), letters(), 0.0}
		population[i] = n
		population[i].score = score(decipher([]byte(ciphertext), n))
	}

	for i := 0; i < maxGens; i++ {
		if i%50 == 0 {
			fmt.Println(i)
		}
		for _, p := range population {
			for j := 0; j < k; j++ {
				newIndividual := mutateIndividual(p)
				newIndividual.score = score(decipher([]byte(ciphertext), newIndividual))
				population = append(population, newIndividual)
			}
		}
		sort.Slice(population, func(i, j int) bool {
			return population[i].score < population[j].score
		})
		population = population[len(population)-genSize-1:]
	}
	return string(decipher([]byte(ciphertext), population[len(population)-1]))
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)

	fmt.Println(geneticAlgorithm(1000, 20, 5, ciphertext))
}
