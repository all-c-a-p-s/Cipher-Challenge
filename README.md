# Cipher Challenge
Code I've used to solve problems from the British National Cipher Challenge. Feel free to use to solve the various ciphers included:

| Solver                              | Description                                                                                                                                                                                                                                    |
|-------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Analysis tool (Go)                  | preforms analysis of ciphertext, giving number of unique letters, letter frequency distribution, and the likelihood that some form of substitution was used. It also has the facility to analyse the frequencies of n-grams in the ciphertext. |
| Baudot decoder (Go)                 | decodes ciphertext with baudot code, with the facility of trying transposition permutations until a valid string of baudot codes is found.                                                                                                     |
| Bifid decoder (Go)                  | implements dictionary attack                                                                                                                                                                                                                   |
| Columnar Transposition decoder (Go) | Decodes text based on trying all permutations up to a certain key size. Also contains facility for rotations of the ciphertext.                                                                                                                |
| Hill Cipher (Go)                    | Bruteforces all 2x2 matrices with numbers 0-10                                                                                                                                                                                                 |
| Monoalphabetic Solvers (Go, Rust)   | implementations of hill climbing, simulated annealing, and the genetic algorithm in Go. Simulated annealing also implemented in rust                                                                 |
| Nihilist decoder (Go)               | Not really much of a decoder, but useful if the user can guess the first keyword                                                                                                                                                               |
| Playfair decoder(Go)                | Implements dictionary attack                                                                                                                                                                                                                   |
| Polybius decoder (Go)               | dictionary attack                                                                                                                                                                                                                              |
| Railfence decoder (Go)              | bruteforce attack                                                                                                                                                                                                                              |
| Trifid decoder (Go)                 | dictionary attack                                                                                                                                                                                                                              |
| Two-Square decoder (Go)             | dictionary attack. only applies to two-square cipher where the squares are placed vertically on top of each other and filled horizontally.                                                                                                     |
| Vigenere Decoder (Go)               | uses statistical attack. It will crash if there are no spikes in IOC, indicating that Vigenere was almost certainly not used.                                                                                                                  |
| Caesar Cipher decoder (Rust)        | bruteforces all shifts from 0-26                                 |
| Four Square decoder (Rust)          | implements both a dictionary attack and simulated annealing with rennealing/reference score based on English sample text generated of same length. The simulated annealing with rennealing is very consistent but its runtime is between 30-60 secs.     |




# Usage: 
To use any of the tools, paste the ciphertext into the ```ciphertext.txt``` file, and run the corresponding source file.

N.B. make sure to run the rust scripts with ```cargo run --release``` for much better performance.
