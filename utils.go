package main

import (
	color "github.com/lucasb-eyer/go-colorful"
	"math/rand"
	"time"
)

// TODO rand.Seed is deprecated, switch to rand.Seed(seed int64)
func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateColor() string {
	c := color.Hsv(rand.Float64()*360.0, 0.75, 0.5)
	return c.Hex()
}
