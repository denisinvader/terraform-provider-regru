package regru

import (
	"fmt"
	"strconv"
)

func (c *Client) GetRecordSRV(zone, service, target string, port, weight, priority int) (ZoneRecord, error) {
	none := ZoneRecord{}
	content := fmt.Sprintf("%d %d %s", weight, port, getNormalizedCanonical(target, zone))

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeSRV && record.Name == service && record.Content == content && record.Priority == priority {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeSRV, service, target, nil}
}

func (c *Client) CreateRecordSRV(zone, service, target string, port, weight, priority int) error {
	res, err := c.fetch(ActionZoneAddRecordSRV, map[string]string{
		"domain_name": zone,
		"service":     service,
		"target":      getNormalizedCanonical(target, zone),
		"port":        strconv.Itoa(port),
		"weight":      strconv.Itoa(weight),
		"priority":    strconv.Itoa(priority),
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordSRV(zone, service, target string, port, weight, priority int) error {
	content := fmt.Sprintf("%d %d %s", weight, port, getNormalizedCanonical(target, zone))

	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeSRV,
		"subdomain":   service,
		"content":     content,
		"priority":    strconv.Itoa(priority),
	})

	return getZoneAnswerError(err, res, zone)
}
