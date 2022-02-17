package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var url string

type HostInfo struct {
	Host string `json:"ip"`
	Port string `json:"port"`
}

func GetIP() string {
	var hostInfo HostInfo
	data, err := ioutil.ReadFile("host.json")
	if err != nil {
		fmt.Println("host json file read error")
		panic(err)
	}
	err = json.Unmarshal(data, &hostInfo)
	if err != nil {
		fmt.Println("host json file unmarshal error")
		panic(err)
	}
	return hostInfo.Host
}

func GetUrl() string {
	if url != "" {
		return url
	} else {
		var hostInfo HostInfo
		data, err := ioutil.ReadFile("host.json")
		if err != nil {
			fmt.Println("host json file read error")
			panic(err)
		}
		err = json.Unmarshal(data, &hostInfo)
		if err != nil {
			fmt.Println("host json file unmarshal error")
			panic(err)
		}
		url = "http://" + hostInfo.Host + ":" + hostInfo.Port
		return url
	}
}
