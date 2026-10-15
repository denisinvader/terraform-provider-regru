package regru

import (
	"encoding/json"
	"fmt"
)

type getResourceRecordsAnswer struct {
	Domains []Zone `json:"domains,omitempty"`
}

func (c *Client) GetZoneInfo(name string) (Zone, error) {
	if c.cachedZones[name] != nil {
		return *c.cachedZones[name], nil
	}

	fmt.Printf("[debug] Requesting zone/get_resource_records for %s\n", name)

	res, err := c.fetch("zone/get_resource_records", map[string]string{"domain_name": name})
	if err != nil {
		return Zone{}, err
	}

	var response Response[getResourceRecordsAnswer]
	if err := json.Unmarshal(res, &response); err != nil {
		return Zone{}, err
	}

	if response.Result == ResultError {
		return Zone{}, RegruError{response.ErrorCode, response.ErrorText}
	}

	for _, zone := range response.Answer.Domains {
		if zone.Name == name {
			if zone.Result == ResultError {
				return Zone{}, RegruError{zone.ErrorCode, zone.ErrorText}
			}

			c.setCachedZone(name, zone)

			return zone, nil
		}
	}

	return Zone{}, fmt.Errorf("Zone %s not found in response", name)
}
