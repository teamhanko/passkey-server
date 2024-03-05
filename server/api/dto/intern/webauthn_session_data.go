package intern

import (
	"encoding/base64"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"github.com/teamhanko/passkey-server/persistence/models"
	"strings"
	"time"
)

func WebauthnSessionDataFromModel(data *models.WebauthnSessionData) *webauthn.SessionData {
	var allowedCredentials [][]byte
	for _, credential := range data.AllowedCredentials {
		credentialId, err := base64.RawURLEncoding.DecodeString(credential.CredentialId)
		if err != nil {
			continue
		}
		allowedCredentials = append(allowedCredentials, credentialId)
	}
	var userId []byte = nil
	if strings.TrimSpace(data.UserId) != "" {
		userId = []byte(data.UserId)
	}
	return &webauthn.SessionData{
		Challenge:            data.Challenge,
		UserID:               userId,
		AllowedCredentialIDs: allowedCredentials,
		UserVerification:     protocol.UserVerificationRequirement(data.UserVerification),
		Expires:              data.ExpiresAt.Time,
	}
}

func WebauthnSessionDataToModel(data *webauthn.SessionData, tenantId uuid.UUID, operation models.Operation, isDiscoverable bool) *models.WebauthnSessionData {
	id, _ := uuid.NewV4()
	now := time.Now()

	var allowedCredentials []models.WebauthnSessionDataAllowedCredential
	for _, credentialID := range data.AllowedCredentialIDs {
		aId, _ := uuid.NewV4()
		allowedCredential := models.WebauthnSessionDataAllowedCredential{
			ID:                    aId,
			CredentialId:          base64.RawURLEncoding.EncodeToString(credentialID),
			WebauthnSessionDataID: id,
			CreatedAt:             now,
			UpdatedAt:             now,
		}

		allowedCredentials = append(allowedCredentials, allowedCredential)
	}

	return &models.WebauthnSessionData{
		ID:                 id,
		Challenge:          data.Challenge,
		UserId:             string(data.UserID),
		UserVerification:   string(data.UserVerification),
		CreatedAt:          now,
		UpdatedAt:          now,
		Operation:          operation,
		AllowedCredentials: allowedCredentials,
		ExpiresAt:          nulls.NewTime(data.Expires),
		TenantID:           tenantId,
		IsDiscoverable:     isDiscoverable,
	}
}
