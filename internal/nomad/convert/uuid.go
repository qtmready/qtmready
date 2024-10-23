package convert

import (
	"github.com/google/uuid"

	commonv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/common/v1"
)

func ProtoToUUID(id *commonv1.UUID) uuid.UUID {
	return uuid.MustParse(id.GetValue())
}

func UUIDToProto(id uuid.UUID) *commonv1.UUID {
	return &commonv1.UUID{Value: id.String()}
}
