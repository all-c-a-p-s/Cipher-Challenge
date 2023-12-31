package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/wcharczuk/go-chart/v2"
)

type character struct {
	char byte
	freq float32
}

type nGram struct {
	nG   string
	freq float32
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func letters() []byte { // initialiser function for letters
	return []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
}

func referenceFrequencies() map[byte]float32 { // initialiser
	return map[byte]float32{ // from https://www3.nd.edu/~busiforc/handouts/cryptography/letterfrequencies.html
		'A': 0.084966,
		'B': 0.020720,
		'C': 0.045388,
		'D': 0.033844,
		'E': 0.111607,
		'F': 0.018121,
		'G': 0.024705,
		'H': 0.030034,
		'I': 0.075448,
		'J': 0.001965,
		'K': 0.011016,
		'L': 0.054893,
		'M': 0.030129,
		'N': 0.066544,
		'O': 0.071635,
		'P': 0.031671,
		'Q': 0.001962,
		'R': 0.075809,
		'S': 0.057351,
		'T': 0.069509,
		'U': 0.036308,
		'V': 0.010074,
		'W': 0.012899,
		'X': 0.002902,
		'Y': 0.017779,
		'Z': 0.002722,
	}
}

func format(original []byte) (formatted []byte) {
	// formats text to remove spaces and characters the are not uppercase letters
	for i := 0; i < len(original); i++ {
		if 65 <= original[i] && original[i] <= 90 {
			formatted = append(formatted, original[i])
		}
	}
	return formatted
}

func removeSpaces(original []byte) (noSpaces []byte) {
	// only removes spaces
	for i := 0; i < len(original); i++ {
		if original[i] != 32 {
			noSpaces = append(noSpaces, original[i])
		}
	}
	return noSpaces
}

func formatBinary(original []byte) (formatted string) {
	// formats binary-encoded ciphertext
	for i := 0; i < len(original); i++ {
		if 48 <= original[i] && original[i] <= 49 {
			formatted += string(original[i])
		}
	}
	return formatted
}

func read() ([]byte, []byte) {
	// reads file, returning both formatted and unformatted versions
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	formatted := format(original) // formatted to remove all chars that are not capital letters
	return original, formatted
}

func randomIOC(text []byte) float32 {
	// takes random indices in the text and compares them to estimate IOC
	var coincidences, total int
	for i := 0; i < 1_000_000; i++ {
		i1 := rand.Intn(len(text))
		i2 := rand.Intn(len(text))

		for i1 == i2 {
			i2 = rand.Intn(len(text))
		}

		if text[i1] == text[i2] {
			coincidences++
		}
		total++
	}
	return float32(coincidences) / float32(total) * float32(26)
	// normalised by * 26 to be more readable
}

func shiftIOC(max int, text []byte) []float32 {
	// calculates coincidences between original text and shifted texts
	var IOCs []float32
	for i := 1; i < max; i++ {
		var coincidences, total float32
		n := make([]byte, len(text))
		for j := 0; j < len(text); j++ {
			k := (j + i) % len(text)
			n[k] = text[j]
		}
		for l := 0; l < len(text); l++ {
			if text[l] == n[l] {
				coincidences++
			}
			total++
		}
		IOCs = append(IOCs, 26*coincidences/total)
	}
	return IOCs
}

func analyseFrequency(text []byte) map[byte]int {
	// scans text and returns hashmap of letter frequencies
	frequencyDistribution := map[byte]int{}
	for _, b := range text {
		if _, ok := frequencyDistribution[b]; !ok {
			frequencyDistribution[b] = 1
		} else {
			frequencyDistribution[b]++
		}
	}
	return frequencyDistribution
}

func analyseNGrams(s string, n int) {
	// function to analyse frequencies of n-grams in the text
	var nGrams []string
	if len(s)%n != 0 {
		panic("text is not divisible by split length")
	}
	for i := 0; i < len(s)-n+1; i += n {
		// loop through text appending all n-grams to slice
		slice := s[i : i+n]
		nGrams = append(nGrams, slice)
	}
	frequencyDistribution := map[string]int{}
	for _, b := range nGrams {
		// go through slice and read into hashmap
		if _, ok := frequencyDistribution[b]; !ok {
			frequencyDistribution[b] = 1
		} else {
			frequencyDistribution[b]++
		}
	}
	var total int
	for _, v := range frequencyDistribution {
		total += v
	}
	relativeFrequency := map[string]float32{}

	for k, v := range frequencyDistribution {
		// where relative frequency = appearances in text/ total number of n-grams
		relativeFrequency[k] = float32(v) / float32(total)
	}

	ngrams := []nGram{}

	for k, v := range relativeFrequency {
		l := nGram{k, v}
		ngrams = append(ngrams, l)
	}

	fmt.Printf("%d Unique n-grams \n", len(relativeFrequency))

	sort.Slice(ngrams, func(i, j int) bool {
		return ngrams[i].freq < ngrams[j].freq
	})

	for i := len(ngrams) - 1; i >= 0; i-- {
		fmt.Printf(`{Value: %f, Label: "%s" },`, ngrams[i].freq, ngrams[i].nG)
		fmt.Println()
	} // printed out in CSS format so it can be pasted into the graphing library
}

func graph() {
	// uses graphing library to create bar chart of n-gram frequencies
	graph := chart.BarChart{
		Title: "Frequency Distribution",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 30,
		// CSS below is output by the program the first time, then can be pasted into here
		Bars: []chart.Value{
			{Value: 0.078197, Label: "A"},
			{Value: 0.047378, Label: "B"},
			{Value: 0.016099, Label: "C"},
			{Value: 0.026679, Label: "D"},
			{Value: 0.028519, Label: "E"},
			{Value: 0.027599, Label: "F"},
			{Value: 0.018629, Label: "G"},
			{Value: 0.056808, Label: "H"},
			{Value: 0.054738, Label: "I"},
			{Value: 0.000000, Label: "J"},
			{Value: 0.048988, Label: "K"},
			{Value: 0.043928, Label: "L"},
			{Value: 0.037489, Label: "M"},
			{Value: 0.055658, Label: "N"},
			{Value: 0.041168, Label: "O"},
			{Value: 0.050598, Label: "P"},
			{Value: 0.068077, Label: "Q"},
			{Value: 0.026449, Label: "R"},
			{Value: 0.064397, Label: "S"},
			{Value: 0.044618, Label: "T"},
			{Value: 0.018399, Label: "U"},
			{Value: 0.034039, Label: "V"},
			{Value: 0.048068, Label: "W"},
			{Value: 0.024839, Label: "X"},
			{Value: 0.019089, Label: "Y"},
			{Value: 0.019549, Label: "Z"},
			{Value: 0.000000, Label: "XAXIS"}, // in here so that the x-axis always starts at 0
		},
	}

	output, err := os.Create("graph.png")
	check(err)
	defer output.Close()
	graph.Render(chart.PNG, output)
}

func referenceScore(relativeFrequency map[byte]float32) float32 { // reference score < 10^5 is usually ony transposition
	// function I made to suggest whether a text has substitution
	// loops throught the map, multiplying the score by the scale factor between
	// frequency in the ciphertext and frequency in English
	var refScore float32 = 1.0 // initialise at 1
	for i := 0; i < len(letters()); i++ {
		if _, ok := relativeFrequency[letters()[i]]; ok {
			r := referenceFrequencies()[letters()[i]]
			n := relativeFrequency[letters()[i]]
			if r > n {
				refScore *= r / n
			} else {
				refScore *= n / r
			} // multiply by the scale factor (always greater than one)
			// between ciphertext frequency: english frequency
		}
	}
	return refScore
}

func dotProduct(vec1 []float32, vec2 []float32) float32 {
	if len(vec1) != len(vec2) {
		panic("Attempt to calculate dot product of vectors of different length")
	}
	var sum float32
	for i := 0; i < len(vec1); i++ {
		sum += vec1[i] * vec2[i]
	}
	return sum
}

func vecLength(v []float32) float32 {
	return float32(math.Sqrt(float64(dotProduct(v, v))))
}

func cosineBetweenVectors(vec1 []float32, vec2 []float32) float32 {
	numerator := dotProduct(vec1, vec2)
	denominator := vecLength(vec1) * vecLength(vec2)
	return numerator / denominator
}

func angleScore(relativeFrequency map[byte]float32) float32 {
	// values > 0.97 suggest no substitution used
	var vec1, vec2 []float32 // frequencies from text, from english
	for _, l := range letters() {
		vec1 = append(vec1, relativeFrequency[l])
		vec2 = append(vec2, referenceFrequencies()[l])
	}
	return cosineBetweenVectors(vec1, vec2)
}

func print(m map[byte]int) {
	var total int
	for _, v := range m {
		total += v
	}
	relativeFrequency := map[byte]float32{}

	for k, v := range m {
		relativeFrequency[k] = float32(v) / float32(total)
	}

	chars := []character{}

	for k, v := range relativeFrequency {
		l := character{k, v}
		chars = append(chars, l)
	}

	fmt.Printf("%d Unique Characters \n", len(relativeFrequency))

	sort.Slice(chars, func(i, j int) bool {
		return chars[i].freq < chars[j].freq
	})

	fmt.Printf("Reference Score: %f \n", referenceScore(relativeFrequency))
	fmt.Printf("Cosine between vectors: %f \n", angleScore(relativeFrequency))

	/** uncomment to print in order of frequncy
	for i := len(chars) - 1; i >= 0; i-- {
		fmt.Printf(`{Value: %f, Label: "%s" },`, chars[i].freq, string(chars[i].char))
		fmt.Println()
	}

	fmt.Println()
	*/

	for i := 0; i < len(letters()); i++ {
		fmt.Printf(`{Value: %f, Label: "%s" },`, relativeFrequency[letters()[i]], string(letters()[i]))
		fmt.Println()
	}
}

func main() { // uncomment whatever you want here
	_, formatted := read()
	// bin := formatBinary(unformatted)
	fmt.Printf("Length: %d \n", len(formatted))
	fmt.Printf("IOC: %f \n", randomIOC(formatted))
	fmt.Println("Shift IOCs: ", shiftIOC(8, formatted))
	// fmt.Println("UNFORMATTED: ")
	// print(analyseFrequency([]byte(strings.TrimSpace(string(unformatted)))))
	// fmt.Print("\n\n\n\n")
	// fmt.Println("FORMATTED: ")
	print(analyseFrequency(formatted))
	// analyseNGrams(bin, 5)
	graph()
}
