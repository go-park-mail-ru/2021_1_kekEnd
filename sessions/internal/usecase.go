package sessions

import (
	"github.com/satori/go.uuid"
)

func addPrefix(id uuid.UUID) string {
	return "sessions:" + id.String()
}
