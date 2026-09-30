package main

import (
	"fmt"
	"photo-export-system/photo"
)

func main() {
	insta := &photo.InstagramExport{}
	pr := &photo.PrintExport{}

	fmt.Println("--- Test 1: plain JPG (standard RGB engine) ---")
	r, err := insta.Export("vacation.jpg", []byte("image_bytes"))
	fmt.Println(r, err)

	fmt.Println("\n--- Test 2: RAW file (goes through the adapter) ---")
	r, err = pr.Export("wedding.cr2", []byte("raw_bytes"))
	fmt.Println(r, err)

	fmt.Println("\n--- Test 3: adapter error handling ---")
	// empty bytes should make the legacy code return -1
	_, err = insta.Export("error.cr2", []byte(""))
	if err != nil {
		fmt.Println("handled:", err)
	}
}
