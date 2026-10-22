package regru

import (
	"fmt"
	"strconv"
)

func (c *Client) GetRecordHTTPS(zone, name, value string, priority int) (ZoneRecord, error) {
	none := ZoneRecord{}
	content := fmt.Sprintf(". %q", value)

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeHTTPS && record.Name == name && record.Content == content && record.Priority == priority {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeHTTPS, name, content, &priority}
}

func (c *Client) CreateRecordHTTPS(zone, name, value string, priority int) error {
	res, err := c.fetch(ActionZoneAddRecordHTTPS, map[string]string{
		"domain_name": zone,
		"subdomain":   name,
		"target":      ".",
		"value":       value,
		"priority":    strconv.Itoa(priority),
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordHTTPS(zone, name, value string, priority int) error {
	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeHTTPS,
		"subdomain":   name,
		"content":     value,
		"priority":    strconv.Itoa(priority),
	})

	return getZoneAnswerError(err, res, zone)
}
