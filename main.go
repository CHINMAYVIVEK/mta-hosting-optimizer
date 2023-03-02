package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type IpConfig struct {
	IP       string `json:"IP"`
	Hostname string `json:"Hostname"`
	Active   bool   `json:"Active"`
}

func main() {
	http.HandleFunc("/hostnames", hostnamesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hostnamesHandler(w http.ResponseWriter, r *http.Request) {
	threshold := 1
	if x, ok := r.URL.Query()["x"]; ok {
		if val, err := strconv.Atoi(x[0]); err == nil {
			threshold = val
		}
	}
	ipConfigs, err := loadIpConfigs("ipconfigs.json")
	if err != nil {
		log.Fatal(err)
	}
	hostnameCounts := make(map[string]int)
	for _, ipConfig := range ipConfigs {
		if ipConfig.Active {
			hostnameCounts[ipConfig.Hostname]++
		}
	}
	var result []string
	for hostname, count := range hostnameCounts {
		if count <= threshold {
			result = append(result, hostname)
		}
	}
	jsonResult, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func loadIpConfigs(filename string) ([]IpConfig, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var ipConfigs []IpConfig
	err = json.Unmarshal(byteValue, &ipConfigs)
	if err != nil {
		return nil, err
	}

	return ipConfigs, nil
}
