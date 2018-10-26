package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	usage = `usage: gethashes 123456 3
  123456 - hashes seed, not less than 6 digits,
  3 - number of iterations, positive integer > 0`
	pause = 3
	conf  = "config.json"
)

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}
type WebSocketConfig struct {
	Location string `json:"location"`
	Origin   string `json:"origin"`
	Port     string `json:"port"`
}
type Config struct {
	Redis RedisConfig     `json:"redis"`
	WS    WebSocketConfig `json:"websocket"`
}

type Hash struct {
	Seed string `json:"number"`
	Hash string `json:"hash"`
}

func loadConfig() (res Config, err error) {
	file, err := os.Open(conf)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		err = fmt.Errorf("error in loadConfig %s", err)
		return res, err
	}
	parser := json.NewDecoder(file)
	err = parser.Decode(&res)
	return res, nil
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

// create long string of random digits
func makePool(n int) string {
	var res string
	for _, v := range rand.Perm(n * 100) {
		res = res + fmt.Sprintf("%d", v)
	}
	return res
}

// compares range of srings from slice
// return true if string is unic
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

// create generator function which returns string of
// four unic digits
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

func makeHashes(number string, laps int, ch chan []byte) {
	gen := makeGenerator(laps)
	ch <- makeHashJSON(number)
	defer close(ch)
	for i := 1; i < laps; i++ {
		replacer := strings.NewReplacer(number[len(number)-4:], gen())
		number = replacer.Replace(number)
		ch <- makeHashJSON(number)
	}
}

func makeHashesServer(ws *websocket.Conn, number string, laps int) func(w *websocket.Conn) {
	return func(ws *websocket.Conn) {
		ch := make(chan []byte, laps)
		counter := 0
		go makeHashes(number, laps, ch)
		for hash := range ch {
			if counter > 0 {
				time.Sleep(pause * time.Second)
			}
			fmt.Printf("%s", hash)
			counter = counter + 1
			// ws.Message.Send(hash)
		}
	}
}

func main() {
	// make sure user input comply to requerements
	if len(os.Args) < 3 {
		fmt.Println(usage)
		os.Exit(1)
	}
	// for now presume Args[1] is digit
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
	// load configs
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ws, err := websocket.Dial(config.WS.Location, "", config.WS.Origin)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	wsServer := makeHashesServer(ws, number, laps)
	http.Handle(config.WS.Location+"/", websocket.Handler(wsServer))
	err = http.ListenAndServe(config.WS.Location, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
