package utils

func CheckLuhn(number string) bool {
	digits := make([]int, len(number))
	for i, char := range number {
		if char < '0' || char > '9' {
			return false
		}
		digits[i] = int(char - '0')
	}

	length := len(digits)
	if length == 0 {
		return false
	}

	sum := 0
	parity := length % 2

	for i := 0; i < length-1; i++ {
		if i%2 != parity {
			sum += digits[i]
		} else if digits[i] > 4 {
			sum += 2*digits[i] - 9
		} else {
			sum += 2 * digits[i]
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	return digits[length-1] == checkDigit
}
