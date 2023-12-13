package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

type key struct {
	k     []int
	score int
}

// Tool used to help solve columnar cipher, user still needs to find key based on key-length decoding
var (
	permutations [][]int
	tg           map[string]int = tetragrams()
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func primes() []int { // initialiser function returning primes up to 1000
	return []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541, 547, 557, 563, 569, 571, 577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643, 647, 653, 659, 661, 673, 677, 683, 691, 701, 709, 719, 727, 733, 739, 743, 751, 757, 761, 769, 773, 787, 797, 809, 811, 821, 823, 827, 829, 839, 853, 857, 859, 863, 877, 881, 883, 887, 907, 911, 919, 929, 937, 941, 947, 953, 967, 971, 977, 983, 991, 997}
}

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

func genNums(n int) (nums []int) {
	for i := 0; i < n; i++ {
		nums = append(nums, i)
	}
	return nums
}

func randomise(n int) (key []int) {
	nums := genNums(n)
	for i := 0; i < n; i++ {
		r := rand.Intn(len(nums))
		key = append(key, nums[r])
		nums = append(nums[:r], nums[r+1:]...) // remove from nums
	}
	return key
}

func removeSpaces(original string) (ciphertext string) {
	for i := 0; i < len(original); i++ {
		if original[i] != ' ' {
			ciphertext += string(original[i])
		}
	}
	return ciphertext
}

func factorise(n int) (factorisation []int) {
	for {
		if n == 1 {
			return factorisation
		}
		for i := 0; i < len(primes()); i++ {
			if n%primes()[i] == 0 {
				factorisation = append(factorisation, primes()[i])
				n /= primes()[i]
				break
			}
		}
	}
}

func possibleRowLengths(factorisation []int) (uniqueFactors []int) {
	for i := 0; i < len(factorisation); i++ {
		new := true
		if len(uniqueFactors) == 0 {
			uniqueFactors = append(uniqueFactors, factorisation[i])
		}
		for j := 0; j < len(uniqueFactors); j++ {
			if uniqueFactors[j] == factorisation[i] {
				new = false
			}
		}
		if new {
			uniqueFactors = append(uniqueFactors, factorisation[i])
		}
	}
	return uniqueFactors
}

func print(b []byte, columnSize int) {
	for i := 0; i < len(b); i++ {
		if i%columnSize == 0 {
			fmt.Println()
		}
		fmt.Print(string(b[i]))
	}
}

func rotate(text []byte, columnSize int) []byte { // rotates by 90 degrees clockwise, can then be called however many times to get the right angle
	var spaces string
	l := len(text)
	for l%columnSize != 0 {
		spaces += " "
		l++
	}

	b := []byte(string(text) + spaces)
	fmt.Println(len(b))
	columns := len(b) / columnSize

	if len(b)%columnSize != 0 {
		columns += 1
	}
	var res []byte
	for i := 0; i < columnSize; i++ {
		j := columnSize * (columns - 1)

		for k := j + i; k >= 0; k -= columnSize { // column currently on
			fmt.Println(k)
			res = append(res, b[k])
		}

	}
	fmt.Println(string(res))
	return res
}

func tryKeyLength(ciphertext string, keyLength int) (decoded string) {
	for i := 0; i < len(ciphertext)/keyLength; i++ {
		for j := i; j < len(ciphertext); j += len(ciphertext) / keyLength {
			decoded += string(ciphertext[j])
		}
	}
	return decoded
}

func decodePermutation(scrambled string, permutation []int) (result string) {
	keyLength := len(permutation)
	for i := 0; i < len(scrambled); i += keyLength {
		if i+keyLength >= len(scrambled) {
			break
		}
		slice := scrambled[i : i+keyLength]
		for j := 0; j < len(permutation); j++ {
			result += string(slice[permutation[j]])
		}
	}
	return result
}

func mutate(parent []int, n int) []int {
	newKey := slices.Clone(parent)
	i1 := rand.Intn(n)
	i2 := rand.Intn(n)
	for i2 == i1 {
		i2 = rand.Intn(n) // keep generating random numbers until not equal to i1
	}

	newKey[i1], newKey[i2] = newKey[i2], newKey[i1]
	return newKey
}

func hillClimb(maxConstant int, n int, ciphertext string) {
	total := 0
	alpha := key{randomise(n), 0}
	for i := 0; i < maxConstant; i++ {
		total++
		p := mutate(alpha.k, n)

		p1 := decodePermutation(ciphertext, alpha.k)
		p2 := decodePermutation(ciphertext, p)
		if score(p2) > score(p1) {
			alpha.k = p
			alpha.score = score(p2)
			i = 0
		}
	}
	fmt.Println(decodePermutation(ciphertext, alpha.k))
	fmt.Println(total)
}

func reverse(s string) (rev string) {
	for i := len(s) - 1; i >= 0; i-- {
		rev += string(s[i])
	}
	return rev
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := removeSpaces(string(original))
	fmt.Println(factorise(len(ciphertext)))
	//_ := tryKeyLength(ciphertext, 17) // change number here based on factorisation
	// fmt.Println(decoded)
	print(rotate([]byte("qwertyuio"), 5), 2)
}
