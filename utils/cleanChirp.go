package utils

import (
	"slices"
	"strings"
)

func CleanChirp(line string) string {
	words := strings.Split(line, " ")
	toDelete := []string{"fornax", "sharbert", "kerfuffle"}
	for i, word := range words {
		if slices.Contains(toDelete, strings.ToLower(word)) {
			words = slices.Replace(words, i, i+1, "****")
		}
	}
	return strings.Join(words, " ")
}
