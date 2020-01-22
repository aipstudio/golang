package main

import (
	"fmt"
	"github.com/OsoianMarcel/claymore-go"
	"encoding/json"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
}

type ExtraResponse struct {
	HighestTemp claymore.TempAndFanReport `json:"highest_temp"`
}

type StatsResponse struct {
	Stats claymore.StatsModel `json:"stats"`
	Extra ExtraResponse       `json:"extra"`
}

type Person struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
}

func get_stats() {
	cc := claymore.NewClient("109.172.77.189:13333")
	//stats, err := cc.GetStats()
	in := `{"firstName":"John","lastName":"Dow"}`
	fmt.Println(cc.GetStats())
	rawIn := json.RawMessage(in)
	bytes, err := rawIn.MarshalJSON()
    if err != nil {
        panic(err)
	}
	var p Person
    err = json.Unmarshal(bytes, &p)
    if err != nil {
        panic(err)
    }
    fmt.Printf(p.FirstName)
	//fmt.Println(statsResp.Stats.EthReport.TotalMhs)
}

func main() {
	get_stats()
	/*cc := claymore.NewClient("109.172.77.189:13333")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		stats, err := cc.GetStats()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}

		extraResp := ExtraResponse{}

		if ht, err := stats.GetHighestTemp(); err == nil {
			extraResp.HighestTemp = ht
		}
		statsResp := StatsResponse{stats, extraResp}
		//json.NewEncoder(w).Encode(DataResponse{statsResp})
		//json.NewEncoder(w).Encode(DataResponse{statsResp.Stats.EthReport.TotalMhs})
		fmt.Println(statsResp.Stats.EthReport.TotalMhs)
		q = statsResp.Stats.EthReport.Shares
	})
	http.ListenAndServe(":8080", nil)
	*/
}
