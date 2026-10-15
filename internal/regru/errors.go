package regru

import (
	"fmt"
)

type RegruError struct {
	Code string
	Text string
}

func (e RegruError) Error() string {
	return fmt.Sprintf("reg.ru error [%s]: %s", e.Code, e.Text)
}
