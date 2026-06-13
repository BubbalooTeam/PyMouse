package gsmarena

import (
	"fmt"
	"io"
	"pymouse/pymouse/helpers/rapidhttp"
	"pymouse/pymouse/helpers/telegram"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

// Device search by name structs
type GSMArenaDeviceSearchResult struct {
	ID          string
	Name        string
	Image       string
	Description string
}

type GSMArenaSearchResult struct {
	Results []GSMArenaDeviceSearchResult
}

// Device fetch by id structs
type Specification struct {
	Attrs string `json:"name"`
	Value string `json:"value"`
}

type PhoneDetail struct {
	Category       string          `json:"category"`
	Specifications []Specification `json:"specifications"`
}

type GSMArenaDeviceBaseResult struct {
	Name         string
	URL          string
	ImageURL     string
	PhoneDetails []PhoneDetail
}

type ParsedSpecs struct {
	Status       string
	Network      string
	Dimensions   string
	Weight       string
	Jack         string
	USB          string
	Sensors      string
	Battery      string
	Charging     string
	Display      string
	Chipset      string
	MainCamera   string
	SelfieCamera string
	Memory       string
}

var GSMArenaHeaders = map[string]string{
	"accept-language": "en-US,en;q=0.9",
	"cache-control":   "max-age=0",
	"priority":        "u=0, i",
	"user-agent":      "Dalvik/2.1.0 (Linux; U; Android 12; M2012K11AG Build/SQ1D.211205.017)",
	"referer":         "https://m.gsmarena.com",
}

const (
	GSMArenaBaseURL      = "https://m.gsmarena.com/%s"
	GSMArenaSearchExtURL = "results.php3?sQuickSearch=yes&sName=%s"
	GSMArenaDeviceExtURL = "%s.php"
	devicesPerPage       = 9
)

func getDataFromURL(urlExt string) (string, error) {
	http_client := rapidhttp.GetHTTPClient()
	r, err := rapidhttp.Request(
		http_client,
		rapidhttp.HTTPStruct{
			Method:  "GET",
			URL:     fmt.Sprintf(GSMArenaBaseURL, urlExt),
			Headers: GSMArenaHeaders,
		},
	)

	if err != nil {
		return "", fmt.Errorf("[BadRequest]: Failed to complete the request in GSMArena: %v", err)
	}
	if r.StatusCode != 200 {
		switch r.StatusCode {
		case 400:
			return "", fmt.Errorf("[BadRequest]: Failed to connect to the GSMArena website: %d.", r.StatusCode)
		case 404:
			return "", fmt.Errorf("[BadRequest]: The GSMArena URL is invalid, please review it: %d.", r.StatusCode)
		case 429:
			return "", fmt.Errorf("[TooManyRequests]: Please wait while we automatically unlock PyMouse access. This process may take minutes, hours, days, or even weeks: %d.", r.StatusCode)
		default:
			return "", fmt.Errorf("[GSMarena]: An unknown error occurred while requesting on the GSMarena website.")
		}

	}

	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("[DecodeError]: Error in decode request body: %v", err)
	}
	return string(bodyBytes), nil
}

func searchDevice(query string) *GSMArenaSearchResult {
	var results []GSMArenaDeviceSearchResult

	html, err := getDataFromURL(
		fmt.Sprintf(GSMArenaSearchExtURL, strings.ReplaceAll(query, " ", "+")),
	)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logrus.Error(err)
		return nil
	}

	doc.Find(".general-menu li").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a")
		img := s.Find("img")
		strong := s.Find("strong")

		id, _ := a.Attr("href")
		image, _ := img.Attr("src")
		description, _ := img.Attr("title")

		// Extract text from strong tag and clean up line breaks
		name := strings.TrimSpace(strong.Text())
		name = strings.ReplaceAll(name, "\n", " ")
		name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")

		results = append(results, GSMArenaDeviceSearchResult{
			ID:          strings.ReplaceAll(id, ".php", ""),
			Name:        name,
			Image:       image,
			Description: description,
		})
	})

	return &GSMArenaSearchResult{
		Results: results,
	}
}

func fetchDevice(deviceID string) *GSMArenaDeviceBaseResult {
	var phoneDetails []PhoneDetail

	html, err := getDataFromURL(fmt.Sprintf(GSMArenaDeviceExtURL, deviceID))
	if err != nil {
		logrus.Error(err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logrus.Error(err)
		return nil
	}

	name := strings.TrimSpace(
		doc.Find(".specs-phone-name-title").First().Text(),
	)

	image, _ := doc.Find(".specs-photo-main img").Attr("src")

	doc.Find("#specs-list table").Each(func(i int, table *goquery.Selection) {

		category := strings.TrimSpace(table.Find("th").First().Text())
		if category == "" {
			return
		}

		var specifications []Specification

		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			ttl := strings.TrimSpace(row.Find("td.ttl").Text())
			nfo := strings.TrimSpace(row.Find("td.nfo").Text())

			if ttl != "" && nfo != "" {
				specifications = append(specifications, Specification{
					Attrs: ttl,
					Value: nfo,
				})
			}
		})

		if len(specifications) > 0 {
			phoneDetails = append(phoneDetails, PhoneDetail{
				Category:       category,
				Specifications: specifications,
			})
		}
	})

	return &GSMArenaDeviceBaseResult{
		Name:         name,
		ImageURL:     image,
		URL:          fmt.Sprintf(GSMArenaBaseURL, fmt.Sprintf(GSMArenaDeviceExtURL, deviceID)),
		PhoneDetails: phoneDetails,
	}
}

