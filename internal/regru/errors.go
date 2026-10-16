package regru

import "fmt"

type RegruError struct {
	Code string
	Text string
}

func (e RegruError) Error() string {
	return fmt.Sprintf("reg.ru error [%s]: %s", e.Code, e.Text)
}

type RegruRecordNotFoundError struct {
	Zone     string
	Type     string
	Name     string
	Content  string
	Priority *int
}

func (e RegruRecordNotFoundError) Error() string {
	if e.Priority != nil {
		return fmt.Sprintf(
			"reg.ru %s record %s -> %s (priority: %d) not found in zone %s",
			e.Type,
			e.Name,
			e.Content,
			*e.Priority,
			e.Zone,
		)
	}

	return fmt.Sprintf(
		"reg.ru %s record %s -> %s not found in zone %s",
		e.Type,
		e.Name,
		e.Content,
		e.Zone,
	)
}
