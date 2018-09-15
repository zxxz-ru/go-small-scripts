package main

import (
	"fmt"
	"os"
	"time"
)

func createArr() [1000]int {
	var arr [1000]int
	i := 0
	for i < 1000 {
		arr[i] = i
		i += 1
	}
	return arr
}

func searchIndexWithLoop(arr [1000]int, q int) (int, string, error) {
	t := time.Now()
	for i, v := range arr {
		if v == q {
			duration := time.Now().Sub(t)
			return i, duration.String(), nil
		}
	}
	return 0, "", fmt.Errorf("Search with loop never finds %d in array", q)
}

func searchIndexWithTree(arr [1000]int, q int) (int, string, error) {
	var (
		a      int = 0
		b      int = 999
		middle int
	)
	done := false
	t := time.Now()
	for !done {
		// check if can find true middle, if not check last element, move pointer
		if b%2 != 0 {
			if q == arr[b] {
				duration := time.Now().Sub(t)
				return b, duration.String(), nil
			}
			b = b - 1
		} else {
			if (b-a)%2 != 0 {
				if q == arr[b] {
					duration := time.Now().Sub(t)
					return b, duration.String(), nil
				}
				b = b - 1
			}
			middle = a + ((b - a) / 2)
			if arr[middle] < q {
				a = middle
			} else if arr[middle] > q {
				b = middle
			} else {
				duration := time.Now().Sub(t)
				return middle, duration.String(), nil
			}
		}
	}
	return 0, "", fmt.Errorf("Search with loop never finds %d in array", q)
}

func main() {
	s_arr := []int{33, 499, 897, 750, 749, 501, 2, 999}
	arr := createArr()
	for _, v := range s_arr {
		if index, duration, err := searchIndexWithLoop(arr, v); err != nil {
			fmt.Printf("Error during execution: %v", err)
			os.Exit(1)
		} else {
			fmt.Printf("Found with Loop %d at index %d. Duration: %v.\n", v, index, duration)
		}
		if index, duration, err := searchIndexWithTree(arr, v); err != nil {
			fmt.Printf("Error during execution: %v", err)
			os.Exit(1)
		} else {
			fmt.Printf("Found with Tree %d at index %d. Duration: %v.\n", v, index, duration)
		}
		fmt.Println("----")
	}
}
