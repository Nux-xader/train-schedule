package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

var agents = []string{
	"nitaratutourandtravel", "tiketkita", "liontravel", "viktoragentiket",
	"tomtourtravel", "yazittour", "tiketkufbtravel", "najwatravel",
	"amanahtravel", "oketrip", "cahayatravel", "toptravel", "mars",
	"totitravel", "tiketjalan", "buroq", "denvitiket", "pergicoid",
	"elbirunitravel", "andratravel", "tanpabatas", "gratiaholiday", "humayora",
	"travelindo", "dewatiket", "balibellatour", "niasticketing", "latikettnt",
	"lancongtravel", "mawantravel", "callystajek", "hasanahtravel",
	"apmikimmdopay", "gudangtiket", "sahabattravel", "tiketok", "tiketmurah99",
	"rajatiketnganjuk", "pergi", "indoagent", "umatourtravel",
	"ptyunapaymentindonesia", "safaritravell", "tiketjr", "eko",
	"sundaymanagement", "salondotravel", "hanstourtravel", "indahtourtravel",
	"paradiso", "nca", "aerotiketing", "susantravel", "sinergitravelindonesia",
	"travelintrips", "sagotourtravel", "aswtourtravel", "kencanatiketmurah",
	"tdtravel", "samuderatnt", "griyatiket", "sahabatiket", "hardilombok",
	"cahayalomboktravel", "nusaya", "cahaya", "ALESHAQUEEN", "louisatourtravel",
	"uwaistravel", "luluticketing", "aflahtour", "tjtiket", "citra55",
	"LINTASWISATA", "universaltravelindo", "rendo21travel",
	"ranszizaflightindonesia", "tourwisataloka", "grahatravelindonesia",
	"alkhairaattiket", "bamstravell", "easytrip", "hashinahtravel", "airatravel",
	"lawu", "haris", "joshuatourandtravel", "jintotravelindo", "vini",
	"jennitiket", "loketku", "nawawitravel", "vinstourntravel", "grahatiket",
	"tiketid", "natravel", "travelbanua", "jecktravel", "sunfelix", "yakinitravel",
	"alkahfitravelindo", "razdatravel", "mitrakelana", "barokahtiket",
	"toptraveling", "bostiket", "bandartiket", "arumtourtravel", "jektravel",
	"nurliwunto", "toscatiketmurah", "rexons", "travelselvyana", "kusumatravel",
	"blitiket", "khusustiketpesawat", "mitramodernutama", "bolangtourtravel",
	"raftiket", "poojaseratravelindo", "warungku", "biznet", "nadhifpay",
	"skylandinklusiutama", "smtourtravel", "gracia", "kiostiket", "wirawiri",
	"globalwonosobo", "airangga", "melancongwisatatravel", "flashmediatravel",
	"indotiket", "cindytravel", "TIKETTRAVEL", "nhasrultravel", "alfarizqitravel",
	"arrazzaq", "wanderlustravel", "tigabersaudaragrouptiketmurahtangerang",
	"transborneotourdantravel", "eticketindonesia", "lintastiket",
	"SutaNusantaraIntermedia", "jalubatravel", "mrtravel", "lsg", "travelisans",
	"pesantiket", "existicket", "bidaratravel", "bhaktitourandtravel", "keihintix",
	"visitix", "pison", "wahanamulyatiket", "jihantravel", "premieretravelindo",
	"globaltravel", "pelangi", "eaitravelindo", "firmanmitraholiday", "norton",
	"tiketkuinfo", "mamindotravel", "twins", "vilovertiketonline",
	"randujavatravel", "salmatravel", "habirtourdantravel", "safiratravel",
	"alltictravelin", "agentiketresmi", "bookingaja", "lomboktiket", "ahtcota",
	"fmtour", "adrtravel", "manogar", "alfazza", "restujubata",
	"darulabidintourtravel", "gracetravel", "agenmila", "baratransport",
	"paradisonesia", "aktharatravel", "wjtravel", "wgmtravel",
	"doucetbarakahtravel", "rajapelangi", "azkatravel", "imransaidtravelindo",
	"grahatiketindonesia", "aloha", "labirutour", "rayyahtravelsystem", "baubau",
	"azzmytrav", "mastravel", "manggalatravelagent", "luxury", "mytraveloka",
	"nauratravel", "initiketindonesia", "ndtravelagency", "yonotravel",
	"butiktravelnabire", "onasistravel", "mbokcikraktravel", "fezztravel",
	"otiket", "haifatravel", "akumauhotel", "landbowtour", "Ratnagemilang",
	"etctravel", "fstravel", "purnamatravel", "wtiket", "holipay", "jakstravel",
	"jelitatravel", "17011995",
}

var SpecialStations = map[string]func(org, dest, departDate string) (trains []Train, err error){
	"TGI": TravelokaTrainsSchedule,
	"LAR": TravelokaTrainsSchedule,
	"RJS": TravelokaTrainsSchedule,
	"SLS": TravelokaTrainsSchedule,
	"PYK": KaiWebSchedule,
}

