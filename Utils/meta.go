package Utils

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type Meta struct {
	PK           string `json:"pk,omitempty"`
	SK           string `json:"sk,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	LastModified int64  `json:"last_modified,omitempty"`
	IsDeleted    bool   `json:"is_deleted,omitempty"`
}

func (s *Meta) SetLastModifiedNow() {
	now := time.Now().Unix()
	s.LastModified = now
}
func (s *Meta) SetCreatedAtNow() {
	now := time.Now().Unix()
	s.CreatedAt = now
}
func (s *Meta) GenerateNewId(prefix string, parents ...string) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if len(parents) == 0 {
		s.SK = prefix + id.String()
	} else {
		if len(parents) == 1 && parents[0] == "" {
			s.SK = prefix + id.String()
		} else {
			s.SK = strings.Join(parents, "_") + "_" + prefix + id.String()
		}
	}
	s.PK = prefix

	return nil
}

// New generates a new PK and SK.
//
// Sets Created at and last modified to now.
func (s *Meta) New(prefix string, parents ...string) error {
	err := s.GenerateNewId(prefix, parents...)
	if err != nil {
		return err
	}
	s.SetCreatedAtNow()
	s.SetLastModifiedNow()

	return nil
}
