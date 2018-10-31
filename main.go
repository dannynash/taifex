package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

type Futures struct {
	Contract     string `json:"contract"`
	Price        string `json:"price"`
	Volume       string `json:"ttlvol"`
	ContractName string `json:"contractName"`
	Updown       string `json:"updown"`
}

const (
	URL = "http://www.taifex.com.tw/cht/quotesApi/getQuotes"
)

func main() {
	var detail *bool
	detail = flag.Bool("detail", false, "detail")
	var time *string
	time = flag.String("time", "auto", "day or night or auto")
	var help *bool
	help = flag.Bool("help", false, "help")
	flag.Parse()

	if *help {
		fmt.Println("--detail for bool default false\n--time for day/night/auto default auto")
		return
	}

	url := getURL(*time)
	futrues, err := fetch(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("isOpen:", isOpen())
	if *detail {
		printDetail(futrues)
	} else {
		printBrief(futrues)
	}
}

func printDetail(futrues []Futures) {
	for _, future := range futrues {
		s := fmt.Sprintf("c:%s, p:%s, vol:%s, name:%s, range:%s", future.Contract, future.Price, future.Volume, future.ContractName, future.Updown)
		fmt.Println(s)
	}
}

func printBrief(futrues []Futures) {
	future := futrues[0]
	s := fmt.Sprintf("p:%s, vol:%s, range:%s", future.Price, future.Volume, future.Updown)
	fmt.Println(s)
}

func fetch(url string) (futrues []Futures, err error) {
	resp, _, errs := gorequest.New().
		Get(url).
		EndStruct(&futrues)

	if resp.StatusCode != 200 {
		err = errors.New("fetch failed")
	}
	if errs != nil {
		return
	}
	return
}

func getURL(time string) (url string) {
	switch time {
	case "day":
		return fmt.Sprintf("%s?objId=2", URL)
	case "night":
		return fmt.Sprintf("%s?objId=12", URL)
	default:
		if isOpen() {
			return fmt.Sprintf("%s?objId=2", URL)
		}
		return fmt.Sprintf("%s?objId=12", URL)
	}
}

func isOpen() bool {
	t := time.Now()
	if t.Weekday() == 0 || t.Weekday() == 6 {
		return false
	}
	h := t.Hour()

	if h < 9 || h >= 14 {
		return false
	}
	return true
}

// 1
// [  {    "futvol": "425,705",    "optvol": "1,114,138"  }]

// id 3 是每分鐘資料
// [  {    "time": "0845",    "price": "9607"  }

// id 13 個股
// 14 k棒 （似乎不止一天）
