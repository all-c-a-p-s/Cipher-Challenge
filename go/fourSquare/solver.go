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
	"time"
)

type individual struct {
	// individual in the population
	a     [25]byte
	b     [25]byte
	c     [25]byte
	d     [25]byte
	score float64
	// where a, d are the unkeyed grids, b and c are the keyed grids
}

type coords struct {
	// coordinates in the grid
	row    int
	column int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func format(original []byte) (ciphertext string) {
	// formats ciphertext to remove spaces/characters that aren't capital letters
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
var tg map[string]float64 = tetragrams() // used instead of initialiser function so the massive file is only read once

func tetragrams() map[string]float64 {
	// function to scan tetragram file into a hashmap of format key: tetragram, value: score
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
	// function to return fitness score based on the deciphered text
	textTGs := map[string]float64{} // hashmap containing frequencies of tetragrams in the text
	for i := 0; i < len(text)-3; i++ {
		// loop through text adding tetragrams to hashmap
		slice := string(text[i : i+4])
		if _, ok := textTGs[slice]; ok {
			textTGs[slice]++ // texts containing more common tetragrams are more likely to be English
		} else {
			textTGs[slice] = 1.0
		}
	}

	for key, val := range textTGs {
		if _, ok := tg[key]; ok {
			// lookup value of tetragram in global hashmap tg
			score += val * tg[key]
		}
	}
	score /= float64((len(text) - 3)) // divide by number of tetragrams
	return score
}

func letters() [25]byte { // initialiser function containing letters except for J
	return [25]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'} // remove whichever letter from here
}

func randomise() [25]byte {
	// shuffles a copy of the letters() array
	// used to generate initial random population
	key := letters()
	rand.Shuffle(len(key), func(i, j int) { key[i], key[j] = key[j], key[i] })
	return key
}

func coordinates(gridIndex int) coords {
	// returns coordinates based on grid index
	var row int = gridIndex / 5
	column := gridIndex - row*5
	return coords{row, column}
}

func gridIndex(c coords) int {
	// grid index based on coordinates
	return 5*c.row + c.column
}

func decipher(ciphertext []byte, key individual) (deciphered []byte) {
	// function to decipher ciphertext based on key
	for i := 0; i < len(ciphertext)-1; i += 2 {
		// scans text in blocks of two letters
		bigram := ciphertext[i : i+2]
		var i1, i2 int // indices in keyed grids
		for i := 0; i < len(key.b); i++ {
			if key.b[i] == bigram[0] {
				i1 = i
			}
		}
		for i := 0; i < len(key.b); i++ {
			if key.c[i] == bigram[1] {
				i2 = i
			}
		}

		coords1 := coordinates(i1)
		coords2 := coordinates(i2) // convert indices to coordinates

		di1 := gridIndex(coords{coords1.row, coords2.column}) // indices in non-keyed grids
		di2 := gridIndex(coords{coords2.row, coords1.column})

		deciphered = append(deciphered, key.a[di1])
		deciphered = append(deciphered, key.d[di2]) // append new bigram to decipher text

	}
	return deciphered
}

func mutateGrid(k [25]byte) [25]byte {
	// mutates a grid by randomly swapping two indices in the grid
	n := k
	i1 := rand.Intn(25)
	i2 := rand.Intn(25)

	n[i1], n[i2] = n[i2], n[i1]

	return n
}

func mutateIndividual(i individual) individual {
	// mutates an individual
	n := i
	x := rand.Intn(2) // random number deciding which grid to mutate
	if x == 0 {
		n.b = mutateGrid(i.b)
	} else {
		n.c = mutateGrid(i.c)
	}
	return n
}

func crossoverGrids(k1 [25]byte, k2 [25]byte) [25]byte {
	// crossover function to combine two grids
	nums := [25]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	// shuffle nums so that which letters come from which key is randomised
	unused := map[byte]bool{} // hashset
	for i := 0; i < len(letters()); i++ {
		unused[letters()[i]] = true
	}
	// initialise hashset with all letters unused
	var child [25]byte
	for i := 0; i < 13; i++ {
		// go through grid 1
		child[nums[i]] = k1[nums[i]]
		// delete element from unused hashmap after that letter is used
	}
	for i := 13; i < 25; i++ {
		child[nums[i]] = k2[nums[i]]
		// same for grid 2
	}
	duplicates := []int{}           // hashset of duplicate indices
	lettersFound := map[byte]bool{} // hashset of letters already found in array
	for i := 0; i < len(child); i++ {
		if _, ok := lettersFound[child[i]]; ok {
			// if letter already found
			duplicates = append(duplicates, i)
			// append index of cuplicate to slice
		} else {
			delete(unused, child[i])
			// delete letters that are used from unused hashmap
		}
		lettersFound[child[i]] = true
		// ^ either adds value to hashset or "updates" it from true to true
	}
	var unusedLetters []byte
	for k := range unused {
		unusedLetters = append(unusedLetters, k)
		// convert hashset into slice
	}

	for len(unusedLetters) != 0 { // Go docs say that hashmap iteration isn't necessarily properly random, just unpredictable
		var r int // random number choosing which unused letter should be inserted
		if len(unusedLetters) != 1 {
			r = rand.Intn(len(unusedLetters) - 1)
		} else {
			r = 0
		}
		child[duplicates[0]] = unusedLetters[r]                           // replace duplicate letter with unsed letter
		unusedLetters = append(unusedLetters[:r], unusedLetters[r+1:]...) // remove eleent from slice
		if len(duplicates) != 1 {
			duplicates = duplicates[1:]
		} else {
			break
		}
	}
	return child
}

func crossover(parent1, parent2 individual) individual {
	// calls crossoverGrids() on both grids
	n := parent1
	n.b = crossoverGrids(parent1.b, parent2.b)
	n.c = crossoverGrids(parent1.c, parent2.c)
	return n
}

func geneticAlgorithm(maxGens, genSize, k int, ciphertext string) string {
	population := make([]individual, genSize)
	for i := range population {
		// randomly generate intial population
		n := individual{letters(), randomise(), randomise(), letters(), 0.0}
		population[i] = n
		population[i].score = score(decipher([]byte(ciphertext), n))
	}

	for gen := 0; gen < maxGens; gen++ {
		// go through generations
		children := make([]individual, genSize*k)
		var wg sync.WaitGroup
		// sync.waitgroup to wait until all goroutines are done
		for i := 0; i < genSize; i++ {
			currentParent := i // necessary because otherise i is "captured by" goroutine
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
					newIndividual := mutateIndividual(crossover(parent1, parent2)) // mixture of crossover and mutation
					newIndividual.score = score(decipher([]byte(ciphertext), newIndividual))
					children[currentParent*k+j] = newIndividual
					// slices are inhernetly references so this works
				}
			}() // concurrently generating children -> faster
		}
		wg.Wait() // wait until all concurrent calls are done before adding children to population
		population = append(population, children...)
		sort.Slice(population, func(i, j int) bool {
			return population[i].score < population[j].score
		}) // sort population by fitness score
		population = population[len(population)-genSize:] // take the genSize fittest in the population
		fmt.Println(score(decipher([]byte(ciphertext), population[len(population)-1])))
	}
	fmt.Print("Final Score: ")
	fmt.Println(score(decipher([]byte(ciphertext), population[len(population)-1])))
	return string(decipher([]byte(ciphertext), population[len(population)-1]))
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)
	start := time.Now()
	fmt.Println(geneticAlgorithm(100, 1000, 5, ciphertext))
	end := time.Now()
	runtime := end.Sub(start)
	fmt.Printf("Done in %v\n", runtime)
}
