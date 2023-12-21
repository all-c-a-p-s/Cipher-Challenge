package main

import (
	"fmt"
	"log"
	"os"
)

//! runtime ~30 mins, almost all of which if format() function

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func uppercase(x byte) byte {
	if x >= 97 && x <= 122 {
		return x - 32
	}
	panic("Attempt to convert non uppercase letter to lowercase")

}

func format(text string) (letters string) {

	for i := 0; i < len(text); i++ {
		if 65 <= text[i] && text[i] <= 90 {
			letters += string(text[i])
		} else if 97 <= text[i] && text[i] <= 122 {
			letters += string(uppercase(text[i]))
		}
	}
	fmt.Println("Finished formatting")
	return letters
}

func scanTetragrams() map[string]float64 {
	tetragrams := map[string]float64{}
	txt, err := os.ReadFile("brown_corpus.txt")
	check(err)
	text := format(string(txt))

  var l float64 = float64(len(text)/4) //total number of tetragrams

	for i := 0; i < len(text)-4; i++ {
		slice := text[i : i+4]
		if _, ok := tetragrams[slice]; ok {
			tetragrams[slice]++
		} else {
			tetragrams[slice] = 1
		}
	}

  for key, val := range tetragrams {
    tetragrams[key] = val / l
  }

	fmt.Println("Finished scanning")
	return tetragrams
}

func main() {
	tetragrams := scanTetragrams()
	file, err := os.Create("tetragrams.txt")
	defer file.Close()
	check(err)
	for key, val := range tetragrams {
		t := key
		v := fmt.Sprintf("%f", val)

		s := fmt.Sprintf(t + ", " + v + "\n")
		file.WriteString(s)
	}

	fmt.Println("Finished writing")

}
