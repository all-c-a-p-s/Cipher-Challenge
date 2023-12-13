package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

type key struct {
	score int
	k     map[byte]byte
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

func sample() string {
	sample, err := os.ReadFile("sample.txt")
	check(err)
	return string(sample)
}

// !const
var (
	tg          map[string]int = tetragrams() // used instead of initialiser function so the massive file is only read once
	sampleText  string         = sample()
	sampleScore                = scoreSample()
)

func scoreSample() (score int) {
	for i := 0; i < len(sampleText)-4; i++ {
		slice := sampleText[i : i+4]
		if val, ok := tg[slice]; ok {
			score += val // texts containing more common tetragrams are more likely to be English
		}
	}

	return score
}

func tetragrams() map[string]int {
	tetragrams := map[string]int{}
	text, err := os.ReadFile("tetragrams.txt")
	check(err)

	scanner := bufio.NewScanner(strings.NewReader(string(text)))

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ", ")
		tetragram := fields[0]
		frequency, err := strconv.Atoi(fields[1])
		check(err)

		tetragrams[tetragram] = frequency
	}
	return tetragrams
}

func encipher(key map[byte]byte, plaintext string) (ciphertext string) {
	for i := 0; i < len(plaintext); i++ {
		ciphertext += string(key[plaintext[i]])
	}
	return ciphertext
}

func score(text string) (score int) {
	for i := 0; i < len(text)-4; i++ {
		slice := text[i : i+4]
		if val, ok := tg[slice]; ok {
			score += val // texts containing more common tetragrams are more likely to be English
		}
	}
	return score
}

func letters() []byte {
	return []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
}

func randomise() map[byte]byte { // used to randomise starting point
	key := map[byte]byte{}
	remaining := letters()
	for i := 0; i < 25; i++ {
		index := rand.Intn(len(remaining))
		key[remaining[index]] = letters()[i]
		copy(remaining[index:], remaining[index+1:])
		remaining[len(remaining)-1] = '/'
		remaining = remaining[:len(remaining)-1] // remove element from slice
	}
	key[remaining[0]] = letters()[25] // assign last letter in key

	return key
}

func mutate(key1 *key, ciphertext string) {
	i1 := rand.Intn(26)
	i2 := rand.Intn(26)

	(*key1).k[letters()[i1]], (*key1).k[letters()[i2]] = (*key1).k[letters()[i2]], (*key1).k[letters()[i1]]

	for {
		r := rand.Intn(100)
		if r >= 95 {
			i1 := rand.Intn(26)
			i2 := rand.Intn(26)

			(*key1).k[letters()[i1]], (*key1).k[letters()[i2]] = (*key1).k[letters()[i2]], (*key1).k[letters()[i1]]
		} else {
			break
		}

	}

	plaintext := encipher((*key1).k, ciphertext)
	(*key1).score = score(plaintext)
}

func generation(keys *[]key, newGen int, genSize int, ciphertext string) {
	for i := 0; i < genSize; i++ {
		for i := 0; i < newGen; i++ {
			child := key{0, (*keys)[i].k}
			mutate(&child, ciphertext)
			(*keys) = append(*keys, child)
		}
	}

	slices.SortFunc(*keys, func(i, j key) int {
		return -cmp.Compare(i.score, j.score)
	})

	//(*keys) = (*keys)[len(*keys)-genSize-1:]
	(*keys) = (*keys)[:genSize]
	// fmt.Println((*keys)[0].score)
}

func geneticAlgorithm(genSize int, newGen int, maxGenerations int, ciphertext string) {
	gensElapsed := 0
	var keys []key
	for i := 0; i < genSize; i++ {
		newKey := key{0, randomise()}
		plaintext := encipher(newKey.k, ciphertext)
		newKey.score = score(plaintext)
		keys = append(keys, newKey)
	}
	for {
		// fmt.Println(gensElapsed)
		if gensElapsed == maxGenerations {
			break
		}

		generation(&keys, newGen, genSize, ciphertext)
		gensElapsed++
	}

	// keys will already be sorted by selectFittest
	fmt.Println((keys)[0].score)
	fmt.Println(encipher(keys[0].k, ciphertext))
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)

	geneticAlgorithm(20, 5, 500, ciphertext)
}
