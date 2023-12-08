package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//!not currently working for ciphers with offset

// !const
var tg map[string]int = tetragrams() // used instead of initialiser function so the massive file is only read once

func tetragrams() map[string]int {
	tetragrams := map[string]int{}
	text, err := os.ReadFile("../monoalphabetic/tetragrams.txt")
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

func score(text string) (score int) {
	for i := 0; i < len(text)-4; i++ {
		slice := text[i : i+4]
		if val, ok := tg[slice]; ok {
			score += val // texts containing more common tetragrams are more likely to be English
		}
	}
	return score
}

func genReference(l int) (ref string) { // generates reference text
	original, err := os.ReadFile("referenceText.txt")
	check(err)
	r := format(original)

	for i := 0; i < l; i++ {
		ref += string(r[i])
	}
	return ref
}

func format(original []byte) (ciphertext string) {
	for i := 0; i < len(original); i++ {
		if 65 <= original[i] && 90 >= original[i] {
			ciphertext += string(original[i])
		}
	}
	return ciphertext
}

func shift(railFence [][]byte, offset int, text string) [][]byte {
	rf := make([][]byte, len(railFence))
	railNum := len(railFence)

	r, d := 0, 1
	for i := 0; i < offset; i++ {
		rf[r] = append(rf[r], '#')
		r += d
		if r == railNum-1 || r == 0 {
			d = -d
		}
	}

	r2, d2 := 0, 1

	for i := 0; i < len(text); i++ {
		rf[r] = append(rf[r], railFence[r2][d2])
		r += d
		if r == railNum-1 || r == 0 {
			d = -d
		}
		r2 += d2
		if r2 == railNum-1 || r2 == 0 {
			d2 = -d2
		}
	}

	return rf
}

func genFence(railNum, offset int, text string) [][]byte {
	railFence := make([][]byte, railNum)

	repeatUnit := 2*railNum - 2 // defined as starting at 0 and ending at 2n - 2
	overflow := len(text) % repeatUnit
	units := (len(text) - overflow) / repeatUnit

	charsPerRow := make([]int, railNum)
	for i := 0; i < railNum; i++ {
		chars := 0
		if i == 0 || i == railNum-1 {
			chars += units
		} else {
			chars += units * 2
		}
		charsPerRow[i] = chars
	}

	d := 1
	r := 0

	for i := 0; i < overflow; i++ {
		charsPerRow[r]++
		r += d
		if r == railNum-1 || r == 0 {
			d = -d
		}
	}

	index := 0

	for row := 0; row < len(charsPerRow); row++ {
		for i := 0; i < charsPerRow[row]; i++ {
			railFence[row] = append(railFence[row], text[index])
			index++
		}
	}

	railFence = shift(railFence, offset, text)
	return railFence
}

func decode(railNum, offset int, text string) (answer string) {
	railFence := genFence(railNum, offset, text)
	for _, r := range railFence {
		for _, b := range r {
			fmt.Print(string(b) + " ")
		}
		fmt.Println()

	}

	indices := make([]int, railNum) // stores current index on each rail

	currentRail, d := 0, 1
	for i := 0; i < offset; i++ {
		indices[currentRail]++
		currentRail += d
		if currentRail == railNum-1 || currentRail == 0 {
			d = -d
		}
	}

	for i := 0; i < len(text); i++ {
		if indices[currentRail] == len(railFence[currentRail]) {
			break
		}
		b := railFence[currentRail][indices[currentRail]]
		if b != '#' {
			answer += string(b)
		}
		indices[currentRail]++
		currentRail += d
		if currentRail == 0 || currentRail == railNum-1 { // last rail
			d = -d
		}
	}
	return answer
}

func bruteForce(ciphertext string, maxRails int) {
	referenceScore := score(genReference(len(ciphertext)))
	for i := 3; i <= maxRails; i++ {
		for j := 0; j < 2*maxRails-2; j++ {
			plaintext := decode(i, j, ciphertext)
			if score(plaintext)*10 >= referenceScore*8 { // within 20% of english
				fmt.Println(plaintext)
			}
		}
	}
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)
	fmt.Println(ciphertext)
	// fmt.Println(decode(5, 0, ciphertext))
	// genFence(3, 3, ciphertext)
	// bruteForce(ciphertext, 20)
}
