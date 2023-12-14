package main

import "os"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func letters() []byte {
	return []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y'} // remove J or Z as necessary
}

func format(original []byte) (result string) {
	for i := 0; i < len(original); i++ {
		if 65 <= original[i] && original[i] <= 90 {
			result += string(original[i])
		}
	}
	return result
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

func encipher(text, keyword string) string {
	grid := fillGrid(keyword)
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)

	ciphertext := format(original)
}
