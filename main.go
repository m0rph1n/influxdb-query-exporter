package main

import (
    "flag"
    "github.com/influxdata/influxdb1-client/v2"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "net/http"
)


func read_json(config string) map[string]interface{} {
    var result map[string]interface{}

    FileJson, err := ioutil.ReadFile(config)
    if err != nil {
        panic(err.Error())
    }
    json.Unmarshal([]byte(FileJson), &result)

    queries := result["queries"].(map[string]interface{})
    return queries
}


func get_query(host, database, query string) string {

    var val string

    c, err := client.NewHTTPClient(client.HTTPConfig{
            Addr: "http://" + host,
    })
    if err != nil {
        panic(err.Error())
    }
    q := client.NewQuery(query, database, "")
    response, err := c.Query(q)
    if err !=nil {
        panic(err.Error())
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

    var (
        host string
        config string
        listen string
    )

    flag.StringVar(&host, "host", "localhost:8086", "host to influxdb")
    flag.StringVar(&config, "config", "", "config with query to influxdb")
    flag.StringVar(&listen, "listen", "", "listen address")

    flag.Parse()
    queries := read_json(config)
    var responseText string
    message := func(w http.ResponseWriter, r *http.Request) {
        for database, _ := range queries {
            query_json := queries[database].(map[string]interface{})
	    for query := range query_json {
	        query_str := fmt.Sprintf("%v", query_json[query])
		responseText += fmt.Sprintf("influxdb_result_query{database=\"%s\", query_name=\"%s\"} %s\n", database, query, get_query(host, database, query_str))
	    }
        fmt.Fprintf(w, responseText)
        responseText = ""
        }
    }
    http.HandleFunc("/metrics", message)
    http.ListenAndServe(listen, nil)
}
