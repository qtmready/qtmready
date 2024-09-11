package auth

import (
	"context"
	"strings"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/quantm/internal/db"
)

type (
	// usrio represents the activities for the user.
	usrio struct{}
)

// UserIO creates and returns a new UserIO object.
//
// Example:
//
//	user_io := auth.UserIO()
func UserIO() *usrio {
	return &usrio{}
}

// Save saves a user to the database. The email address is converted to
// lowercase before saving.
//
// Note: Always provide a pointer to the complete User object to avoid
// creating a copy. The Save method will update the provided User object
// with any changes made by the database.
//
// Example:
//
//	user, err := auth.UserIO().Save(ctx, user)
func (a *usrio) Save(ctx context.Context, user *User) (*User, error) {
	user.Email = strings.ToLower(user.Email) // Lowercase email before saving

	return user, db.Save(user)
}

// Get retrieves a user from the database based on the provided parameters.
// If the "email" parameter is present, it is converted to lowercase before
// the query.
//
// Example:
//
//	user, err := auth.UserIO().Get(ctx, db.QueryParams{"id": user_id})
//	user, err := auth.UserIO().Get(ctx, db.QueryParams{"email": "user@example.com"})
func (a *usrio) Get(ctx context.Context, params db.QueryParams) (*User, error) {
	user := &User{}

	if email, ok := params["email"]; ok {
		params["email"] = strings.ToLower(email)
	}

	return user, db.Get(user, params)
}

// GetByID retrieves a user from the database by their ID.
//
// Example:
//
//	user, err := auth.UserIO().GetByID(ctx, user_id)
func (a *usrio) GetByID(ctx context.Context, id string) (*User, error) {
	user := &User{}

	return user, db.Get(user, db.QueryParams{"id": id})
}

// GetByEmail retrieves a user from the database by their email address.
// The email address is converted to lowercase before the query.
//
// Example:
//
//	user, err := auth.UserIO().GetByEmail(ctx, "user@example.com")
func (a *usrio) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}

	return user, db.Get(user, db.QueryParams{"email": strings.ToLower(email)})
}

// GetActiveTeam retrieves a team from the database associated with the given user.
//
// Example:
//
//	team, err := auth.UserIO().GetActiveTeam(ctx, user)
func (a *usrio) GetActiveTeam(ctx context.Context, user *User) (*Team, error) {
	team := &Team{}

	return team, db.Get(team, db.QueryParams{
		"id": user.TeamID.String(), // Convert TeamID to string
	})
}

// GetTeams retrieves all teams associated with a user.
//
// Example:
//
//	teams, err := auth.UserIO().GetTeams(ctx, user)
func (a *usrio) GetTeams(ctx context.Context, user *User) ([]Team, error) {
	entity := &Team{}
	teams := make([]Team, 0)
	ids := make([]string, 0)

	tus, err := TeamUserIO().GetByUserID(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	for _, tu := range tus {
		ids = append(ids, tu.TeamID.String())
	}

	query := db.
		SelectBuilder(entity.GetTable().Name()).
		Where(qb.In("id"))

	err = db.DB().
		Session.
		Query(query.ToCql()).BindMap(qb.M{"id": ids}).
		GetRelease(teams)

	return teams, err
}

// GetTeamUser retrieves a user from the database associated with the given team ID and user ID.
//
// Example:
//
//	teamUser, err := auth.UserIO().GetTeamUser(ctx, user_id, team_id)
func (a *usrio) GetTeamUser(ctx context.Context, user_id string, team_id string) (*TeamUser, error) {
	teamUser := &TeamUser{}

	return teamUser, db.Get(teamUser, db.QueryParams{
		"id":      user_id,
		"team_id": team_id,
	})
}

// SetPassword hashes the clear text password using bcrypt.
//
// Example:
//
//	user, err := auth.UserIO().SetPassword(ctx, user, "password")
func (a *usrio) SetPassword(ctx context.Context, user *User, password string) (*User, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(p)

	return user, nil
}

// VerifyPassword verifies the plain text password against the hashed password.
//
// Example:
//
//	isValid := auth.UserIO().VerifyPassword(ctx, user, "password")
func (a *usrio) VerifyPassword(ctx context.Context, user *User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}

	return nil
}

// SetActiveTeam sets the active team for the given user.
//
// Example:
//
//	user, err := auth.UserIO().SetActiveTeam(ctx, user, id)
func (a *usrio) SetActiveTeam(ctx context.Context, user *User, id string) (*User, error) {
	parsed, err := gocql.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	user.TeamID = parsed

	return user, nil
}
