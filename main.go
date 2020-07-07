package main

import (
    "flag"
    "github.com/influxdata/influxdb1-client/v2"
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
//    "net/http"
)


func get_query(host, database, query string) {

    c, err := client.NewHTTPClient(client.HTTPConfig{
	    Addr: "http://" + host,
    })
    if err != nil {
	fmt.Println(err.Error())
    }

    q := client.NewQuery(query, database, "")
    if response, err := c.Query(q); err == nil && response.Error() == nil {
        values := response.Results[0].Series[0].Values
        for i := 0; i < len(values); i++ {
            value := values[i]

            fmt.Println(value[1])
	}
    }
}


func main() {
    var host string
    var database string
    var config string

    flag.StringVar(&host, "host", "localhost:8086", "host to influxdb")
    flag.StringVar(&database, "database", "", "database to influxdb")
    flag.StringVar(&config, "config", "", "config with query to influxdb")

    flag.Parse()

    jsonFile, err := os.Open("config.json")

    if err != nil {
        fmt.Println(err)
    }
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var result map[string]map[string]string
    json.Unmarshal([]byte(byteValue), &result)
    for key := range result["queries"] {
        query := result["queries"][key]
         get_query(host, database, query)
    }
//    http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
//        fmt.Fprintf(w, test)
//    })
//    http.ListenAndServe(":8003", nil)


}