type Train struct {
	Name  string
	Class string
	Stock int
}

type TrainApiResp struct {
	Data struct {
		Html string `json:"jkHtml"`
	} `json:"data"`
}

type TravelokaReq struct {
	Fields []any `json:"fields"`
	Data   struct {
		DepartureDate Date   `json:"departureDate"`
		ReturnDate    *Date  `json:"returnDate"`
		Destination   string `json:"destination"`
		Origin        string `json:"origin"`
		NumOfAdult    int    `json:"numOfAdult"`
		NumOfInfant   int    `json:"numOfInfant"`
		ProviderType  string `json:"providerType"`
		Currency      string `json:"currency"`
		TrackingMap   struct {
			UtmID              *string `json:"utmId"`
			UtmEntryTimeMillis int64   `json:"utmEntryTimeMillis"`
		} `json:"trackingMap"`
	} `json:"data"`
	ClientInterface string `json:"clientInterface"`
}

type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type TravelokaResp struct {
	Data struct {
		List []struct {
			Name     string `json:"trainBrandLabel"`
			Stock    string `json:"numSeatsAvailable"`
			Segments []struct {
				Summary struct {
					Class string `json:"subClass"`
				} `json:"productSummary"`
			} `json:"trainSegments"`
		} `json:"departTrainInventories"`
	} `json:"data"`
}

var monthID = []string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func GetTrainSchedule(org, dest, departDate string) (trains []Train, err error) {
	f, ok := SpecialStations[org]
	if ok {
		f_, ok := SpecialStations[dest]
		if ok {
			return f_(org, dest, departDate)
		}
		return f(org, dest, departDate)
	}

	f, ok = SpecialStations[dest]
	if ok {
		return f(org, dest, departDate)
	}

	return VelositaTrainsSchedule(org, dest, departDate)
}

func VelositaTrainsSchedule(org, dest, departDate string) (trains []Train, err error) {
	data := url.Values{
		"kereta_city_from_code": {org},
		"kereta_city_to_code":   {dest},
		"kereta_depart_date":    {departDate},
		"adult":                 {"1"},
		"infant":                {"0"},
	}

	agent := agents[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(agents))]
	req, err := http.NewRequest("POST", "https://velotiket.com/"+agent+"/tiket-kereta/jadwal-kai", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "id-ID,id;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.5 Mobile/15E148 Safari/604.1")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var trainResp = &TrainApiResp{}
	if err = json.Unmarshal(respBody, trainResp); err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(trainResp.Data.Html)))
	if err != nil {
		return
	}

	for _, node := range doc.Find(".train-result").Nodes {
		var train = Train{}

		train.Name = strings.Replace(goquery.NewDocumentFromNode(node).Find(
			"div.text-center.text-lg-start > strong",
		).Text(), "- ", "", 1)

		train.Class = strings.Trim(goquery.NewDocumentFromNode(node).Find(
			"div.col-12.col-6.col-sm-2.text-center.text-sm-end.text-dark.d-flex.flex-column.text-6.price > div > span",
		).Text(), "( )")

		train.Stock, err = strconv.Atoi(strings.Split(strings.TrimSpace(goquery.NewDocumentFromNode(node).Find(
			"div > div.col-12.col-6.col-sm-2.text-center.text-sm-end.text-dark.d-flex.flex-column.text-6.price > span.text-1.text-primary",
		).Text()), " ")[0])

		if err != nil {
			return
		}

		trains = append(trains, train)
	}

	return
}

