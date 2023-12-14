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

func gridCells() []int {
	return []int{11, 12, 13, 14, 15, 21, 22, 23, 24, 25, 31, 32, 33, 34, 35, 41, 42, 43, 44, 45, 51, 52, 53, 54, 55}
}

func letters() []byte {
	return []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y'} // remove J or Z as necessary
}

func format(original []byte, min, max byte) (ciphertext string) {
	for i := 0; i < len(original); i++ {
		if min <= original[i] && max >= original[i] { // polybius only conatins numbers 1-5
			ciphertext += string(original[i])
		}
	}
	return ciphertext
}

func scanDict() (words []string) {
	dict, err := os.Open("../google-10000-english.txt")
	check(err)

	words = append(words, "polybius") // worth a try lol

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
	r := format(original, 65, 90)

	for i := 0; i < l; i++ {
		ref += string(r[i])
	}
	return ref
}

func genKey(keyword string) map[int]byte {
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

	key := map[int]byte{}

	for index, gc := range gridCells() {
		key[gc] = grid[index]
	}
	return key
}

func encipher(text, keyword string) []int {
	key := invert(genKey(keyword))
	var ciphertext []int
	for i := 0; i < len(text); i++ {
		ciphertext = append(ciphertext, key[text[i]])
	}
	return ciphertext
}

func decipher(text string, keyword string) (result string) {
	var nums []int
	for i := 0; i < len(text)-1; i += 2 {
		slice := text[i : i+2]
		n, err := strconv.Atoi(slice)
		check(err)
		nums = append(nums, n)
	}
	key := genKey(keyword)
	for _, num := range nums {
		result += string(key[num])
	}
	return result
}

func polybius() map[byte]int {
	return map[byte]int{
		'A': 11,
		'B': 12,
		'C': 13,
		'D': 14,
		'E': 15,
		'F': 21,
		'G': 22,
		'H': 23,
		'I': 24,
		'J': 25,
		'K': 31,
		'L': 32,
		'M': 33,
		'N': 34,
		'O': 35,
		'P': 41,
		'Q': 42,
		'R': 43,
		'S': 44,
		'T': 45,
		'U': 51,
		'V': 52,
		'W': 53,
		'X': 54,
		'Y': 55,
	}
}

func invert[K comparable, V comparable](m map[K]V) map[V]K {
	r := make(map[V]K, len(m))
	for k, v := range m {
		r[v] = k
	}
	return r
}

func convertToPolybius(s string) (p []int) {
	for i := 0; i < len(s); i++ {
		p = append(p, polybius()[s[i]])
	}
	return p
}

func convertFromPolybius(nums []int) (res string) {
	pToEng := invert(polybius())
	for _, num := range nums {
		res += string(pToEng[num])
	}
	return res
}

func dictionaryAttack(ciphertext string) {
	dict := scanDict()
	referenceScore := score(genReference(len(ciphertext) / 2)) // divided by 2 because 2 digits correspond to one number
	fmt.Println(referenceScore)
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

	ciphertext := format(original, 49, 53)
	dictionaryAttack(ciphertext)
}
