# Testing Simulated Annealing Versus the Genetic Algorithm

I used three plaintexts to compare the monoalphabetic solvers. Each of these were generated using the website (blindtextgenerator.com)[blindtextgenerator.com]:
- Cicero(en) - 100 words
- Li Europan Lingues (en) - 50 words
- Kafka - 25 words


I used the file randomGenerator/randomCiphertext.go to randomly generate keys to encipher these ciphertexts with random keys
I also added an iteration limit of 10000 iterations to the simulated annealing script because otherwise it would run indefinitely for the 25 word text.

# Results:
| Test                | Simulated Annealing | Genetic Algorithm |
| ----------------- | ------------------- | ----------------- |
| Cicero 100 (1)                     |  ✅ 5.76                                |  ✅ 5.76                            |
| Cicero 100 (2)                     | ❌ 3.65                                 | ✅ 5.76                             |
| Cicero 100 (3)                     | ✅ 5.76                                 | ✅ 5.76                             |
| Cicero 100 (4)                     | ❌ 5.16                                 | ✅ 5.76                             |
| Cicero 100 (5)                     | ✅ 5.76                                 | ✅ 5.76                             |
| Li Europan 50 (1)                  | ✅ 5.67                                 | ✅ 5.67                             |
| Li Europan 50 (2)                  | ❌ 2.89                                 | ❌ 3.39                             |
| Li Europan 50 (3)                  | ❌ 3.27                                 | ✅ 5.67                             |
| Li Europan 50 (4)                  | ❌ 5.33                                 | ✅ 5.67                             |
| Li Europan 50 (5)                  | ✅ 5.67                                 | ✅ 5.67                             |
| Kafka 25 (1)                       | ❌ 3.48                                 | ❌ 4.24                             |
| Kafka 25 (2)                       | ❌ 4.17                                 | ✅ 5.34                             |
| Kafka 25 (3)                       | ❌ 3.39                                 | ❌ 3.75                             |
| Kafka 25 (4)                       | ❌ 3.59                                 | ❌ 4.08                             |
| Kafka 25 (5)                       | ❌ 4.00                                 | ❌ 4.01                             |
| Mean Score:                        | 4.50                                   | 5.09                               |

# Conclusion

Overall, the genetic algorithm was much more consistent:
- It found a fully correct solution 10 times compared to 5
- It had a higher mean fitness score for its solution
- It was the only algorithm to solve the 25 word ciphertext 
