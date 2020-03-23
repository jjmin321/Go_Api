package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Result struct {
	Count      int         `json:"count"`
	Page       string      `json:"page"`
	StoreInfos []StoreInfo `json:"storeInfos"`
	Sales      []Sales     `json:"sales"`
}

type StoreInfo struct {
	Code string `json:"code"`
	Addr string `json:"addr"`
	Name string `json:"name"`
}

type Sales struct {
	Code       string `json:"code"`
	CreatedAt  string `json:"created_at"`
	RemainStat string `json:"remain_stat"`
	StockAt    string `json:"stock_at"`
}

// 약국 이름을 받아서 약국 코드번호 출력
func Drugstore(name string, page int, wg *sync.WaitGroup, ch chan interface{}) {
	//GET 호출
	stringpage := strconv.Itoa(page)
	address := "https://8oi9s0nnth.apigw.ntruss.com/corona19-masks/v1/stores/json?page=" + stringpage
	// fmt.Println(address)
	resp, err := http.Get(address)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var result Result
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&result)

	for _, v := range result.StoreInfos {
		if v.Name == name {
			// for i := 1; i <= 51; i++ {
			// 	go masks(v.Code, v.Addr, v.Name, i, wg)
			// }
			// jsonBytes, _ := json.Marshal(v)
			// jsonString := string(jsonBytes)
			ch <- v
		}
	}
	wg.Done()
}

func Masks(code string, page int, wg *sync.WaitGroup, ch chan interface{}) {
	stringpage := strconv.Itoa(page)
	address := "https://8oi9s0nnth.apigw.ntruss.com/corona19-masks/v1/sales/json?page=" + stringpage
	// fmt.Println(address)
	resp, err := http.Get(address)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var result Result
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&result)

	for _, v := range result.Sales {
		if v.Code == code {
			fmt.Println(v.Code, v.CreatedAt, v.RemainStat, v.StockAt)
			// jsonBytes, _ := json.Marshal(v)
			// jsonString := string(jsonBytes)
			ch <- v
		}
	}
	wg.Done()
}
