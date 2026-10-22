package regru

func (c *Client) GetRecordTXT(zone, name, text string) (ZoneRecord, error) {
	none := ZoneRecord{}

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeTXT && record.Name == name && record.Content == text {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeTXT, name, text, nil}
}

func (c *Client) CreateRecordTXT(zone, name, text string) error {
	res, err := c.fetch(ActionZoneAddRecordTXT, map[string]string{
		"domain_name": zone,
		"subdomain":   name,
		"text":        text,
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordTXT(zone, name, text string) error {
	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeTXT,
		"subdomain":   name,
		"content":     text,
	})

	return getZoneAnswerError(err, res, zone)
}
