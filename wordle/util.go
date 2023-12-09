package wordle

func isCapitalLetter(char rune) bool {
	return char >= 'A' && char <= 'Z'
}

func isNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

func isSpecialChar(char rune) bool {
	charSet := "!@#$%^&*()-_=+[]{}|;:'\",.<>/?"
	for _, c := range charSet {
		if char == c {
			return true
		}
	}
	return false
}
