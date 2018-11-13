package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
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

var pre_vol int
var pre_ud int
var pre_future []Futures

func main() {
	var detail *bool
	detail = flag.Bool("detail", false, "detail")
	var s_time *string
	s_time = flag.String("time", "auto", "day or night or auto")
	var help *bool
	help = flag.Bool("help", false, "help")
	flag.Parse()

	if *help {
		fmt.Println("--detail for bool default false\n--time for day/night/auto default auto")
		return
	}

	fmt.Println("isOpen:", isOpen())
	if !isOpen() {
		return
	}
	for true {
		url := getURL(*s_time)
		futrues, err := fetch(url)

		if err != nil {
			fmt.Println(err)
			return
		}

		if len(pre_future) > 0 {
			if futrues[0].Volume == pre_future[0].Volume {
				continue
			}
		}

		if *detail {
			printDetail(futrues)
		} else {
			printBrief(futrues)
		}
		time.Sleep(30 * time.Second)
		pre_future = futrues
	}
}

func printDetail(futrues []Futures) {
	for _, future := range futrues {
		s := fmt.Sprintf("c:%s, p:%s, vol:%s, name:%s, range:%s", future.Contract, future.Price, future.Volume, future.ContractName, future.Updown)
		fmt.Println(s)
	}
}

func StrToInt(s string) int {
	num, err := strconv.Atoi(strings.Replace(s, ",", "", -1))
	if err == nil {
		return num
	}
	return -1
}

func printBrief(futrues []Futures) {
	future := futrues[0]
	vol := StrToInt(future.Volume)
	price := StrToInt(future.Price)
	updown := StrToInt(future.Updown)
	s := fmt.Sprintf("[%s] p:%d, vol:%d, range:%d, vol_dif:%d, ud_dif:%d",
		time.Now().Format("2006-01-02 15:04:05"),
		price,
		vol,
		updown,
		vol-pre_vol,
		updown-pre_ud)
	pre_vol = vol
	pre_ud = updown
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
		if isDay() {
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
	h := float64(t.Hour())
	m := float64(t.Minute())
	s := float64(t.Second())

	t_in_min := h*60 + m + s/60

	// 05:00 = 300
	// 08:45 = 525
	// 13:45 = 825
	// 15:00 = 900

	if t_in_min < 525 && t_in_min > 300 {
		return false
	}
	if t_in_min > 825 && t_in_min < 900 {
		return false
	}
	return true
}

func isDay() bool {
	t := time.Now()
	if t.Weekday() == 0 || t.Weekday() == 6 {
		return false
	}
	h := t.Hour()
	m := t.Minute()
	s := t.Second()

	t_in_min := h*60 + m + s/60.0

	// 05:00 = 300
	// 08:45 = 525
	// 13:45 = 825
	// 15:00 = 900

	if t_in_min >= 525 && t_in_min <= 825 {
		return true
	}
	return false
}

// 1
// [  {    "futvol": "425,705",    "optvol": "1,114,138"  }]

// id 3 是每分鐘資料
// [  {    "time": "0845",    "price": "9607"  }

// id 13 個股
// 14 k棒 （似乎不止一天）
