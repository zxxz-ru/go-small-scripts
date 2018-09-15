package main
import(
"log"
"time"
"net/http"
)

type TimeHandler struct {
format string
}

func (th *TimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
    tm := time.Now().Format(th.format)
    w.Write([]byte("Current time is: " + tm))
}
func main() {
    mux := http.NewServeMux()
    rt := http.RedirectHandler("http://google.com", 302)
    mux.Handle("/foo", rt)
    th := &TimeHandler{format:time.RFC1123}
    mux.Handle("/time", th)
    log.Println("Start server")
    http.ListenAndServe(":8080", mux)
}
