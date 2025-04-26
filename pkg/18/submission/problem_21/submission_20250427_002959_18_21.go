package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func countPalindromes(input string) int64 {
	// Create modified string: $#a#b#c#@ for input "abc"
	modifiedString := "$#"
	for _, char := range input {
		modifiedString += string(char) + "#"
	}
	modifiedString += "@"

	n := len(modifiedString)
	P := make([]int, n)
	center, right := 0, 0

	// Manacher's algorithm
	for i := 1; i < n-1; i++ {
		mirror := 2*center - i
		if i < right {
			if P[mirror] < right-i {
				P[i] = P[mirror]
			} else {
				P[i] = right - i
			}
		}
		for i+(1+P[i]) < n && i-(1+P[i]) >= 0 && modifiedString[i+(1+P[i])] == modifiedString[i-(1+P[i])] {
			P[i]++
		}
		if i+P[i] > right {
			center = i
			right = i + P[i]
		}
	}

	// Count palindromes
	var count int64
	for i := 1; i < n-1; i++ {
		count += int64((P[i] + 1) / 2)
	}
	return count
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	fmt.Println(countPalindromes(strings.ToLower(input)))
}