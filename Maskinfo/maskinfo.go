package maskinfo

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo"
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
	Type string `json:"type"`
}

type Sales struct {
	Code       string `json:"code"`
	CreatedAt  string `json:"created_at"`
	RemainStat string `json:"remain_stat"`
	StockAt    string `json:"stock_at"`
}

func Drugstore(c echo.Context) error {
	name := c.Param("name")
	slice1 := []interface{}{}
	ch := make(chan interface{}, 500)
	var wg sync.WaitGroup
	for i := 1; i <= 54; i++ {
		wg.Add(1)
		go FindDrugstore(name, i, &wg, ch)
	}
	wg.Wait()
	close(ch)
	for i := range ch {
		slice1 = append(slice1, i)
	}
	return c.JSON(http.StatusOK, slice1)
}

func FindDrugstore(name string, page int, wg *sync.WaitGroup, ch chan interface{}) {
	stringpage := strconv.Itoa(page)
	address := "https://8oi9s0nnth.apigw.ntruss.com/corona19-masks/v1/stores/json?page=" + stringpage
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
			ch <- v
		}
	}
	wg.Done()
}

func Masks(c echo.Context) error {
	code := c.Param("code")
	slice1 := []interface{}{}
	ch := make(chan interface{}, 500)
	var wg sync.WaitGroup
	for i := 1; i <= 51; i++ {
		wg.Add(1)
		go FindMasks(code, i, &wg, ch)
	}
	wg.Wait()
	close(ch)
	for i := range ch {
		slice1 = append(slice1, i)
	}
	return c.JSON(http.StatusOK, slice1)
}

func FindMasks(code string, page int, wg *sync.WaitGroup, ch chan interface{}) {
	stringpage := strconv.Itoa(page)
	address := "https://8oi9s0nnth.apigw.ntruss.com/corona19-masks/v1/sales/json?page=" + stringpage
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
			ch <- v
		}
	}
	wg.Done()
}
