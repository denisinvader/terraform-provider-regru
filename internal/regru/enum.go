package regru

const (
	ResultSuccess string = "success"
	ResultError   string = "error"
)

type Action string

const (
	ActionZoneInfo         Action = "zone/get_resource_records"
	ActionZoneRemoveRecord Action = "zone/remove_record"
	ActionZoneAddRecordA   Action = "zone/add_alias"
)

const (
	RecordTypeA     string = "A"
	RecordTypeAAAA  string = "AAAA"
	RecordTypeCAA   string = "CAA"
	RecordTypeCNAME string = "CNAME"
	RecordTypeMX    string = "MX"
	RecordTypeNS    string = "NS"
	RecordTypeSRV   string = "SRV"
	RecordTypeTXT   string = "TXT"
)
