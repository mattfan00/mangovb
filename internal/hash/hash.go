package hash

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	vb "github.com/mattfan00/nycvbtracker"
)

func Hash(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

type uniqueEvent struct {
	Source    string    `json:"source"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	StartDate time.Time `json:"startDate"`
}

func HashEvent(event vb.Event) (string, error) {
	b, err := json.Marshal(uniqueEvent{
		Source:    event.Source,
		Name:      event.Name,
		Location:  event.Location,
		StartDate: event.StartDate,
	})
	if err != nil {
		return "", err
	}

	return Hash(string(b)), nil
}
