package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type matrix struct {
	a int // top left
	b int // top right
	c int // bottom left
	d int // bottom right
}

type vector struct {
	x int
	y int
}

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

func genReference(l int) (ref string) { // generates reference text
	original, err := os.ReadFile("../railfence/referenceText.txt")
	check(err)
	r := format(original)

	for i := 0; i < l; i++ {
		ref += string(r[i])
	}
	return ref
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

func format(original []byte) (ciphertext string) {
	for i := 0; i < len(original); i++ {
		if 65 <= original[i] && 90 >= original[i] {
			ciphertext += string(original[i])
		}
	}
	return ciphertext
}

func multiply(v vector, m matrix) vector {
	return vector{m.a*v.x + m.b*v.y, m.c*v.x + m.d*v.y}
}

func modulo(v vector) vector { // mod of negative returns negative, which isn't helpful here
	x, y := v.x, v.y

	for x < 0 {
		x += 26
	}
	for y < 0 {
		y += 26
	}

	return vector{x % 26, y % 26}
}

func convertToVectors(s string) (vecs []vector) {
	for i := 0; i < len(s)-1; i += 2 {
		slice := s[i : i+2]
		vecs = append(vecs, vector{int(slice[0] - 65), int(slice[1] - 65)})
	}
	return vecs
}

func convertToString(vecs []vector) (s string) {
	for i := 0; i < len(vecs); i++ {
		s += string(vecs[i].x + 65)
		s += string(vecs[i].y + 65)
	}
	return s
}

func inverse(m matrix) matrix {
	return matrix{m.d, -m.b, -m.c, m.a}
}

func decipher(ciphertext string, m matrix) string {
	vecs := convertToVectors(ciphertext)
	var shifted []vector
	for _, v := range vecs {
		shifted = append(shifted, modulo(multiply(v, m)))
	}
	// fmt.Println(shifted)
	s := convertToString(shifted)
	return s
}

func bruteForce(ciphertext string) {
	referenceScore := score(genReference(len(ciphertext)))
	for i := 0; i <= 9999; i++ { // where each decimal bit represent one number in the matrix
		a := i / 1000                    // integer divide
		b := (i - a*1000) / 100          // integer divide, while subtracting thousands
		c := (i - a*1000 - b*100) / 10   // integer divide while subtracting hundreds as well
		d := (i - a*1000 - b*100 - c*10) // units digit

		m := matrix{a, b, c, d}

		plaintext := decipher(ciphertext, inverse(m))
		if score(plaintext)*10 >= referenceScore*8 {
			fmt.Println(m)
			fmt.Println(plaintext)
		}
	}
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)
	bruteForce(ciphertext)
}
