package main

import (
    "flag"
    "github.com/influxdata/influxdb1-client/v2"
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
    "net/http"
)


func get_query(host, database, query string) string {

    var val string

    c, err := client.NewHTTPClient(client.HTTPConfig{
            Addr: "http://" + host,
    })
    if err != nil {
        return err.Error()
    }
    q := client.NewQuery(query, database, "")
    response, err := c.Query(q)
    if err !=nil {
        return err.Error()
    }
    if len(response.Results[0].Series) != 0 {
        values := response.Results[0].Series[0].Values
        for i := 0; i < len(values); i++ {
            value := values[i]
	    val = fmt.Sprint(value[1])
        }
    }
    return val
}


func main() {
    var host string
    var database string
    var config string
    var listen string
    var result map[string]map[string]string


    flag.StringVar(&host, "host", "localhost:8086", "host to influxdb")
    flag.StringVar(&database, "database", "", "database to influxdb")
    flag.StringVar(&config, "config", "", "config with query to influxdb")
    flag.StringVar(&listen, "listen", "", "listen address")

    flag.Parse()

    jsonFile, err := os.Open(config)
    defer jsonFile.Close()

    if err != nil {
        return
    }
    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal([]byte(byteValue), &result)
    http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        for key := range result["queries"] {
            query := result["queries"][key]
            get_query(host, database, query)
	    fmt.Fprintf(w, key + "\t" + get_query(host, database, query) + "\n")
        }
    })
    http.ListenAndServe(listen, nil)
}
