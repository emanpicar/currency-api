package xmldata

import "encoding/xml"

type (
	Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		Gesmes  string   `xml:"gesmes,attr"`
		Xmlns   string   `xml:"xmlns,attr"`
		Subject string   `xml:"subject"`
		Sender  struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
		} `xml:"Sender"`
		Cube struct {
			Text string `xml:",chardata"`
			Cube []struct {
				Text string `xml:",chardata"`
				Time string `xml:"time,attr"`
				Cube []struct {
					Text     string  `xml:",chardata"`
					Currency string  `xml:"currency,attr"`
					Rate     float64 `xml:"rate,attr"`
				} `xml:"Cube"`
			} `xml:"Cube"`
		} `xml:"Cube"`
	}
)
