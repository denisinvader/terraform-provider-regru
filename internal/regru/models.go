package regru

type Result = string

const (
	ResultSuccess Result = "success"
	ResultError   Result = "error"
)

type RecordType = string

const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCAA   RecordType = "CAA"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypeNS    RecordType = "NS"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypeTXT   RecordType = "TXT"
)

type Zone struct {
	Name        string       `json:"dname"`
	Result      Result       `json:"result"`
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
	Content  string     `json:"content"`
	Priority int        `json:"prio"`
	Type     RecordType `json:"rectype"`
	Name     string     `json:"subname"`
}

type Response[T any] struct {
	Result      Result
	Answer      T              `json:"answer,omitempty"`
	InputParams map[string]any `json:"input_params,omitempty"`
	ErrorCode   string         `json:"error_code,omitempty"`
	ErrorText   string         `json:"error_text,omitempty"`
	ErrorParams map[string]any `json:"error_params,omitempty"`
}
