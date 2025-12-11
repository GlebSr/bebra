package utils

import (
	"math/rand"
	"time"
)

func GenerateDiscriminator() int {
	// Implement a simple random discriminator generator (e.g., between 0001 and 9999)
	// In a real-world scenario, you might want to ensure uniqueness within a certain scope
	rand.Seed(time.Now().UnixNano())
	return 1000 + rand.Intn(9000)
}
