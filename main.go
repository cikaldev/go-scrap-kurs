package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Bank struct {
	Name      string    `json:"bank"`
	TimeStamp time.Time `json:"timestamp"`
	Data      []Dataset `json:"data"`
}

type Dataset struct {
	Currency string  `json:"currency"`
	Buy      float32 `json:"buy"`
	Sell     float32 `json:"sell"`
}

var bank Bank

func normalize(str string) string {
	return strings.Replace(str, ",", "", -1)
}

func normalizeDot(str string) string {
	return strings.Replace(string(strings.Replace(str, ".", "", -1)), ",", ".", -1)
}

func httpGet(targetUri string) (goquery.Document, error) {
	res, err := http.Get(targetUri)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return *doc, err
}

func getBCA() string {
	doc, err := httpGet("https://www.bca.co.id/id/informasi/kurs")
	if err != nil {
		log.Fatalln(err)
	}

	bank.Name = "BCA - Bank Central Asia"
	bank.TimeStamp = time.Now()
	rows := make([]Dataset, 0)

	pageCounts := doc.Find(".m-table-kurs tbody tr")
	pageCounts.Each(func(_ int, tr *goquery.Selection) {
		td := new(Dataset)
		td.Currency = tr.Find("td:first-child span p").Text()
		buy, _ := strconv.ParseFloat(normalizeDot(tr.Find("td:nth-child(2) p").Text()), 32)
		td.Buy = float32(buy)
		sell, _ := strconv.ParseFloat(normalizeDot(tr.Find("td:nth-child(3) p").Text()), 32)
		td.Sell = float32(sell)
		rows = append(rows, *td)
	})
	bank.Data = rows
	hasil, _ := json.Marshal(bank)
	return string(hasil)
}

func getBI() string {
	doc, err := httpGet("https://www.bi.go.id/id/statistik/informasi-kurs/transaksi-bi/default.aspx")
	if err != nil {
		log.Fatalln(err)
	}

	bank.Name = "BI - Bank Indonesia"
	bank.TimeStamp = time.Now()
	rows := make([]Dataset, 0)

	pageCounts := doc.Find(".table-lg tbody tr")
	pageCounts.Each(func(_ int, tr *goquery.Selection) {
		td := new(Dataset)
		td.Currency = strings.Trim(tr.Find("td:first-child").Text(), " ")
		buy, _ := strconv.ParseFloat(normalizeDot(tr.Find("td:nth-child(4)").Text()), 32)
		td.Buy = float32(buy)
		sell, _ := strconv.ParseFloat(normalizeDot(tr.Find("td:nth-child(3)").Text()), 32)
		td.Sell = float32(sell)
		rows = append(rows, *td)
	})
	bank.Data = rows
	hasil, _ := json.Marshal(bank)
	return string(hasil)
}

func getBNI() string {
	doc, err := httpGet("https://www.bni.co.id/id-id/beranda/informasivalas")
	if err != nil {
		log.Fatalln(err)
	}

	bank.Name = "BNI - Bank Negara Indonesia"
	bank.TimeStamp = time.Now()
	rows := make([]Dataset, 0)

	pageCounts := doc.Find("#dnn_ctr3510_BNIValasInfoView_divBankNotes table tbody tr")
	pageCounts.Each(func(_ int, tr *goquery.Selection) {
		td := new(Dataset)
		td.Currency = tr.Find("td:first-child").Text()
		buy, _ := strconv.ParseFloat(normalizeDot(tr.Find("td:nth-child(2)").Text()), 32)
		td.Buy = float32(buy)
		sell, _ := strconv.ParseFloat(normalizeDot(tr.Find("td:nth-child(3)").Text()), 32)
		td.Sell = float32(sell)
		rows = append(rows, *td)
	})
	bank.Data = rows
	hasil, _ := json.Marshal(bank)
	return string(hasil)
}

func getMEGA() string {
	doc, err := httpGet("https://www.bankmega.com/id/bisnis/treasury/")
	if err != nil {
		log.Fatalln(err)
	}

	bank.Name = "MEGA - Bank Mega"
	bank.TimeStamp = time.Now()
	rows := make([]Dataset, 0)

	pageCounts := doc.Find("table tbody tr")
	pageCounts.Each(func(_ int, tr *goquery.Selection) {
		td := new(Dataset)
		td.Currency = tr.Find("td:first-child").Text()
		buy, _ := strconv.ParseFloat(normalize(tr.Find("td:nth-child(2)").Text()), 32)
		td.Buy = float32(buy)
		sell, _ := strconv.ParseFloat(normalize(tr.Find("td:nth-child(3)").Text()), 32)
		td.Sell = float32(sell)
		rows = append(rows, *td)
	})
	bank.Data = rows
	hasil, _ := json.Marshal(bank)
	return string(hasil)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rp := html.EscapeString(r.URL.Path)

		log.Printf("GET: %s", rp)
		w.Header().Set("Content-Type", "application/json")
		if rp == "/" {
			fmt.Fprintf(w, `{"code": 200, "message": "Go Scrap Kurs"}`)
		} else {
			fmt.Fprintf(w, `{"code": 404, "message": "Page Not Found", "path": %q}`, rp)
		}
	})

	mux.HandleFunc("/bca", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET: /bca")
		w.Header().Set("Content-Type", "application/json")
		bca := getBCA()
		fmt.Fprint(w, bca)
	})

	mux.HandleFunc("/bi", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET: /bi")
		w.Header().Set("Content-Type", "application/json")
		bi := getBI()
		fmt.Fprint(w, bi)
	})

	mux.HandleFunc("/bni", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET: /bni")
		w.Header().Set("Content-Type", "application/json")
		bni := getBNI()
		fmt.Fprint(w, bni)
	})

	mux.HandleFunc("/mega", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET: /mega")
		w.Header().Set("Content-Type", "application/json")
		mega := getMEGA()
		fmt.Fprint(w, mega)
	})

	log.Println("server started")
	log.Fatal(http.ListenAndServe(":2021", mux))
}
