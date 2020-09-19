package stores

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func key(IDs ...uuid.UUID) string {
	key := ""
	for _, ID := range IDs {
		if key != "" {
			key += ":"
		}
		key += ID.String()
	}
	return key
}

func scopedKey(scope string, IDs ...uuid.UUID) string {
	return fmt.Sprintf("%s/%s", scope, key(IDs...))
}