func extractSpecifications(
	phone_details []PhoneDetail,
	categoryName string,
	attrs []string,
) string {
	for _, phone_detail := range phone_details {

		if phone_detail.Category != categoryName {
			continue
		}

		if len(attrs) == 0 {
			var values []string
			for _, spec := range phone_detail.Specifications {
				if spec.Value != "" {
					values = append(values, spec.Value)
				}
			}
			return strings.Join(values, "\n")
		}

		var values []string

		for _, spec := range phone_detail.Specifications {
			for _, attr := range attrs {
				if spec.Attrs == attr && spec.Value != "" {
					values = append(values, spec.Value)
				}
			}
		}

		return strings.Join(values, "\n")
	}

	return ""
}

func parseSpecifications(phone_details []PhoneDetail) ParsedSpecs {
	return ParsedSpecs{
		Status:       extractSpecifications(phone_details, "Launch", []string{"Status"}),
		Network:      extractSpecifications(phone_details, "Network", []string{"Technology"}),
		Dimensions:   extractSpecifications(phone_details, "Body", []string{"Dimensions"}),
		Weight:       extractSpecifications(phone_details, "Body", []string{"Weight"}),
		Jack:         extractSpecifications(phone_details, "Sound", []string{"3.5mm jack"}),
		USB:          extractSpecifications(phone_details, "Comms", []string{"USB"}),
		Sensors:      extractSpecifications(phone_details, "Features", []string{"Sensors"}),
		Battery:      extractSpecifications(phone_details, "Battery", []string{"Type"}),
		Charging:     extractSpecifications(phone_details, "Battery", []string{"Charging"}),
		Display:      extractSpecifications(phone_details, "Display", []string{"Type", "Size", "Resolution"}),
		Chipset:      extractSpecifications(phone_details, "Platform", []string{"Chipset", "CPU", "GPU"}),
		MainCamera:   extractSpecifications(phone_details, "Main Camera", nil),
		SelfieCamera: extractSpecifications(phone_details, "Selfie camera", nil),
		Memory:       extractSpecifications(phone_details, "Memory", []string{"Internal"}),
	}
}

func formatGSMarenaMessage(
	device *GSMArenaDeviceBaseResult,
	l func(string) string,
) string {

	parsed := parseSpecifications(device.PhoneDetails)

	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf(
		"<a href='%s'>\u2000</a>"+
			"<a href='%s'><b>%s</b></a>\n\n",
		device.ImageURL,
		device.URL,
		device.Name,
	))

	writeField := func(key, value string) {
		if strings.TrimSpace(value) == "" || value == "-" {
			return
		}

		b.WriteString(fmt.Sprintf(
			"<b>%s:</b> <i>%s</i>\n\n",
			l(fmt.Sprintf("gsmarena.phone-formatter.%s", key)),
			value,
		))
	}

	writeField("status", parsed.Status)
	writeField("network", parsed.Network)
	writeField("dimensions", parsed.Dimensions)
	writeField("weight", parsed.Weight)
	writeField("jack", parsed.Jack)
	writeField("usb", parsed.USB)
	writeField("sensors", parsed.Sensors)
	writeField("battery", parsed.Battery)
	writeField("charging", parsed.Charging)
	writeField("display", parsed.Display)
	writeField("chipset", parsed.Chipset)
	writeField("main-camera", parsed.MainCamera)
	writeField("selfie-camera", parsed.SelfieCamera)
	writeField("memory", parsed.Memory)

	return b.String()
}

func formatDeviceName(name string) string {
	reLowerNumber := regexp.MustCompile(`([a-z])([0-9])`)
	name = reLowerNumber.ReplaceAllString(name, `$1 $2`)

	reLowerUpper := regexp.MustCompile(`([a-z])([A-Z])`)
	name = reLowerUpper.ReplaceAllString(name, `$1 $2`)

	words := strings.Fields(name)
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
	}

	return strings.Join(words, " ")
}

func GSMarenaCreateKeyboard(
	search []GSMArenaDeviceSearchResult,
	userID int,
	searchID string,
	page int,
) *telego.InlineKeyboardMarkup {
	total := len(search)
	if total == 0 {
		return &telego.InlineKeyboardMarkup{}
	}

	totalPages := (total + devicesPerPage - 1) / devicesPerPage
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * devicesPerPage
	end := start + devicesPerPage
	if end > total {
		end = total
	}

	var rows [][]telego.InlineKeyboardButton
	var row []telego.InlineKeyboardButton

	for _, device := range search[start:end] {
		btn := telegoutil.InlineKeyboardButton(
			formatDeviceName(device.Name),
		).WithCallbackData(
			fmt.Sprintf("d|%s|%d|%s", device.ID, userID, searchID),
		)

		row = append(row, btn)
		if len(row) == 2 {
			rows = append(rows, row)
			row = []telego.InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		rows = append(rows, row)
	}
	if totalPages > 1 {
		callbackPattern := fmt.Sprintf("gsm_page|{number}|%d|%s", userID, searchID)

		pagination := telegram.KeyboardPaginate(
			totalPages,
			page,
			callbackPattern,
		)

		rows = append(rows, pagination...)
	}

	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
