package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

var client *elastic.Client
var ctx context.Context

type metric struct {
	ClientID   string `json:"clientid"`
	UserAgent  string `json:"useragent"`
	XUserAgent string `json:"xuseragent"`
	Date       int64  `json:"date"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"metric":{
			"properties":{
				"clientid":{
					"type":"keyword"
				},
				"useragent":{
                    "type":"keyword"
                },
                "xuseragent":{
                    "type":"keyword"
                },
                "date": {
                    "type":   "date",
                    "format": "epoch_millis"
                }
			}
		}
	}
}`

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println("UserAgent", r.UserAgent())
	fmt.Println("ClientID", r.Header.Get("ClientID"))
	fmt.Println("XUserAgent", r.Header.Get("X-User-Agent"))
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	metric1 := metric{ClientID: r.Header.Get("ClientID"), UserAgent: r.UserAgent(), XUserAgent: r.Header.Get("X-User-Agent"), Date: time.Now().UnixNano() / 1000000}
	put1, err := client.Index().
		Index("metric").
		Type("metric").
		BodyJson(metric1).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed metric to index %s, type %s\n", put1.Index, put1.Type)

	fmt.Fprintf(w, "Hello astaxie!") // send data to client side
}

func main() {

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx = context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client1, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}
	client = client1

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("metric").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("metric").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	http.HandleFunc("/", sayhelloName)      // set router
	err = http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
