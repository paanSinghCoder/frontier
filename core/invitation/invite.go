package invitation

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/raystack/frontier/pkg/metadata"
)

var (
	ErrNotFound = errors.New("invitation not found")
)

const (
	DefaultExpiryDuration = 24 * time.Hour * 7
)

type Invitation struct {
	ID        uuid.UUID
	UserID    string
	OrgID     string
	GroupIDs  []string
	RoleIDs   []string
	Metadata  metadata.Metadata
	CreatedAt time.Time
	ExpiresAt time.Time
}
