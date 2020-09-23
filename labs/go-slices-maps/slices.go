package main

import "golang.org/x/tour/pic"

// Pic function from Golang tour.
func Pic(dx, dy int) [][]uint8 {
	img := make([][]uint8, dy)
	for x := range img {
		img[x] = make([]uint8, dx)
		for y := 0; y < dx; y++ {
			img[x][y] = uint8(y * x + x * y * 98 + 85)
		}
	}
	return img
}

func main() {
	pic.Show(Pic)
}