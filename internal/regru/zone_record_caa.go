package regru

import "fmt"

func (c *Client) GetRecordCAA(zone, name, flag, tag, value string) (ZoneRecord, error) {
	none := ZoneRecord{}

	content := fmt.Sprintf("%s %s %q", flag, tag, value)

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeCAA && record.Name == name && record.Content == content {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeCAA, name, content, nil}
}

func (c *Client) CreateRecordCAA(zone, name, flag, tag, value string) error {
	res, err := c.fetch(ActionZoneAddRecordCAA, map[string]string{
		"domain_name": zone,
		"subdomain":   name,
		"flags":       flag,
		"tag":         tag,
		"value":       value,
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordCAA(zone, name, flag, tag, value string) error {
	content := fmt.Sprintf("%s %s %q", flag, tag, value)

	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeCAA,
		"subdomain":   name,
		"content":     content,
	})

	return getZoneAnswerError(err, res, zone)
}
