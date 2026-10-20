package regru

import "encoding/json"

type Zone struct {
	Name        string       `json:"dname"`
	Result      string       `json:"result"`
	Records     []ZoneRecord `json:"rrs,omitempty"`
	ServiceID   string       `json:"service_id,omitempty"`
	ServiceType string       `json:"servtype,omitempty"`
	Soa         ZoneSoa      `json:"soa"`
	ErrorCode   string       `json:"error_code,omitempty"`
	ErrorText   string       `json:"error_text,omitempty"`
}

type ZoneSoa struct {
	MinimumTtl string `json:"minimum_ttl"`
	Ttl        string `json:"ttl"`
}

type ZoneRecord struct {
	Content  string `json:"content"`
	Priority int    `json:"prio"`
	Type     string `json:"rectype"`
	Name     string `json:"subname"`
}

type Response struct {
	Result      string
	Answer      json.RawMessage `json:"answer,omitempty"`
	InputParams map[string]any  `json:"input_params,omitempty"`
	ErrorCode   string          `json:"error_code,omitempty"`
	ErrorText   string          `json:"error_text,omitempty"`
	ErrorParams map[string]any  `json:"error_params,omitempty"`
}

type ZonesAnswer struct {
	Zones []Zone `json:"domains,omitempty"`
}
