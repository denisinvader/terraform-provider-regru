package regru

import (
	"encoding/json"
	"fmt"
	"strings"
)

func getNormalizedCanonical(canonical, zone string) string {
	if !strings.Contains(canonical, ".") {
		return fmt.Sprintf("%s.%s.", canonical, zone)
	}

	if !strings.HasSuffix(canonical, ".") {
		return canonical + "."
	}

	return canonical
}

func getZoneAnswer(response Response, name string) (Zone, error) {
	var o Zone

	var answer ZonesAnswer
	if err := json.Unmarshal(response.Answer, &answer); err != nil {
		return o, err
	}

	for _, zone := range answer.Zones {
		if zone.Name == name {
			if zone.Result == ResultError {
				return o, RegruError{zone.ErrorCode, zone.ErrorText}
			}

			return zone, nil
		}
	}

	return o, fmt.Errorf("Zone %s not found in response.answer.domains", name)
}

func getZoneAnswerError(err error, response Response, zone string) error {
	if err == nil {
		_, err = getZoneAnswer(response, zone)
	}

	return err
}
