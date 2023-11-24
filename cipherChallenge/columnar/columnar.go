package main

import (
	"fmt"
	"log"
	"os"
)

var permutations [][]int

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func primes() []int { //initialiser function returning primes up to 1000
	return []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541, 547, 557, 563, 569, 571, 577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643, 647, 653, 659, 661, 673, 677, 683, 691, 701, 709, 719, 727, 733, 739, 743, 751, 757, 761, 769, 773, 787, 797, 809, 811, 821, 823, 827, 829, 839, 853, 857, 859, 863, 877, 881, 883, 887, 907, 911, 919, 929, 937, 941, 947, 953, 967, 971, 977, 983, 991, 997}

}

func dictionary() []string {
	return []string{"THE", "OF", "TO", "AND", "A", "IN", "IS", "IT", "YOU", "THAT", "HE", "WAS", "FOR", "ON", "ARE", "WITH", "AS", "I", "HIS", "THEY", "BE", "AT", "ONE", "HAVE", "THIS", "FROM", "OR", "HAD", "BY", "NOT", "WORD", "BUT", "WHAT", "SOME", "WE", "CAN", "OUT", "OTHER", "WERE", "ALL", "THERE", "WHEN", "UP", "USE", "YOUR", "HOW", "SAID", "AN", "EACH", "SHE", "WHICH", "DO", "THEIR", "TIME", "IF", "WILL", "WAY", "ABOUT", "MANY", "THEN", "THEM", "WRITE", "WOULD", "LIKE", "SO", "THESE", "HER", "LONG", "MAKE", "THING", "SEE", "HIM", "TWO", "HAS", "LOOK", "MORE", "DAY", "COULD", "GO", "COME", "DID", "NUMBER", "SOUND", "NO", "MOST", "PEOPLE", "MY", "OVER", "KNOW", "WATER", "THAN", "CALL", "FIRST", "WHO", "MAY", "DOWN", "SIDE", "BEEN", "NOW", "FIND", "ANY", "NEW", "WORK", "PART", "TAKE", "GET", "PLACE", "MADE", "LIVE", "WHERE", "AFTER", "BACK", "LITTLE", "ONLY", "ROUND", "MAN", "YEAR", "CAME", "SHOW", "EVERY", "GOOD", "ME", "GIVE", "OUR", "UNDER", "NAME", "VERY", "THROUGH", "JUST", "FORM", "SENTENCE", "GREAT", "THINK", "SAY", "HELP", "LOW", "LINE", "DIFFER", "TURN", "CAUSE", "MUCH", "MEAN", "BEFORE", "MOVE", "RIGHT", "BOY", "OLD", "TOO", "SAME", "TELL", "DOES", "SET", "THREE", "WANT", "AIR", "WELL", "ALSO", "PLAY", "SMALL", "END", "PUT", "HOME", "READ", "HAND", "PORT", "LARGE", "SPELL", "ADD", "EVEN", "LAND", "HERE", "MUST", "BIG", "HIGH", "SUCH", "FOLLOW", "ACT", "WHY", "ASK", "MEN", "CHANGE", "WENT", "LIGHT", "KIND", "OFF", "NEED", "HOUSE", "PICTURE", "TRY", "US", "AGAIN", "ANIMAL", "POINT", "MOTHER", "WORLD", "NEAR", "BUILD", "SELF", "EARTH", "FATHER", "HEAD", "STAND", "OWN", "PAGE", "SHOULD", "COUNTRY", "FOUND", "ANSWER", "SCHOOL", "GROW", "STUDY", "STILL", "LEARN", "PLANT", "COVER", "FOOD", "SUN", "FOUR", "BETWEEN", "STATE", "KEEP", "EYE", "NEVER", "LAST", "LET", "THOUGHT", "CITY", "TREE", "CROSS", "FARM", "HARD", "START", "MIGHT", "STORY", "SAW", "FAR", "SEA", "DRAW", "LEFT", "LATE", "RUN", "DON'T", "WHILE", "PRESS", "CLOSE", "NIGHT", "REAL", "LIFE", "FEW", "NORTH", "OPEN", "SEEM", "TOGETHER", "NEXT", "WHITE", "CHILDREN", "BEGIN", "GOT", "WALK", "EXAMPLE", "EASE", "PAPER", "GROUP", "ALWAYS", "MUSIC", "THOSE", "BOTH", "MARK", "OFTEN", "LETTER", "UNTIL", "MILE", "RIVER", "CAR", "FEET", "CARE", "SECOND", "BOOK", "CARRY", "TOOK", "SCIENCE", "EAT", "ROOM", "FRIEND", "BEGAN", "IDEA", "FISH", "MOUNTAIN", "STOP", "ONCE", "BASE", "HEAR", "HORSE", "CUT", "SURE", "WATCH", "COLOR", "FACE", "WOOD", "MAIN", "ENOUGH", "PLAIN", "GIRL", "USUAL", "YOUNG", "READY", "ABOVE", "EVER", "RED", "LIST", "THOUGH", "FEEL", "TALK", "BIRD", "SOON", "BODY", "DOG", "FAMILY", "DIRECT", "POSE", "LEAVE", "SONG", "MEASURE", "DOOR", "PRODUCT", "BLACK", "SHORT", "NUMERAL", "CLASS", "WIND", "QUESTION", "HAPPEN", "COMPLETE", "SHIP", "AREA", "HALF", "ROCK", "ORDER", "FIRE", "SOUTH", "PROBLEM", "PIECE", "TOLD", "KNEW", "PASS", "SINCE", "TOP", "WHOLE", "KING", "SPACE", "HEARD", "BEST", "HOUR", "BETTER", "TRUE", "DURING", "HUNDRED", "FIVE", "REMEMBER", "STEP", "EARLY", "HOLD", "WEST", "GROUND", "INTEREST", "REACH", "FAST", "VERB", "SING", "LISTEN", "SIX", "TABLE", "TRAVEL", "LESS", "MORNING", "TEN", "SIMPLE", "SEVERAL", "VOWEL", "TOWARD", "WAR", "LAY", "AGAINST", "PATTERN", "SLOW", "CENTER", "LOVE", "PERSON", "MONEY", "SERVE", "APPEAR", "ROAD", "MAP", "RAIN", "RULE", "GOVERN", "PULL", "COLD", "NOTICE", "VOICE", "UNIT", "POWER", "TOWN", "FINE", "CERTAIN", "FLY", "FALL", "LEAD", "CRY", "DARK", "MACHINE", "NOTE", "WAIT", "PLAN", "FIGURE", "STAR", "BOX", "NOUN", "FIELD", "REST", "CORRECT", "ABLE", "POUND", "DONE", "BEAUTY", "DRIVE", "STOOD", "CONTAIN", "FRONT", "TEACH", "WEEK", "FINAL", "GAVE", "GREEN", "OH", "QUICK", "DEVELOP", "OCEAN", "WARM", "FREE", "MINUTE", "STRONG", "SPECIAL", "MIND", "BEHIND", "CLEAR", "TAIL", "PRODUCE", "FACT", "STREET", "INCH", "MULTIPLY", "NOTHING", "COURSE", "STAY", "WHEEL", "FULL", "FORCE", "BLUE", "OBJECT", "DECIDE", "SURFACE", "DEEP", "MOON", "ISLAND", "FOOT", "SYSTEM", "BUSY", "TEST", "RECORD", "BOAT", "COMMON", "GOLD", "POSSIBLE", "PLANE", "STEAD", "DRY", "WONDER", "LAUGH", "THOUSAND", "AGO", "RAN", "CHECK", "GAME", "SHAPE", "EQUATE", "HOT", "MISS", "BROUGHT", "HEAT", "SNOW", "TIRE", "BRING", "YES", "DISTANT", "FILL", "EAST", "PAINT", "LANGUAGE", "AMONG", "GRAND", "BALL", "YET", "WAVE", "DROP", "HEART", "AM", "PRESENT", "HEAVY", "DANCE", "ENGINE", "POSITION", "ARM", "WIDE", "SAIL", "MATERIAL", "SIZE", "VARY", "SETTLE", "SPEAK", "WEIGHT", "GENERAL", "ICE", "MATTER", "CIRCLE", "PAIR", "INCLUDE", "DIVIDE", "SYLLABLE", "FELT", "PERHAPS", "PICK", "SUDDEN", "COUNT", "SQUARE", "REASON", "LENGTH", "REPRESENT", "ART", "SUBJECT", "REGION", "ENERGY", "HUNT", "PROBABLE", "BED", "BROTHER", "EGG", "RIDE", "CELL", "BELIEVE", "FRACTION", "FOREST", "SIT", "RACE", "WINDOW", "STORE", "SUMMER", "TRAIN", "SLEEP", "PROVE", "LONE", "LEG", "EXERCISE", "WALL", "CATCH", "MOUNT", "WISH", "SKY", "BOARD", "JOY", "WINTER", "SAT", "WRITTEN", "WILD", "INSTRUMENT", "KEPT", "GLASS", "GRASS", "COW", "JOB", "EDGE", "SIGN", "VISIT", "PAST", "SOFT", "FUN", "BRIGHT", "GAS", "WEATHER", "MONTH", "MILLION", "BEAR", "FINISH", "HAPPY", "HOPE", "FLOWER", "CLOTHE", "STRANGE", "GONE", "JUMP", "BABY", "EIGHT", "VILLAGE", "MEET", "ROOT", "BUY", "RAISE", "SOLVE", "METAL", "WHETHER", "PUSH", "SEVEN", "PARAGRAPH", "THIRD", "SHALL", "HELD", "HAIR", "DESCRIBE", "COOK", "FLOOR", "EITHER", "RESULT", "BURN", "HILL", "SAFE", "CAT", "CENTURY", "CONSIDER", "TYPE", "LAW", "BIT", "COAST", "COPY", "PHRASE", "SILENT", "TALL", "SAND", "SOIL", "ROLL", "TEMPERATURE", "FINGER", "INDUSTRY", "VALUE", "FIGHT", "LIE", "BEAT", "EXCITE", "NATURAL", "VIEW", "SENSE", "EAR", "ELSE", "QUITE", "BROKE", "CASE", "MIDDLE", "KILL", "SON", "LAKE", "MOMENT", "SCALE", "LOUD", "SPRING", "OBSERVE", "CHILD", "STRAIGHT", "CONSONANT", "NATION", "DICTIONARY", "MILK", "SPEED", "METHOD", "ORGAN", "PAY", "AGE", "SECTION", "DRESS", "CLOUD", "SURPRISE", "QUIET", "STONE", "TINY", "CLIMB", "COOL", "DESIGN", "POOR", "LOT", "EXPERIMENT", "BOTTOM", "KEY", "IRON", "SINGLE", "STICK", "FLAT", "TWENTY", "SKIN", "SMILE", "CREASE", "HOLE", "TRADE", "MELODY", "TRIP", "OFFICE", "RECEIVE", "ROW", "MOUTH", "EXACT", "SYMBOL", "DIE", "LEAST", "TROUBLE", "SHOUT", "EXCEPT", "WROTE", "SEED", "TONE", "JOIN", "SUGGEST", "CLEAN", "BREAK", "LADY", "YARD", "RISE", "BAD", "BLOW", "OIL", "BLOOD", "TOUCH", "GREW", "CENT", "MIX", "TEAM", "WIRE", "COST", "LOST", "BROWN", "WEAR", "GARDEN", "EQUAL", "SENT", "CHOOSE", "FELL", "FIT", "FLOW", "FAIR", "BANK", "COLLECT", "SAVE", "CONTROL", "DECIMAL", "GENTLE", "WOMAN", "CAPTAIN", "PRACTICE", "SEPARATE", "DIFFICULT", "DOCTOR", "PLEASE", "PROTECT", "NOON", "WHOSE", "LOCATE", "RING", "CHARACTER", "INSECT", "CAUGHT", "PERIOD", "INDICATE", "RADIO", "SPOKE", "ATOM", "HUMAN", "HISTORY", "EFFECT", "ELECTRIC", "EXPECT", "CROP", "MODERN", "ELEMENT", "HIT", "STUDENT", "CORNER", "PARTY", "SUPPLY", "BONE", "RAIL", "IMAGINE", "PROVIDE", "AGREE", "THUS", "CAPITAL", "WON'T", "CHAIR", "DANGER", "FRUIT", "RICH", "THICK", "SOLDIER", "PROCESS", "OPERATE", "GUESS", "NECESSARY", "SHARP", "WING", "CREATE", "NEIGHBOR", "WASH", "BAT", "RATHER", "CROWD", "CORN", "COMPARE", "POEM", "STRING", "BELL", "DEPEND", "MEAT", "RUB", "TUBE", "FAMOUS", "DOLLAR", "STREAM", "FEAR", "SIGHT", "THIN", "TRIANGLE", "PLANET", "HURRY", "CHIEF", "COLONY", "CLOCK", "MINE", "TIE", "ENTER", "MAJOR", "FRESH", "SEARCH", "SEND", "YELLOW", "GUN", "ALLOW", "PRINT", "DEAD", "SPOT", "DESERT", "SUIT", "CURRENT", "LIFT", "ROSE", "CONTINUE", "BLOCK", "CHART", "HAT", "SELL", "SUCCESS", "COMPANY", "SUBTRACT", "EVENT", "PARTICULAR", "DEAL", "SWIM", "TERM", "OPPOSITE", "WIFE", "SHOE", "SHOULDER", "SPREAD", "ARRANGE", "CAMP", "INVENT", "COTTON", "BORN", "DETERMINE", "QUART", "NINE", "TRUCK", "NOISE", "LEVEL", "CHANCE", "GATHER", "SHOP", "STRETCH", "THROW", "SHINE", "PROPERTY", "COLUMN", "MOLECULE", "SELECT", "WRONG", "GRAY", "REPEAT", "REQUIRE", "BROAD", "PREPARE", "SALT", "NOSE", "PLURAL", "ANGER", "CLAIM", "CONTINENT", "OXYGEN", "SUGAR", "DEATH", "PRETTY", "SKILL", "WOMEN", "SEASON", "SOLUTION", "MAGNET", "SILVER", "THANK", "BRANCH", "MATCH", "SUFFIX", "ESPECIALLY", "FIG", "AFRAID", "HUGE", "SISTER", "STEEL", "DISCUSS", "FORWARD", "SIMILAR", "GUIDE", "EXPERIENCE", "SCORE", "APPLE", "BOUGHT", "LED", "PITCH", "COAT", "MASS", "CARD", "BAND", "ROPE", "SLIP", "WIN", "DREAM", "EVENING", "CONDITION", "FEED", "TOOL", "TOTAL", "BASIC", "SMELL", "VALLEY", "NOR", "DOUBLE", "SEAT", "ARRIVE", "MASTER", "TRACK", "PARENT", "SHORE", "DIVISION", "SHEET", "SUBSTANCE", "FAVOR", "CONNECT", "POST", "SPEND", "CHORD", "FAT", "GLAD", "ORIGINAL", "SHARE", "STATION", "DAD", "BREAD", "CHARGE", "PROPER", "BAR", "OFFER", "SEGMENT", "SLAVE", "DUCK", "INSTANT", "MARKET", "DEGREE", "POPULATE", "CHICK", "DEAR", "ENEMY", "REPLY", "DRINK", "OCCUR", "SUPPORT", "SPEECH", "NATURE", "RANGE", "STEAM", "MOTION", "PATH", "LIQUID", "LOG", "MEANT", "QUOTIENT", "TEETH", "SHELL", "NECK"}
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
		slice := scrambled[i : i+keyLength]
		for j := 0; j < len(permutation); j++ {
			result += string(slice[permutation[j]])
		}
	}
	return result
}

func main() {
	original, err := os.ReadFile("ciphertext.txt")
	check(err)
	ciphertext := removeSpaces(string(original))
	decoded := tryKeyLength(ciphertext, 7)
	fmt.Println(decodePermutation(decoded, []int{5, 3, 6, 0, 2, 1, 4}))
}