func TravelokaTrainsSchedule(org, dest, departDate string) (trains []Train, err error) {
	departDateTime, _ := time.Parse("02-01-2006", departDate)
	reqStruct := TravelokaReq{
		Fields:          []any{},
		ClientInterface: "mobile",
	}

	reqStruct.Data.Destination = dest
	reqStruct.Data.Origin = org
	reqStruct.Data.NumOfAdult = 1
	reqStruct.Data.NumOfInfant = 0
	reqStruct.Data.ProviderType = "KAI"
	reqStruct.Data.Currency = "IDR"

	reqStruct.Data.DepartureDate = Date{
		Year:  departDateTime.Year(),
		Month: int(departDateTime.Month()),
		Day:   departDateTime.Day(),
	}

	payload, _ := json.Marshal(reqStruct)

	req, _ := http.NewRequest(
		"POST",
		"https://www.traveloka.com/api/v2/train/search/inventoryv2",
		bytes.NewBuffer(payload),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "id-ID,id;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://www.traveloka.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.5 Mobile/15E148 Safari/604.1")
	req.Header.Set("www-app-version", "release_webgtr_20251208-163cd17da7")
	req.Header.Set("x-client-interface", "mobile")
	req.Header.Set("x-did", "MDFLNjNNTjRYSkswTTVFR1lHS0NYNzFITlg=")
	req.Header.Set("x-domain", "train")
	req.Header.Set("x-route-prefix", "id-id")
	req.Header.Set("cookie", "tvl=PsMhBYzGArN88x++KFopFRFOX40j/KgUv28OBBa9VNp0m+0onfOIZNUfyjlCGoOLjdaQwhNIgdIoDpqKrEu0uaIsz6WboGxGDjrrDo9xXauxJk8RuysdEjH3Wg4Ipvws7Ly7zkxiC1uOdmpiTW0fAd6RQvJO5ebcghP0mfXrrYl0mjylVTf5gzXC/DZ5apq4PrWJzhEEK6G2GBWi776W6gDpvLoddwCDZ8OjvFDdrjEedR4NztmI8lIAtMmMJPwKoKXfAT2cfKxBv0+ewW4tWVxNQUY8S5vpxeF8xq4rJmnrYp1u7F5ezA6pqNwAUxAqv8r9/3O2XIQkDVNSuTEOPkFCBXSAEATIxRYSArR93+Lq4o2JelD5kVvcBrxgzxbmhkGY118inX83dd0tog3qv9gUUIPhFKTjTV21i1XfEYLW+rnm2i1g3sl7iRfZG/XfySj7Hu4y/iCrxfIGOnzBeSebP6zMAhrBkd/Yn0z1HtCIHZVE53lGMv+ahgzuurs3YrQdX9UcU5tJ4U+q+vVQgWjphu3aJKbZnDH4anP+197q6dO71ISJ5xxW3lmaIsoYWLY=~djAy; tvs=QNgs+aTLreIrHeD1GTgzoTuTEPVaSyaxJpgMA9v7Mvl4gljky88nP50AVIa0pZGPAFl9TxkyVI2V8Z9RDslqOl0Ehvf+4RXwWoHtpz0RjeVceKR7pPyt52WqV95SMakX8no6L6FBtD/a7hgA1G5NuEAK1xDY1pthh/dkxD5VRRZlG8wM3MetFq0lACzbDXH1jW8QVv2orPya4yiJMk6KGlk9slqliW8hWwRLWxW6TU20zZmw4VD4ha7tuBPWo4+SLTOEeJKN3URHP/ru94riKS9siHSkVVyKm1AUMuFvxv44+G3H/wpJvuS2WDlPdWsNgdnr6lqViOY78hgXynm3A1isQiqIK3Ihklr7jIVgm/Exn0aXMt9aeaahR5DintHErOr3+wGZ3wr+5G9JVKA/MDq4VDXPwCSvsQ==~djAy")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var trainResp = &TravelokaResp{}
	if err = json.Unmarshal(respBody, trainResp); err != nil {
		return
	}

	for _, t := range trainResp.Data.List {
		stock, _ := strconv.Atoi(t.Stock)
		trains = append(trains, Train{
			Name:  strings.ReplaceAll(strings.ReplaceAll(t.Name, "(", ""), ")", ""),
			Class: t.Segments[0].Summary.Class,
			Stock: stock,
		})
	}

	return
}

func KaiWebSchedule(org, dest, departDate string) (trains []Train, err error) {
	departDateTime, _ := time.Parse("02-01-2006", departDate)

	queryParam := fmt.Sprintf(
		"origination=%v&flexdatalist-origination=%v&destination=%v&flexdatalist-destination=%v&tanggal=%v&adult=1&infant=0", // &submit=Cari & Pesan Tiket
		org, org, dest, dest,
		fmt.Sprintf("%02d-%s-%d", departDateTime.Day(), monthID[departDateTime.Month()-1], departDateTime.Year()),
	)

	client := resty.New()
	// client.SetDebug(true)
	client.SetHeaders(map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Encoding": "utf-8",
	})

	// var redirectHistories []string
	client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		// redirectHistories = append(redirectHistories, req.URL.String())
		return nil
	}))

	resp, err := client.R().Get("https://booking.kai.id/?" + queryParam)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return
	}
	for _, node := range doc.Find(".list-kereta").Nodes {
		doc = goquery.NewDocumentFromNode(node)
		seatRaw := strings.Split(doc.Find(".sisa-kursi").Text(), " ")
		var stock int
		if len(seatRaw) > 1 {
			stock, err = strconv.Atoi(seatRaw[1])
			if err != nil {
				return
			}
		} else {
			if seatRaw[0] == "Habis" {
				stock = 0
			} else {
				stock = 99
			}
		}

		trainNumber := doc.Find(`[name="nokereta"]`).AttrOr("value", "")
		trainName := doc.Find(`[name="kereta"]`).AttrOr("value", "")
		subClass := doc.Find(`[name="subkelas"]`).AttrOr("value", "")

		trains = append(trains, Train{
			Name:  trainName + " " + trainNumber,
			Class: subClass,
			Stock: stock,
		})
	}

	return
}
