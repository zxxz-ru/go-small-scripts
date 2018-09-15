package main
import(
"./sum"
"fmt"
"net/http"
"os"
)

func DoubleHandler(w http.ResponseWriter, r *http.Request){
    _, err := fmt.Fprintf(w, "This is temporary response.\n")
    if err != nil{
    os.Exit(1)
    }

}
func main(){
fmt.Printf("Doule of 2 is %v\n", sum.Double(2))
http.HandleFunc("/double/", DoubleHandler)
err := http.ListenAndServe(":8080", nil)
if err == nil{
panic("Panic in main can not srart server")
}
}
