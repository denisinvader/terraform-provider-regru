package regru

import "fmt"

func (c *Client) GetRecordA(zone, name string) (ZoneRecord, error) {
	res, err := c.GetZoneInfo(zone)
	if err != nil {
		return ZoneRecord{}, err
	}

	for _, record := range res.Records {
		if record.Type == RecordTypeA && record.Name == name {
			return record, nil
		}
	}

	return ZoneRecord{}, fmt.Errorf("Record %s not found in zone %s", name, zone)
}

func (c *Client) CreateRecordA(zone, name, ip string) error {
	data := map[string]string{
		"subdomain": name,
		"ipaddr":    ip,
	}

	_, err := c.fetch("zone/add_alias", data)

	if err != nil {
		c.deleteCachedZone(zone)
	}

	return err
}
