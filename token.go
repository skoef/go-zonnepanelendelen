package zonnepanelendelen

import "fmt"

func (t AuthToken) String() string {
	return fmt.Sprintf("Token %s", t.Token)
}
