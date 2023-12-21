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

func gridIndex(row, column int) int {
	return row*5 + column
}

func rowAndColumn(grid int) (int, int) {
	row := grid / 5
	column := grid - row*5
	return row, column
}

func decipher(text, keyword string) (deciphered string) {
	grid := fillGrid(keyword)
	for i := 0; i < len(text)-1; i += 2 {
		var idx1, idx2 int
		for j := 0; j < len(grid); j++ {
			if text[i] == grid[j] {
				idx1 = j
			} else if text[i+1] == grid[j] {
				idx2 = j
			}
		}
		row1, column1 := rowAndColumn(idx1)
		row2, column2 := rowAndColumn(idx2)

		var dRow1, dColumn1, dRow2, dColumn2 int

		if row1 == row2 && column1 == column2 {
			fmt.Println(i)
			fmt.Println(string(text[len(deciphered)]), string(text[len(deciphered)+1]))
			panic("Two letters on same grid square")
		} else if row1 == row2 {
			dRow1 = row1
			dRow2 = row2
			dColumn1 = (column1 - 1 + 5) % 5 // adding 5 bc negative mod doesn't work
			dColumn2 = (column2 - 1 + 5) % 5
		} else if column1 == column2 {
			dColumn1 = column1
			dColumn2 = column2
			dRow1 = (row1 - 1 + 5) % 5
			dRow2 = (row2 - 1 + 5) % 5
		} else { // no same
			dRow1 = row1
			dRow2 = row2
			dColumn1 = column2
			dColumn2 = column1
		}

		letter1 := grid[gridIndex(dRow1, dColumn1)]
		letter2 := grid[gridIndex(dRow2, dColumn2)]
		deciphered += string(letter1)
		deciphered += string(letter2)

	}
	return deciphered
}

func dictionaryAttack(ciphertext string) {
	dict := scanDict()
	referenceScore := score(genReference(len(ciphertext)))
	for i := 0; i < len(dict); i++ {
		p := decipher(ciphertext, dict[i])
		s := score(p)
		if 10*s >= referenceScore*8 {
			fmt.Println(p)
		}
	}
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)

	ciphertext := format(original)
	dictionaryAttack(ciphertext)
}
