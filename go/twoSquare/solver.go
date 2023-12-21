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

func fillGridVH(keyword1, keyword2 string) []byte { // fills two vertical grids horizontally
	unused := letters()
	var grid1, grid2, grid []byte
	for i := 0; i < len(keyword1); i++ {
		l := keyword1[i]
		for j := 0; j < len(unused); j++ {
			if unused[j] == l {
				grid1 = append(grid1, l) // only append if it isnt already used
				unused = append(unused[:j], unused[j+1:]...)
			}
		}
	}
	for i := 0; i < len(unused); i++ {
		grid1 = append(grid1, unused[i])
	}
	unused = letters()
	for i := 0; i < len(keyword2); i++ {
		l := keyword2[i]
		for j := 0; j < len(unused); j++ {
			if unused[j] == l {
				grid2 = append(grid2, l) // only append if it isnt already used
				unused = append(unused[:j], unused[j+1:]...)
			}
		}
	}
	for i := 0; i < len(unused); i++ {
		grid2 = append(grid2, unused[i])
	}
	grid = append(grid1, grid2...)
	return grid
}

func gridIndexV(row, column int) int {
	return row*5 + column
}

func rowAndColumnV(grid int) (int, int) {
	row := grid / 5
	column := grid - row*5
	return row, column
}

func gridIndexH(row, column int) int {
	return row*10 + column
}

func rowAndColumnH(grid int) (int, int) {
	row := grid / 10
	column := grid - row*10
	return row, column
}

func decipherVH(text, keyword1, keyword2 string) (deciphered string) {
	grid := fillGridVH(keyword1, keyword2)
	for i := 0; i < len(text)-1; i += 2 {
		var idx1, idx2 int
		for j := 0; j < len(grid)/2; j++ {
			if text[i] == grid[j] {
				idx1 = j
			}
		}
		for k := len(grid) / 2; k < len(grid); k++ {
			if text[i+1] == grid[k] {
				idx2 = k
			}
		}
		if idx1 == idx2 {
			panic("Indexes are the same")
		}
		row1, column1 := rowAndColumnV(idx1)
		row2, column2 := rowAndColumnV(idx2)

		var dRow1, dColumn1, dRow2, dColumn2 int

		if row1 == row2 && column1 == column2 {
			panic("Row and Column are the same")
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

		letter1 := grid[gridIndexV(dRow1, dColumn1)]
		letter2 := grid[gridIndexV(dRow2, dColumn2)]
		deciphered += string(letter1)
		deciphered += string(letter2)

	}
	return deciphered
}

func dictionaryAttack(ciphertext string) {
	dict := scanDict()[:1000]
	referenceScore := score(genReference(len(ciphertext)))
	for i := 0; i < len(dict); i++ {
		fmt.Println(i)
		for j := 0; j < len(dict); j++ {
			p := decipherVH(ciphertext, dict[i], dict[j])
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
}
