package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

//https://affiliate.dmm.com/api/v3/itemlist.html

func TestAPI(t *testing.T) {

	apiID := os.Getenv("AV_API_ID")
	offset := 1
	layout := "2006-01-02T00:00:00"

	var params = map[string]string{
		"service":      "digital",
		"output":       "json",
		"offset":       strconv.Itoa(offset),
		"hits":         "20",
		"api_id":       apiID,
		"affiliate_id": "10278-996",
		"article":      "actress",
		"gte_date":     time.Now().AddDate(0, 0, -7).Format(layout),
	}

	data := ClientRequest(AV_PRODUCT_API_URL, nil, params)

	// res := &AvAPIResponse{}
	// err := json.Unmarshal(data, res)
	// if err != nil {
	// 	fmt.Println("av error: ", err)
	// }

	// if res.Result.ResultCount == 0 {
	// 	fmt.Println("retry: ", err)
	// 	return
	// }

	// index := RandomInt(res.Result.ResultCount-1, 0)
	// actress := res.Result.Actress[index]

	fmt.Println("data: ", data)
}
