package main

type SankeyData struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int    `json:"value"`
}

type SourceDropdownData struct {
	Affiliate string `json:"utm_affiliate"`
	Campaign  string `json:"utm_campaign"`
	Source    string `json:"utm_source"`
	Content   string `json:"utm_content"`
}

type NodelinkData struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type FunnelInputData struct {
	Type string `json:"type"`
	Node string `json:"node"`
}

type FunnelOutputData struct {
	Count    int64  `json:"count"`
	PageName string `json:"pagename"`
	PageType string `json:"pagetype"`
	Source   string `json:"utm_source"`
}

// Response schema
type Response struct {
	Name                string               `json:"name"`
	StatusCode          string               `json:"statuscode"`
	Status              string               `json:"status"`
	Message             string               `json:"message"`
	Sankey              []SankeyData         `json:"sankey"`
	SourceDropdown      []SourceDropdownData `json:"sourcedropdown"`
	DestinationDropdown []NodelinkData       `json:"destinationdropdown"`
	Nodelink            []NodelinkData       `json:"nodelink"`
	Funnel              []FunnelOutputData   `json:"funnel"`
}
