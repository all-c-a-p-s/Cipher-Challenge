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

func letters() []byte {
	return []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'} // remove J or Z as necessary
}

func format(original []byte) (result string) {
	for i := 0; i < len(original); i++ {
		if 65 <= original[i] && original[i] <= 90 {
			result += string(original[i])
		}
	}
	return result
}

func scanDict() (words []string) {
	dict, err := os.Open("../google-10000-english.txt")
	check(err)

	words = append(words, "POLYBIUS") // worth a try lol
	words = append(words, "PLAYFAIR")

	scanner := bufio.NewScanner(dict)
	for scanner.Scan() {
		words = append(words, strings.ToUpper(scanner.Text()))
	}
	return words
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
	original, err := os.ReadFile("../railfence/referenceText.txt")
	check(err)
	r := format(original)

	for i := 0; i < l; i++ {
		ref += string(r[i])
	}
	return ref
}

func fillGrid(keyword string) []byte {
	unused := letters()
	var grid []byte
	for i := 0; i < len(keyword); i++ {
		l := keyword[i]
		for j := 0; j < len(unused); j++ {
			if unused[j] == l {
				grid = append(grid, l) // only append if it isnt already used
				unused = append(unused[:j], unused[j+1:]...)
			}
		}
	}
	for i := 0; i < len(unused); i++ {
		grid = append(grid, unused[i])
	}
	return grid
}

func coordinates(idx int) (row, column int) {
	row = idx / 5
	column = idx - row*5
	return row, column
}

func gridIndex(r, c int) int {
	return 5*r + c
}

func decipher(text, keyword string, n int) (deciphered string) {
	grid := fillGrid(keyword)
	numBoxes := 1
	if n != 0 {
		numBoxes = len(text) / n // where a box is 3 rows of n digits
		if len(text)%n != 0 {
			numBoxes++
		}
	}
	boxes := make([]string, numBoxes)
	var idx int
	var currentBox int
	for i := 0; i < len(text); i++ {

		for j := 0; j < len(grid); j++ {
			if grid[j] == text[i] {
				idx = j
			}
		}
		r, c := coordinates(idx)
		boxes[currentBox] += fmt.Sprintf("%d%d", r, c) // len of box is always a mutliple of 3 so this is fine
		if len(boxes[currentBox]) == (2 * n) {         // current box filled
			currentBox++
		}
	}

	var row1, row2 string
	for _, box := range boxes {
		row1 += box[:len(box)/2]
		row2 += box[len(box)/2:]
	}

	var unboxed []string // 10/10 variable name
	for i := 0; i < len(row1); i++ {
		concat := string(row1[i]) + string(row2[i])
		unboxed = append(unboxed, concat)
	}

	for _, u := range unboxed {
		r, err := strconv.Atoi(string(u[0]))
		check(err)
		c, err := strconv.Atoi(string(u[1]))
		check(err)
		deciphered += string(grid[gridIndex(r, c)])
	}
	return deciphered
}

func dictionaryAttack(ciphertext string) {
	dict := scanDict()[:1000]
	referenceScore := score(genReference(len(ciphertext)))
	for n := 0; n < 20; n++ {
		for i := 0; i < len(dict); i++ {
			p := decipher(ciphertext, dict[i], n)
			s := score(p)
			if 10*s >= referenceScore*8 {
				fmt.Println(p)
			}
		}
	}
}

func wordListAttack(ciphertext string) { // words list generated from previous plaintexts of the challenge
	w, err := os.ReadFile("../challenge_word_list.txt")
	check(err)
	f := strings.ToUpper(string(w))
	words := strings.Fields(f)
	referenceScore := score(genReference(len(ciphertext)))
	for n := 0; n < 20; n++ {
		for i := 0; i < len(words); i++ {
			p := decipher(ciphertext, words[i], n)
			s := score(p)
			if 10*s >= referenceScore*8 {
				fmt.Println(p)
			}
		}
	}
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)
	dictionaryAttack(ciphertext)
	wordListAttack(ciphertext)
}
