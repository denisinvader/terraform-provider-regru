package regru

import "strconv"

func (c *Client) GetRecordMX(zone, name, mailServer string, priority int) (ZoneRecord, error) {
	none := ZoneRecord{}
	content := getNormalizedCanonical(mailServer, zone)

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeMX && record.Name == name && record.Content == content && record.Priority == priority {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeMX, name, mailServer, &priority}
}

func (c *Client) CreateRecordMX(zone, name, mailServer string, priority int) error {
	res, err := c.fetch(ActionZoneAddRecordMX, map[string]string{
		"domain_name": zone,
		"subdomain":   name,
		"mail_server": getNormalizedCanonical(mailServer, zone),
		"priority":    strconv.Itoa(priority),
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordMX(zone, name, mailServer string, priority int) error {
	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeMX,
		"subdomain":   name,
		"content":     getNormalizedCanonical(mailServer, zone),
		"priority":    strconv.Itoa(priority),
	})

	return getZoneAnswerError(err, res, zone)
}
