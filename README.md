<div align="center">
  
![GitHub contributors](https://img.shields.io/github/contributors/jjmin321/MaskInfo_api)
![GitHub forks](https://img.shields.io/github/forks/jjmin321/MaskInfo_api?label=Forks)
![GitHub stars](https://img.shields.io/github/stars/jjmin321/MaskInfo_api?style=Stars)
![GitHub issues](https://img.shields.io/github/issues-raw/jjmin321/MaskInfo_api)
[![Go Report Card](https://goreportcard.com/badge/github.com/jjmin321/MaskInfo_api)](https://goreportcard.com/badge/github.com/jjmin321/MaskInfo_api)

</div>

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

## 🦍 채널을 통해 일치하는 값을 전송하는 코드 (마지막 업데이트 - 코드 수정)
```go
// 고루틴(쓰레드)가 돌아가는 도중 채널을 통해 일치하는 값이 나왔을 때 반환이 아닌 전송함으로써 고루틴(쓰레드)간의 교착 상태가 발생하지 않음 
for _, v := range result.StoreInfos {
		if v.Name == name {
			// for i := 1; i <= 51; i++ {
			// 	go masks(v.Code, v.Addr, v.Name, i, wg)
			// }
			// jsonBytes, _ := json.Marshal(v)	// 마지막 업데이트 - 코드 삭제 
			// jsonString := string(jsonBytes)  // 마지막 업데이트 - 코드 삭제 
			// ch <- jsonString    //채널을 통한 값 전송 // 마지막 업데이트 - 코드 삭제
			ch <- v // JSON으로 다시 디코딩 할 필요 없이 인터페이스(모든 타입 담을 수 있는 타입) 형으로 전송 
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
			slice1 = append(slice1, i)	// append함수는 첫번째 인자 + 두번 째 인자 의 값을 반환 // 첫 번쨰 인자와 두 번째 인자 타입 같아야함 
		}
```

## 📔 인코딩을 하여 슬라이스를 다시 JSON 값으로 변경하여 반환 (마지막 업데이트 - 코드 수정)
```go
// jsonBytes, _ := json.Marshal(slice1)	// 코드 삭제 
// return c.JSONBlob(http.StatusOK, jsonBytes)	// 코드 삭제 
return c.JSON(http.StatusOk, slice1) // @api /drugstore
return c.JSON(http.StatusOk, <-ch)	// @api /masks
```

## 🍭 어떤 시간에 사용자가 어떤 웹 브라우저를 통하여 어떤 api 를 사용했는지 로그 출력 
```go
// Echo 프레임워크에서 코드를 읽어보면 내가 포맷 형식을 설정할 수 있음.
e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `Time = ${time_rfc3339}], Status = ${status}, Method = ${method}, Uri = ${uri}, User_agent = ${user_agent}, Error = ${error},  Remote_ip = ${remote_ip}` + "\n",
}))
```
<img width="1440" alt="Screen Shot 2020-03-23 at 15 18 22" src="https://user-images.githubusercontent.com/52072077/77288466-3ccf8900-6d1b-11ea-8a9b-ac2fc808e218.png">

# API - FUNCTION

### @GET /drugstore
- Request : name (QueryParam) [약국 이름]
- Response : 약국 이름, 약국 주소, 약국 코드 (JSON)
<img width="1440" alt="Screen Shot 2020-03-23 at 15 23 05" src="https://user-images.githubusercontent.com/52072077/77288467-3e00b600-6d1b-11ea-974c-e442b0c75b77.png">

### @GET /masks
- Request : code (QueryParam) [약국 코드]
- Response : 해당 약국 마스크 재고량이 갱신된 시간, 마스크 재고량, 최근 재고 갱신 시간 (JSON)

<img width="1440" alt="Screen Shot 2020-03-23 at 15 27 30" src="https://user-images.githubusercontent.com/52072077/77288470-3f31e300-6d1b-11ea-8199-2e30ab041f08.png">


# 2020/3/23 - 마지막 코드 수정 
1. JSON 값이 들어있는 웹페이지 1 ~ 50 개를 고루틴(쓰레드)을 통해 들어간 후 사용자가 검색한 값만 찾아 채널을 통해 비동기 전송 : 속도 증가 
2. waitgroup, channel을 통해 2단 동기화 : waitgroup을 통해 1 ~ 50 개의 고루틴(쓰레드)이 모두 끝나기 전까지 기다림, channel을 통해 모든 비동기 전송 값을 받으면 클라이언트에 값을 반환
3. channel의 버퍼 크기 수정 : 채널의 버퍼 크기를 50으로 해놨더니 신세게약국(80개 존재)을 검색하면 오류가 발생했다. 채널 버퍼 크기를 50에서 500으로 수정해줬다.
4. JSON으로 다시 인코딩할 필요 없이 slice 와 채널 값 그대로 반환 : Echo 프레임워크에서 제공하는 JSON 반환은 어떤 값이든 자동으로 JSON으로 바꿔줌
5. 채널 값을 반환하는 mask api 또한 slice로 변환 : 채널에 값이 들어오지 않으면 무한 로딩되는 것을 개선 