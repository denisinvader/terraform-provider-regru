package regru

import "strconv"

func (c *Client) GetRecordNS(zone, name, canonical string, index int) (ZoneRecord, error) {
	none := ZoneRecord{}
	content := getNormalizedCanonical(canonical, zone)

	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return none, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeNS && record.Name == name && record.Content == content && record.Priority == index {
			return record, nil
		}
	}

	return none, RegruRecordNotFoundError{zone, RecordTypeNS, name, canonical, &index}
}

func (c *Client) CreateRecordNS(zone, name, canonical string, index int) error {
	res, err := c.fetch(ActionZoneAddRecordNS, map[string]string{
		"domain_name":   zone,
		"subdomain":     name,
		"dns_server":    getNormalizedCanonical(canonical, zone),
		"record_number": strconv.Itoa(index),
	})

	return getZoneAnswerError(err, res, zone)
}

func (c *Client) DeleteRecordNS(zone, name, canonical string, index int) error {
	res, err := c.fetch(ActionZoneRemoveRecord, map[string]string{
		"domain_name": zone,
		"record_type": RecordTypeNS,
		"subdomain":   name,
		"content":     getNormalizedCanonical(canonical, zone),
		"priority":    strconv.Itoa(index),
	})

	return getZoneAnswerError(err, res, zone)
}
