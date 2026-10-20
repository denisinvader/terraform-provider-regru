package regru

func (c *Client) GetRecordCNAME(zone, name, canonical string) (ZoneRecord, error) {
	none := ZoneRecord{}
	content := getNormalizedCanonical(canonical, zone)

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeCNAME && record.Name == name && record.Content == content {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeCNAME, name, canonical, nil}
}

func (c *Client) CreateRecordCNAME(zone, name, canonical string) error {
	res, err := c.fetch(ActionZoneAddRecordCNAME, map[string]string{
		"domain_name":    zone,
		"subdomain":      name,
		"canonical_name": getNormalizedCanonical(canonical, zone),
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordCNAME(zone, name, canonical string) error {
	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeCNAME,
		"subdomain":   name,
		"content":     getNormalizedCanonical(canonical, zone),
	})

	return getZoneAnswerError(err, res, zone)
}
