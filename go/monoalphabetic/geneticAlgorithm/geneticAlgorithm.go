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
	key   [26]byte
	score float64
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
	text, err := os.ReadFile("../../logTetragrams.txt")
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

func encipher(key [26]byte, plaintext string) (ciphertext []byte) {
	// also used as deciphering function with inverse key
	for i := 0; i < len(plaintext); i++ {
		for j := 0; j < len(letters()); j++ {
			if letters()[j] == plaintext[i] {
				ciphertext = append(ciphertext, key[j])
			}
		}
	}
	return ciphertext
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

func letters() [26]byte {
	return [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
}

func randomise() [26]byte { // used to randomise starting point
	key := letters()
	rand.Shuffle(len(key), func(i, j int) { key[i], key[j] = key[j], key[i] })
	return key
}

func mutate(k [26]byte) [26]byte {
	n := k
	i1 := rand.Intn(26)
	i2 := rand.Intn(26)

	n[i1], n[i2] = n[i2], n[i1]

	return n
}

func geneticAlgorithm(maxGens, genSize, k int, ciphertext string) string {
	population := make([]individual, genSize)
	for i := range population { // generate initial population
		n := randomise()
		population[i].key = n
		population[i].score = score(encipher(n, ciphertext))
	}

	for i := 0; i < maxGens; i++ {
		for _, p := range population {
			for j := 0; j < k; j++ {
				newKey := mutate(p.key)
				newScore := score(encipher(newKey, ciphertext))
				population = append(population, individual{newKey, newScore})
			}
		}
		sort.Slice(population, func(i, j int) bool {
			return population[i].score < population[j].score
		})
		population = population[len(population)-genSize-1:]
	}

	return string(encipher(population[len(population)-1].key, ciphertext))
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)

	fmt.Println(geneticAlgorithm(500, 20, 5, ciphertext))
}
