package Utils

import (
	"github.com/google/uuid"
	"time"
)

type Meta struct {
	PK           string `json:"pk,omitempty"`
	SK           string `json:"sk,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	LastModified int64  `json:"last_modified,omitempty"`
}

func (s *Meta) SetLastModifiedNow() {
	now := time.Now().Unix()
	s.LastModified = now
}
func (s *Meta) SetCreatedAtNow() {
	now := time.Now().Unix()
	s.LastModified = now
}
func (s *Meta) GenerateNewId(parentId string, prefix string) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if parentId == "" {
		s.SK = prefix + id.String()
	} else {
		s.SK = parentId + "_" + prefix + id.String()
	}
	s.PK = prefix

	return nil
}

func (s *Meta) New(parentId string, prefix string) error {
	err := s.GenerateNewId(parentId, prefix)
	if err != nil {
		return err
	}
	s.SetCreatedAtNow()
	s.SetLastModifiedNow()

	return nil
}
