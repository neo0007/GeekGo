package main

import (
	"errors"
	"fmt"
)

func deleteElement[T any](slice []T, index int) ([]T, error) {
	if index < 0 || index >= len(slice) {
		return nil, errors.New("invalid index")
	}

	if len(slice) == 1 {
		return nil, errors.New("cannot delete the only element in the slice")
	}

	copy(slice[index:], slice[index+1:])
	return slice[:len(slice)-1], nil
}

func main() {
	var slice = []int{1, 2, 6, 8, 5, 6}
	indexToRemove := 2

	fmt.Println("Original slice:", slice)
	updatedSlice, err := deleteElement(slice, indexToRemove)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if updatedSlice == nil {
		fmt.Println("The slice became empty.")
	} else {
		fmt.Println("Modified slice:", updatedSlice)

		// Perform manual resize if needed
		if len(updatedSlice) < cap(updatedSlice)/2 {
			newSlice := make([]int, len(updatedSlice), len(updatedSlice)*2)
			copy(newSlice, updatedSlice)
			updatedSlice = newSlice
		}

		fmt.Println("Resized slice:", updatedSlice)
	}
}
