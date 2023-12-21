package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

//!Not completely reliable, tends to work around 1 in 3 times on ciphertexts ~500 chars

type key struct {
	k     [26]byte
	score int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func format(original []byte) (ciphertext string) {
	for i := 0; i < len(original); i++ {
		if 65 <= original[i] && 90 >= original[i] {
			ciphertext += string(original[i])
		}
	}
	return ciphertext
}

// !const
var tg map[string]int = tetragrams() // used instead of initialiser function so the massive file is only read once

func tetragrams() map[string]int {
	tetragrams := map[string]int{}
	text, err := os.ReadFile("../tetragrams.txt")
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

func encipher(key [26]byte, plaintext string) (ciphertext string) {
	for i := 0; i < len(plaintext); i++ {
		for j := 0; j < 26; j++ {
			if letters()[j] == plaintext[i] {
				ciphertext += string(key[j])
			}
		}
	}
	return ciphertext
}

func letters() [26]byte {
	return [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
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

func randomise() [26]byte { // used to randomise starting point
	key := letters()
	rand.Shuffle(len(key), func(i, j int) { key[i], key[j] = key[j], key[i] })
	return key
}

func mutate(parent [26]byte) [26]byte {
	newKey := parent
	i1 := rand.Intn(26)
	i2 := rand.Intn(26)
	for i2 == i1 {
		i2 = rand.Intn(26) // keep generating random numbers until not equal to i1
	}

	newKey[i1], newKey[i2] = newKey[i2], newKey[i1]

	return newKey
}

func acceptanceProbability(dE, temp float64) float64 {
	return dE * temp
}

func simulatedAnnealing(maxConstant int, maxTemp int, ciphertext string, k float64) {
	alpha := key{randomise(), 0}
	constant := 0
	temp := float64(maxTemp)
	for constant < maxConstant {
		n := mutate(alpha.k)

		p1 := encipher(alpha.k, ciphertext)
		p2 := encipher(n, ciphertext)

		s1 := 1
		s2 := float64(score(p2)) / float64(score(p1))
		deltaE := float64(s2 / float64(s1))
		if deltaE >= 1 {
			alpha.k = n
			alpha.score = score(p2)
			constant = 0
		} else {
			p := acceptanceProbability(deltaE, temp)
			x := rand.Float64()
			if x <= p {
				alpha.k = n
				alpha.score = score(p2)
				constant = 0
			} else {
				constant++
			}
		}
		temp *= k
	}
	fmt.Println(encipher(alpha.k, ciphertext))
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)
	simulatedAnnealing(500, 1, ciphertext, 0.95)
}
