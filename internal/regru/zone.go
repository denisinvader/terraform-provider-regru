package regru

func (c *Client) GetZoneInfo(name string) (Zone, error) {
	var zone Zone

	res, err := c.fetch(ActionZoneInfo, map[string]string{"domain_name": name})

	if err == nil {
		zone, err = getZoneAnswer(res, name)
	}

	return zone, err
}

func (c *Client) UpdateZoneSoa(zone, ttl, minimumTtl string) error {
	res, err := c.fetch(ActionZoneUpdateSoa, map[string]string{
		"domain_name": zone,
		"ttl":         ttl,
		"minimum_ttl": minimumTtl,
	})

	return getZoneAnswerError(err, res, zone)
}
