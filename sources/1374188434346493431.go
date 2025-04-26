package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func countPalindromes(input string) int64 {
	modifiedString := "$#"
	for _, ch := range input {
		modifiedString += string(ch) + "#"
	}
	modifiedString += "@"

	n := len(modifiedString)
	P := make([]int, n)
	center, right := 0, 0

	for i := 1; i < n-1; i++ {
		mirror := 2*center - i
		if i < right {
			P[i] = min(right-i, P[mirror])
		}
		for i+1+P[i] < n && i-1-P[i] >= 0 && modifiedString[i+1+P[i]] == modifiedString[i-1-P[i]] {
			P[i]++
		}
		if i+P[i] > right {
			center = i
			right = i + P[i]
		}
	}

	var count int64 = 0
	for i := 1; i < n-1; i++ {
		count += int64((P[i] + 1) / 2)
	}
	return count
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	fmt.Println(countPalindromes(strings.ToLower(input)))
}
