package devices

import (
	"encoding/json"
	"io"
	"pymouse/pymouse/helpers/rapidhttp"
	"strings"

	"github.com/sirupsen/logrus"
)

type WhatISBaseResult struct {
	Name   string
	Brand  string
	Device string
	Model  string
}

type deviceEntry struct {
	Name   string `json:"name"`
	Brand  string `json:"brand"`
	Device string `json:"device"`
	Model  string `json:"model"`
}

func getDevice(codename string) (WhatISBaseResult, bool) {
	var url string
	var key string

	httpClient := rapidhttp.GetHTTPClient()

	if strings.HasPrefix(strings.ToLower(codename), "sm-") {
		url = "https://raw.githubusercontent.com/androidtrackers/certified-android-devices/master/by_model.json"
		key = strings.ToUpper(codename)
	} else {
		url = "https://raw.githubusercontent.com/androidtrackers/certified-android-devices/master/by_device.json"

		newDevice := strings.ToLower(codename)
		if strings.HasPrefix(codename, "beyond") {
			newDevice = strings.TrimSuffix(newDevice, "lte")
		}
		key = newDevice
	}

	resp, err := rapidhttp.Request(
		httpClient,
		rapidhttp.HTTPStruct{
			Method: "GET",
			URL:    url,
		},
	)
	if err != nil {
		logrus.Error("failed to fetch device information:", err)
		return WhatISBaseResult{}, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("failed to read response body:", err)
		return WhatISBaseResult{}, false
	}

	var db map[string][]deviceEntry
	if err := json.Unmarshal(body, &db); err != nil {
		logrus.Error("failed to parse json:", err)
		return WhatISBaseResult{}, false
	}

	entries, ok := db[key]
	if !ok || len(entries) == 0 {
		return WhatISBaseResult{}, false
	}

	entry := entries[0]

	result := WhatISBaseResult{
		Name:   entry.Name,
		Brand:  entry.Brand,
		Model:  entry.Model,
		Device: entry.Device,
	}

	return result, true
}
