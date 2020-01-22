package main

import (
	"fmt"
	"github.com/OsoianMarcel/claymore-go"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TroloLO!\n")
	cc := claymore.NewClient("109.172.77.189:13333")
	fmt.Fprintln(cc.GetStats())
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
