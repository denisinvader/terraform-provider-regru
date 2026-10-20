package regru

func (c *Client) GetRecordAAAA(zone, name, ip string) (ZoneRecord, error) {
	none := ZoneRecord{}

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeAAAA && record.Name == name && record.Content == ip {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeAAAA, name, ip, nil}
}

func (c *Client) CreateRecordAAAA(zone, name, ip string) error {
	res, err := c.fetch(ActionZoneAddRecordAAAA, map[string]string{
		"domain_name": zone,
		"subdomain":   name,
		"ipaddr":      ip,
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordAAAA(zone, name, ip string) error {
	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeAAAA,
		"subdomain":   name,
		"content":     ip,
	})

	return getZoneAnswerError(err, res, zone)
}
