package nomad

import (
	"context"

	"github.com/google/uuid"
)

type (
	AuthContext string
)

const (
	AuthContextUser AuthContext = "user_id"
	AuthContextOrg  AuthContext = "org_id"
)

func GetAuthContext(ctx context.Context) (uuid.UUID, uuid.UUID) {
	user_id := ctx.Value(AuthContextUser)
	org_id := ctx.Value(AuthContextOrg)

	return uuid.MustParse(user_id.(string)), uuid.MustParse(org_id.(string))
}
