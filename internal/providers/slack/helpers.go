package slack

import (
	"fmt"
	"time"
)

// returns a time.Time and error.
func getDateFromSubmittedValue(submittedValue string) (time.Time, error) {
	return time.Parse("2006-01-02", submittedValue)
}

// GenerateUniqueExternalID generates a unique external ID based on user ID and timestamp.
func GenerateUniqueExternalID(userID string) string {
	return fmt.Sprintf("user-%s-%d", userID, time.Now().UnixNano())
}
