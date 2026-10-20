package regru

func (c *Client) GetRecordA(zone, name, ip string) (ZoneRecord, error) {
	none := ZoneRecord{}

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeA && record.Name == name && record.Content == ip {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeA, name, ip, nil}
}

func (c *Client) CreateRecordA(zone, name, ip string) error {
	res, err := c.fetch(ActionZoneAddRecordA, map[string]string{
		"domain_name": zone,
		"subdomain":   name,
		"ipaddr":      ip,
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordA(zone, name, ip string) error {
	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeA,
		"subdomain":   name,
		"content":     ip,
	})

	return getZoneAnswerError(err, res, zone)
}
