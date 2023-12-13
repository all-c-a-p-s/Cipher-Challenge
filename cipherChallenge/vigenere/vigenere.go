package main

import (
	"fmt"
	"log"
	"os"
)

func removePunctuation(original string) (ciphertext string) {
	for i := 0; i < len(original); i++ {
		for j := 0; j < len(letters()); j++ {
			if original[i] == letters()[j] {
				ciphertext += string(original[i])
			}
		}
	}
	return ciphertext
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func letterFrequencies() map[byte]float32 { // initialiser function
	freq := map[byte]float32{
		'A': 8.12,
		'B': 1.49,
		'C': 2.71,
		'D': 4.32,
		'E': 12.02,
		'F': 2.30,
		'G': 2.03,
		'H': 5.92,
		'I': 7.31,
		'J': 0.1,
		'K': 0.69,
		'L': 3.98,
		'M': 2.61,
		'N': 6.95,
		'O': 7.68,
		'P': 1.82,
		'Q': 0.11,
		'R': 6.02,
		'S': 6.28,
		'T': 9.10,
		'U': 2.88,
		'V': 1.11,
		'W': 2.09,
		'X': 0.17,
		'Y': 2.11,
		'Z': 0.07,
	}

	return freq
}

func letters() []byte { // initialiser function
	letters := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

	return letters
}

func keyLength(ciphertext string) int {
	var keyLength int
	var coincidences []int
	var total int

	for i := 1; i < len(ciphertext)/2; i++ {
		c := 0 // coincidences on current shift
		for j := 0; j < len(ciphertext)-i; j++ {
			if ciphertext[j+i] == ciphertext[j] {
				c++
			}
		}
		coincidences = append(coincidences, c)
		total += c
	}

	var highIndices []int
	mean := total / len(coincidences)

	for i := 0; i < len(coincidences); i++ {
		if float32(coincidences[i]) > float32(mean)*1.5 {
			highIndices = append(highIndices, i)
		}
	}

	var bestLength int

	for i := 1; i < len(highIndices); i++ { // where i is increment between steps
		currentTotal := 0
		for j := 0; j < len(highIndices); j++ {
			if j == len(highIndices)-1 {
				break
			}
			if highIndices[j]+i == highIndices[j+1] { // uses modulo to check for multiples, doesn't matter if this is negative
				currentTotal++
			} else {
				currentTotal-- // used so that factors of actual key don't get equal or higher total
			}
		}
		if currentTotal > bestLength {
			keyLength = i
			bestLength = currentTotal
		}
	}

	return keyLength
}

func findKey(length int, ciphertext string) (key []int) {
	for i := 0; i < length; i++ { // loop through characters shifted by various amounts
		bestShift := 0
		var bestScore float32
		for shift := 0; shift < 26; shift++ {
			frequenciesFound := map[byte]int{
				'A': 0,
				'B': 0,
				'C': 0,
				'D': 0,
				'E': 0,
				'F': 0,
				'G': 0,
				'H': 0,
				'I': 0,
				'J': 0,
				'K': 0,
				'L': 0,
				'M': 0,
				'N': 0,
				'O': 0,
				'P': 0,
				'Q': 0,
				'R': 0,
				'S': 0,
				'T': 0,
				'U': 0,
				'V': 0,
				'W': 0,
				'X': 0,
				'Y': 0,
				'Z': 0,
			}
			for j := i; j < len(ciphertext); j += length {
				frequenciesFound[string(rune((int(ciphertext[j])-65+shift)%26) + 65)[0]]++ // convert to 0-25, then add shift and % 26, than add 65 back, then convert to byte by string[0]
			}
			shiftScore := frequencyScore(frequenciesFound)
			if shiftScore > bestScore {
				bestScore = shiftScore
				bestShift = shift
			}

		}

		key = append(key, (26-bestShift)%26)
	}

	return key
}

func frequencyScore(frequencies map[byte]int) float32 {
	freq := letterFrequencies()
	letters := letters()

	var total float32

	for _, letter := range letters {
		total += float32(frequencies[letter]) * freq[letter]
	}

	return total
}

func decrypt(ciphertext string, key []int) (decrypted string) {
	step := 0
	for i := 0; i < len(ciphertext); i++ {
		cipherLetterIndex := int(ciphertext[i]) - 65
		letterIndex := (cipherLetterIndex + 26 - key[step%len(key)]) % 26
		decrypted += string(letters()[letterIndex])
		step++
	}
	return decrypted
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	ciphertext := removePunctuation(string(original))
	check(err)
	key := findKey(keyLength(string(ciphertext)), string(ciphertext))

	fmt.Println(decrypt(ciphertext, key))
}
