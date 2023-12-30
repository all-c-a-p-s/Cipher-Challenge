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
	"sync"
)

type individual struct {
	key   [26]byte
	score float64
	// individual made up of key, and score for purpose of sorting
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func format(original []byte) (ciphertext string) {
	// formats text to remove character that are not capital letters
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
	// scans tetragram file into a hashmap
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
	// score based on tetragram frequencies
	textTGs := map[string]float64{}
	for i := 0; i < len(text)-3; i++ {
		// read all tetragrams in text, including overlaps
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
			// lookup in global hashmap
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

func crossover(k1 [26]byte, k2 [26]byte) [26]byte {
	// crossover function to combine two keys
	nums := [26]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}
	rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	// shuffle nums so that which letters come from which key is randomised
	unused := map[byte]bool{} // hashset
	for i := 0; i < len(letters()); i++ {
		unused[letters()[i]] = true
	} // initialise unused with all letters
	var child [26]byte
	for i := 0; i < 13; i++ {
		child[nums[i]] = k1[nums[i]]
	}
	for i := 13; i < 26; i++ {
		child[nums[i]] = k2[nums[i]]
	} // fill new key with 13 from each
	duplicates := []int{}           // hashset of duplicate indices
	lettersFound := map[byte]bool{} // hashset of letters already found in array
	for i := 0; i < len(child); i++ {
		if _, ok := lettersFound[child[i]]; ok {
			// if letter already found
			duplicates = append(duplicates, i)
		} else {
			delete(unused, child[i])
			// if found for the first time, remove from unsed hashset
		}
		lettersFound[child[i]] = true
	}
	var unusedLetters []byte
	for k := range unused {
		unusedLetters = append(unusedLetters, k)
	} // convert hashset into slice

	for len(unusedLetters) != 0 { // Go docs say that hashmap iteration isn't necessarily properly random, just unpredictable
		var r int // random number choosing which unused letter should be inserted
		if len(unusedLetters) != 1 {
			r = rand.Intn(len(unusedLetters) - 1)
		} else {
			r = 0
		}
		child[duplicates[0]] = unusedLetters[r]                           // fill duplicate index with unused letter
		unusedLetters = append(unusedLetters[:r], unusedLetters[r+1:]...) // remove unused letter from slice
		if len(duplicates) != 1 {
			duplicates = duplicates[1:]
		} else {
			break
		}
	}
	return child
}

func mutate(k [26]byte) [26]byte {
	// basic mutation function which swaps two indices in a key
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

	for gen := 0; gen < maxGens; gen++ {
		// go through generations
		children := make([]individual, genSize*k)
		var wg sync.WaitGroup
		// sync.WaitGroup to wait until all goroutines are done
		for i := 0; i < genSize; i++ {
			currentParent := i
			// necessary because otherwise i is "captured by" goroutines
			wg.Add(1)
			go func() {
				defer wg.Done()
				r1 := rand.Intn(genSize)
				r2 := rand.Intn(genSize)
				for r2 == r1 {
					r2 = rand.Intn(genSize)
				}
				parent1 := population[r1]
				parent2 := population[r2]
				for j := 0; j < k; j++ {
					newKey := mutate(crossover(parent1.key, parent2.key)) // crossover plus random mutation
					newScore := score(encipher(newKey, ciphertext))
					children[currentParent*k+j] = individual{newKey, newScore}
				}
			}() // concurrent call to generate children
		}
		wg.Wait()
		// wait until all concurrent calls are done before adding to population
		population = append(population, children...)
		sort.Slice(population, func(i, j int) bool {
			return population[i].score < population[j].score
		})
		population = population[len(population)-genSize:]
	}

	fmt.Print("Final Score: ")
	fmt.Println(score(encipher(population[len(population)-1].key, ciphertext)))
	return string(encipher(population[len(population)-1].key, ciphertext))
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)
	fmt.Println(geneticAlgorithm(40, 500, 5, ciphertext))
}
