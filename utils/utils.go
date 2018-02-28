package utils

// Returns N number of chars starting form a given position
// (recall that in go the char at position i might not be the same as what we find in the same position within a string)
func GetChars(line string, startPos int, endPos int) string {
	var result string
	// start index
	i := 0
	// number of chars concatenated
	total := 0
	for _, char := range line {
		if (i >= startPos) && (total <= (endPos - startPos)) {
			result += string(char)
			total++
		}
		i++
	}
	return result
}
