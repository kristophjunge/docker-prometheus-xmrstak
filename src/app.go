package main

import (
    "io"
    "net/http"
    "log"
    "os"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "errors"
)

const LISTEN_ADDRESS = ":9209"

var apiUrl string
var minerId string
var testMode string

type XMRStakStatistics struct {
    Hashrate struct {
        Total []float64 `json:"total"`
    } `json:"hashrate"`
}

func integerToString(value int64) string {
    return strconv.FormatInt(value, 10)
}

func floatToString(value float64, precision int64) string {
    return strconv.FormatFloat(value, 'f', int(precision), 64)
}

func stringToInteger(value string) int64 {
    if value == "" {
        return 0
    }
    result, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
        log.Fatal(err)
    }
    return result
}

func stringToFloat(value string) float64 {
    if value == "" {
        return 0
    }
    result, err := strconv.ParseFloat(value, 64)
    if err != nil {
        log.Fatal(err)
    }
    return result
}

func formatValue(key string, meta string, value string) string {
    result := key;
    if (meta != "") {
        result += "{" + meta + "}";
    }
    result += " "
    result += value
    result += "\n"
    return result
}

func queryData() (string, error) {
    var err error

    // Perform HTTP request
    resp, err := http.Get(apiUrl);
    if err != nil {
        return "", err;
    }

    // Parse response
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return "", errors.New("HTTP returned code " + integerToString(int64(resp.StatusCode)))
    }
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    if err != nil {
        return "", err;
    }

    return bodyString, nil;
}

func getTestData() (string, error) {
    dir, err := os.Getwd()
    if err != nil {
        return "", err;
    }
    body, err := ioutil.ReadFile(dir + "/test.json")
    if err != nil {
        return "", err;
    }
    return string(body), nil
}

func metrics(w http.ResponseWriter, r *http.Request) {
    log.Print("Serving /metrics")

    var up int64 = 1
    var jsonString string
    var err error

    if (testMode == "1") {
        jsonString, err = getTestData()
    } else {
        jsonString, err = queryData()
    }
    if err != nil {
        log.Print(err)
        up = 0
    }

    // Parse JSON
    jsonData := XMRStakStatistics{}
    json.Unmarshal([]byte(jsonString), &jsonData)

    // Output
    io.WriteString(w, formatValue("xmrstak_up", "miner=\"" + minerId + "\"", integerToString(up)))
    io.WriteString(w, formatValue("xmrstak_hashrate", "miner=\"" + minerId + "\"", floatToString(jsonData.Hashrate.Total[1], 1)))
}

func index(w http.ResponseWriter, r *http.Request) {
    log.Print("Serving /index")
    html := `<!doctype html>
<html>
    <head>
        <meta charset="utf-8">
        <title>XMR-Stark Exporter</title>
    </head>
    <body>
        <h1>XMR-Stark Exporter</h1>
        <p><a href="/metrics">Metrics</a></p>
    </body>
</html>`
    io.WriteString(w, html)
}

func main() {
    testMode = os.Getenv("TEST_MODE")
    if (testMode == "1") {
        log.Print("Test mode is enabled")
    }

    apiUrl = os.Getenv("API_URL")
    log.Print("API URL: " + apiUrl)

    minerId = os.Getenv("MINER_ID")
    log.Print("Miner ID: " + minerId)

    log.Print("XMR-Stark exporter listening on " + LISTEN_ADDRESS)
    http.HandleFunc("/", index)
    http.HandleFunc("/metrics", metrics)
    http.ListenAndServe(LISTEN_ADDRESS, nil)
}
