package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func TrainsSchedule(org, dest, departDate string) (trains []Train, err error) {
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

		train.Name = strings.ToUpper(strings.Replace(goquery.NewDocumentFromNode(node).Find(
			"div.text-center.text-lg-start > strong",
		).Text(), "- ", "", 1))

		train.Class = strings.ToUpper(strings.Trim(goquery.NewDocumentFromNode(node).Find(
			"div.col-12.col-6.col-sm-2.text-center.text-sm-end.text-dark.d-flex.flex-column.text-6.price > div > span",
		).Text(), "( )"))

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
