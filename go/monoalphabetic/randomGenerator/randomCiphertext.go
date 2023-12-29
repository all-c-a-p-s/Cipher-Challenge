package main

import (
	"math/rand"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func uppercase(x byte) byte {
	if x >= 97 && x <= 122 {
		return x - 32
	}
	panic("Attempt to convert non uppercase letter to lowercase")
}

func letters() [26]byte {
	return [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
}

func format(text []byte) (letters []byte) {
	for i := 0; i < len(text); i++ {
		if 65 <= text[i] && text[i] <= 90 {
			letters = append(letters, text[i])
		} else if 97 <= text[i] && text[i] <= 122 {
			letters = append(letters, uppercase(text[i]))
		}
	}
	return letters
}

func encipher(key [26]byte, plaintext []byte) (ciphertext []byte) {
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

func randomise() [26]byte { // used to randomise starting point
	key := letters()
	rand.Shuffle(len(key), func(i, j int) { key[i], key[j] = key[j], key[i] })
	return key
}

func main() {
	original, err := os.ReadFile("plaintext.txt")
	check(err)
	plaintext := format(original)
	ciphertext := encipher(randomise(), plaintext)
	file, err := os.Create("randomCiphertext.txt")
	check(err)
	defer file.Close()
	file.WriteString(string(ciphertext))
}
