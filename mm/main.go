package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OsoianMarcel/claymore-go"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AutoGenerated struct {
	Data struct {
		Stats struct {
			MinerVersion   string `json:"miner_version"`
			RunningMinutes int    `json:"running_minutes"`
			EthReport      struct {
				TotalMhs       int `json:"total_mhs"`
				Shares         int `json:"shares"`
				RejectedShares int `json:"rejected_shares"`
				InvalidShares  int `json:"invalid_shares"`
				PoolSwitches   int `json:"pool_switches"`
				MhsPerGpu      []struct {
					Mhs int `json:"mhs"`
					Gpu int `json:"gpu"`
				} `json:"mhs_per_gpu"`
			} `json:"eth_report"`
			AltReport struct {
				TotalMhs       int `json:"total_mhs"`
				Shares         int `json:"shares"`
				RejectedShares int `json:"rejected_shares"`
				InvalidShares  int `json:"invalid_shares"`
				PoolSwitches   int `json:"pool_switches"`
				MhsPerGpu      []struct {
					Mhs int `json:"mhs"`
					Gpu int `json:"gpu"`
				} `json:"mhs_per_gpu"`
			} `json:"alt_report"`
			TempAndFanReports []struct {
				Temp int `json:"temp"`
				Fan  int `json:"fan"`
				Gpu  int `json:"gpu"`
			} `json:"temp_and_fan_reports"`
			Pools []string `json:"pools"`
		} `json:"stats"`
		Extra struct {
			HighestTemp struct {
				Temp int `json:"temp"`
				Fan  int `json:"fan"`
				Gpu  int `json:"gpu"`
			} `json:"highest_temp"`
		} `json:"extra"`
	} `json:"data"`
}

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

var Mhs, T, C []int
var Mhssum, Tmax, Tmin, Cmax, Cmin, Trmax, Trmin, Crmax, Crmin int
var count_message int = 0

func get_json_miner() {
	Tmax = 80
	Tmin = 0
	Cmax = 100
	Cmin = 0
	Mhssum = 0
	Mhs = nil
	T = nil
	C = nil

	mm := [3]string{"109.172.77.189:33332", "109.172.77.189:33333", "109.172.77.189:33337"}
	//mm := []string{"109.172.77.189:33332"}
	for i := 0; i < len(mm); i++ {
		cc := claymore.NewClient(mm[i])
		resp, err := cc.GetStats()
		if err != nil {
			// TODO: Handle error.
		}
		writer := bytes.NewBuffer([]byte{})
		extraResp := ExtraResponse{}
		statsResp := StatsResponse{resp, extraResp}
		json.NewEncoder(writer).Encode(DataResponse{statsResp})
		//fmt.Println(writer.String())
		writer1 := []byte(writer.String())
		var app = AutoGenerated{}
		err1 := json.Unmarshal(writer1, &app)
		if err1 != nil {
			log.Fatal("error")
		}
		//fmt.Println(app.Data.Stats.EthReport.MhsPerGpu)
		for _, row := range app.Data.Stats.EthReport.MhsPerGpu {
			//fmt.Println(row.Mhs)
			Mhs = append(Mhs, row.Mhs)
			Mhssum = Mhssum + row.Mhs
		}
		for _, row := range app.Data.Stats.TempAndFanReports {
			T = append(T, row.Temp)
			if row.Temp > Tmin {
				Trmax = row.Temp
				Tmin = row.Temp
			}
			if row.Temp < Tmax {
				Trmin = row.Temp
				Tmax = row.Temp
			}
			C = append(C, row.Fan)
			if row.Fan > Cmin {
				Crmax = row.Fan
				Cmin = row.Fan
			}
			if row.Fan < Cmax {
				Crmin = row.Fan
				Cmax = row.Fan
			}
		}
	}
	if Mhssum < 380000 || Trmax > 80 || Trmin < 40 || Crmax > 80 || Crmin < 40 {
		count_message = count_message + 1
		fmt.Println(count_message)
		if count_message > 9 {
			count_message = 0
			MakeRequest()
		}
	}

}

func MakeRequest() {
	var y string = time.Now().Format("15:04:05 01-02-2006") + " Mhs=" + strconv.Itoa(Mhssum) + " Tmax=" + strconv.Itoa(Trmax) + " Tmin=" + strconv.Itoa(Trmin) + " Cmax=" + strconv.Itoa(Crmax) + " Cmin=" + strconv.Itoa(Crmin)
	var yy string = ``
	for i := 0; i < len(Mhs); i++ {
		yy = yy + strconv.Itoa(Mhs[i]) + ` - ` + strconv.Itoa(T[i]) + ` - ` + strconv.Itoa(C[i]) + `\n`
	}
	data := []byte(`{
		"embeds": [
		  {
			"title": "MM",
			"description": "` + y + `",
			"color": 564300,
			"fields": [{"name": "Report","value": "` + yy + `"}]
		}
		]
	  }`)
	r := bytes.NewReader(data)
	resp, err := http.Post("https://discordapp.com/api/webhooks/642252309182808074/Rnx0v-0hKSuX-GYmNrpegEVjUsXm0I6K703L2G85lsp2kM-TmYqN3-zYcR4IuFAv_blh", "application/json", r)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	//log.Println(result)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, time.Now().Format("15:04:05 01-02-2006"), "\nMhs=", Mhssum, "\nTmax=", Trmax, "\nTmin=", Trmin, "\nCmax=", Crmax, "\nCmin=", Crmin, "\n")
	for i := 0; i < len(Mhs); i++ {
		fmt.Fprintln(w, Mhs[i], T[i], C[i])
	}

	//for _, v := range Mhs {
	//	fmt.Fprintln(w, v)
	//}
}
func get_json_miner_timer() {
	for {
		get_json_miner()
		time.Sleep(60 * time.Second)
	}
}

func main() {
	go get_json_miner_timer()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8009", nil))
}
