# 😷MaskInfo_API - 3/17/2020 ~ 3/19/2020 😷

🌟약국 이름을 검색하면 약국 이름이 들어간 모든 약국 코드와 위치를 알려주며, 약국 코드를 검색하면 그에 맞는 마스크 재고량과 재고 갱신 시간 등을 알려주는 API Server 🌟

## 📦 사용한 패키지
```go
    "encoding/json"
	"fmt"
	"net/http"
	"strconv"
    "sync"
    "github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
```

## 🖥️ 사용한 서버 프레임워크
https://echo.labstack.com [GO - echo]

## 🧚‍♂️모든 사이트를 고루틴을 통해 이동하는 코드
```go
// 전국 약국 정보들이 들어간 54개의 사이트에 고루틴(쓰레드)를 통해 동시에 접속
for i := 1; i <= 54; i++ {
			wg.Add(1)
			go api.Drugstore(name, i, &wg, ch)
		}
```

## 🌱 사이트의 코드를 읽어오는 코드 
```go
// 사이트 코드를 모두 읽어온 후 defer 함수를 통해 작업이 끝나면 자동으로 사이트 접속 종료
resp, err := http.Get(address)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
```

## 📃 사이트의 코드에서 원하는 값을 필터링하는 코드 
```go
// 구조체 생성 후 JSON 태그를 통해 원하는 값만 GO value로 변경
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
```

## 🦍 채널을 통해 일치하는 값을 전송하는 코드 
```go
// 고루틴(쓰레드)가 돌아가는 도중 채널을 통해 일치하는 값이 나왔을 때 반환이 아닌 전송함으로써 고루틴(쓰레드)간의 교착 상태가 발생하지 않음 
for _, v := range result.StoreInfos {
		if v.Name == name {
			// for i := 1; i <= 51; i++ {
			// 	go masks(v.Code, v.Addr, v.Name, i, wg)
			// }
			jsonBytes, _ := json.Marshal(v)
			jsonString := string(jsonBytes)
			ch <- jsonString    //채널을 통한 값 전송 
		}
	}
```

## ⛓️ waitgroup을 통해 동기화 
```go
// wg.Add(1) = 작업이 시작할 때마다 wg 값을 1 증가시킴 
// wg.Done() = 작업이 끝날때마다 wg 값을 1 감소시킴
// wg.Wait() = wg 값이 0이 되기 전까지 다음 코드로 넘어가지 않음 
wg.Add(1)
wg.Done()
wg.Wait()
```

## 🎮 채널을 통해 들어온 값을 순회하며 슬라이스에 추가 시킴 
```go
//append를 통해 기존의 배열 값에 추가된 채널 값을 합침 
for i := range ch {
			slice1 = append(slice1, i)
		}
```

## 📔 인코딩을 하여 슬라이스를 다시 JSON 값으로 변경하여 반환 
```go
//json.Marshal() 을 통해 인코딩
jsonBytes, _ := json.Marshal(slice1)
return c.JSONBlob(http.StatusOK, jsonBytes)
```

## API - FUNCTION

### @GET /drugstore
- Request : name (QueryParam) [약국 이름]
- Response : 약국 이름, 약국 주소, 약국 코드 (JSON)

### @GET /masks
- Request : code (QueryParam) [약국 코드]
- Response : 해당 약국 마스크 재고량이 갱신된 시간, 마스크 재고량, 최근 재고 갱신 시간 (JSON)