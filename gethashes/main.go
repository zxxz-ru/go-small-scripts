package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	usage = `usage: gethashes 123456 3
  123456 - hashes seed, not less than 6 digits,
  3 - number of iterations, positive integer > 0`
	pause = 3
)

type Hash struct {
	Seed string `json:"number"`
	Hash string `json:"hash"`
}

func makeHashJSON(n string) []byte {
	hash := sha256.Sum256([]byte(n))
	hashStr := fmt.Sprintf("%x", hash)
	hs := Hash{n, hashStr}
	res, err := json.Marshal(hs)
	if err != nil {
		fmt.Println("can not Marshal json")
	}
	return res
}
func makePool(n int) string {
	var res string
	for _, v := range rand.Perm(n * 100) {
		res = res + fmt.Sprintf("%d", v)
	}
	return res
}
func isUnic(s string, sl []string) bool {
	res := false
	for _, str := range sl {
		res = strings.Compare(s, str) == 0
		if res {
			break
		}
	}
	return !res
}
func makeGenerator(n int) func() string {
	var unics []string
	pool := makePool(n)
	return func() string {
	Loop:
		try := pool[:4]
		if isUnic(try, unics) {
			unics = append(unics, try)
			pool = pool[1:]
			return try
		}
		goto Loop
	}
}

func main() {
	// make sure user input comply to requerements
	if len(os.Args) < 3 {
		fmt.Println(usage)
		os.Exit(1)
	}
	if len(os.Args[1]) < 6 {
		fmt.Println(usage)
		os.Exit(1)
	}
	// get variables to work with
	number := os.Args[1]
	laps, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(usage)
		os.Exit(1)
	}
	if laps <= 0 {
		fmt.Println(usage)
		os.Exit(1)
	}

	gen := makeGenerator(laps)
	for i := 0; i < laps; i++ {
		replacer := strings.NewReplacer(number[len(number)-4:], gen())
		number = replacer.Replace(number)
		fmt.Println(number)
	}
	res := makeHashJSON(os.Args[1])
	fmt.Println(string(res))
}
