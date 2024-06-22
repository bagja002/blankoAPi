package tools

import (
	"math/rand"
	"time"
)

// Function to generate a random string of specified length
func RandomString(length int) string {
	// Define the characters to be used in the random string
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a slice to hold the random characters
	result := make([]byte, length)

	// Generate the random characters
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
