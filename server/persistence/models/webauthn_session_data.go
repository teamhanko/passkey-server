package models

import (
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type Operation string

var (
	WebauthnOperationRegistration   Operation = "registration"
	WebauthnOperationAuthentication Operation = "authentication"
	WebauthnOperationTransaction    Operation = "transaction"
)

// WebauthnSessionData is used by pop to map your webauthn_session_data database table to your go code.
type WebauthnSessionData struct {
	ID                 uuid.UUID                              `db:"id"`
	UserId             string                                 `db:"user_id" json:"-"`
	Challenge          string                                 `db:"challenge"`
	UserVerification   string                                 `db:"user_verification"`
	CreatedAt          time.Time                              `db:"created_at"`
	UpdatedAt          time.Time                              `db:"updated_at"`
	Operation          Operation                              `db:"operation"`
	AllowedCredentials []WebauthnSessionDataAllowedCredential `has_many:"webauthn_session_data_allowed_credentials"`
	ExpiresAt          nulls.Time                             `db:"expires_at"`
	IsDiscoverable     bool                                   `db:"is_discoverable"`

	TenantID uuid.UUID `db:"tenant_id"`
	Tenant   *Tenant   `belongs_to:"tenants"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (sd *WebauthnSessionData) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "ID", Field: sd.ID},
		&validators.StringIsPresent{Name: "Challenge", Field: sd.Challenge},
		&validators.StringInclusion{Name: "Operation", Field: string(sd.Operation), List: []string{string(WebauthnOperationRegistration), string(WebauthnOperationAuthentication), string(WebauthnOperationTransaction)}},
		&validators.TimeIsPresent{Name: "UpdatedAt", Field: sd.UpdatedAt},
		&validators.TimeIsPresent{Name: "CreatedAt", Field: sd.CreatedAt},
	), nil
}
