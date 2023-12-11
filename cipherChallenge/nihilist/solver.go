package main

import (
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

func format(b []byte) (s string) {
	for i := 0; i < len(b); i++ {
		if (b[i] >= 48 && b[i] <= 57) || b[i] == 32 { // numbers or space
			s += string(b[i])
		}
	}
	return s
}

func convertToPolybius(s string) (p []int) {
	polybius := map[byte]int{
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
	for i := 0; i < len(s); i++ {
		p = append(p, polybius[s[i]])
	}
	return p
}

func decipher(key1 []int) {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := format(original)

	numsStr := strings.Fields(ciphertext)
	var nums []int
	for _, n := range numsStr {
		num, err := strconv.Atoi(n)
		check(err)
		nums = append(nums, num)
	}

	var d1 []int // partially deciphered

	for i := 0; i < len(nums); i++ {
		shift := key1[i%len(key1)]
		d1 = append(d1, nums[i]-shift)
	}

	for _, v := range d1 {
		fmt.Printf("%d ", v)
	}
}

func main() {
	decipher(convertToPolybius("NOTHING"))
}
