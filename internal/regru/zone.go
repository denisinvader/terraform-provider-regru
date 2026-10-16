package regru

import (
	"encoding/json"
	"fmt"
)

type getZoneInfoAnswer struct {
	Domains []Zone `json:"domains,omitempty"`
}

func (c *Client) GetZoneInfo(name string) (Zone, error) {
	res, err := c.fetch(ActionZoneInfo, map[string]string{"domain_name": name})
	if err != nil {
		return Zone{}, err
	}

	var answer getZoneInfoAnswer
	if err := json.Unmarshal(res.Answer, &answer); err != nil {
		return Zone{}, err
	}

	for _, zone := range answer.Domains {
		if zone.Name == name {
			if zone.Result == ResultError {
				return Zone{}, RegruError{zone.ErrorCode, zone.ErrorText}
			}

			return zone, nil
		}
	}

	return Zone{}, fmt.Errorf("zone %s not found in response", name)
}

func (c *Client) UpdateZoneSoa(name, ttl, minimumTtl string) error {
	_, err := c.fetch(ActionZoneUpdateSoa, map[string]string{
		"domain_name": name,
		"ttl":         ttl,
		"minimum_ttl": minimumTtl,
	})

	return err
}
