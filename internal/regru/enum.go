package regru

const (
	ResultSuccess string = "success"
	ResultError   string = "error"
)

type Action string

const (
	ActionZoneInfo           Action = "zone/get_resource_records"
	ActionZoneUpdateSoa      Action = "zone/update_soa"
	ActionZoneRemoveRecord   Action = "zone/remove_record"
	ActionZoneAddRecordA     Action = "zone/add_alias"
	ActionZoneAddRecordAAAA  Action = "zone/add_aaaa"
	ActionZoneAddRecordCNAME Action = "zone/add_cname"
	ActionZoneAddRecordMX    Action = "zone/add_mx"
	ActionZoneAddRecordNS    Action = "zone/add_ns"
	ActionZoneAddRecordTXT   Action = "zone/add_txt"
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
