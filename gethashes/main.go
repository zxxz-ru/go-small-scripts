package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	usage = `usage: gethashes 123456 3
  123456 - hashes seed, not less than 6 digits,
  3 - number of iterations, positive integer > 0`
	pause = 3
	conf  = "config.json"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

// get string and create hash using param n as seed
// it's not realy SOLID, but for simplicity decided not
// to break it in two funcs
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

// send created hash json to chan with provided param number,
// create random number for laps times sleep for 3 seconds
// send created json to chan in loop.
func makeHashes(number string, laps int, ch chan []byte) {
	defer close(ch)
	gen := makeGenerator(laps)
	ch <- makeHashJSON(number)
	for i := 1; i < laps; i++ {
		time.Sleep(pause * time.Second)
		replacer := strings.NewReplacer(number[len(number)-4:], gen())
		number = replacer.Replace(number)
		ch <- makeHashJSON(number)
	}
}

// send hashes to web script using websocket connection
// number and laps are parameters passed with command line
func sendHashes(conn *websocket.Conn, number string, laps int) {
	ch := make(chan []byte, laps)
	go makeHashes(number, laps, ch)
	for hash := range ch {
		w, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			fmt.Println(err.Error())
		}
		// fmt.Printf("%s\n", hash)
		w.Write(hash)
		w.Close()
	}
	conn.Close()
	fmt.Println("Done, exit now.")
	os.Exit(0)
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
	// config, err := loadConfig()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "hashes.html")
	})

	http.HandleFunc("/hashes", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		go sendHashes(conn, number, laps)
	})

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
