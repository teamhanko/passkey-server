package intern

import (
	"encoding/base64"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofrs/uuid"
	"github.com/teamhanko/passkey-server/mapper"
	"github.com/teamhanko/passkey-server/persistence/models"
	"time"
)

func WebauthnCredentialToModel(credential *webauthn.Credential, userId string, webauthnUserId uuid.UUID, backupEligible bool, backupState bool, authenticatorMetadata mapper.AuthenticatorMetadata, isMFACredential bool) *models.WebauthnCredential {
	now := time.Now().UTC()
	aaguid, _ := uuid.FromBytes(credential.Authenticator.AAGUID)
	credentialID := base64.RawURLEncoding.EncodeToString(credential.ID)
	name := authenticatorMetadata.GetNameForAaguid(aaguid)
	if name == nil {
		genericName := fmt.Sprintf("cred-%s", credentialID)
		name = &genericName
	}

	c := &models.WebauthnCredential{
		ID:              credentialID,
		Name:            name,
		UserId:          userId,
		PublicKey:       base64.RawURLEncoding.EncodeToString(credential.PublicKey),
		AttestationType: credential.AttestationType,
		AAGUID:          aaguid,
		SignCount:       int(credential.Authenticator.SignCount),
		LastUsedAt:      &now,
		CreatedAt:       now,
		UpdatedAt:       now,
		BackupEligible:  backupEligible,
		BackupState:     backupState,
		IsMFA:           isMFACredential,

		WebauthnUserID: webauthnUserId,
	}

	for _, name := range credential.Transport {
		if string(name) != "" {
			id, _ := uuid.NewV4()
			t := models.WebauthnCredentialTransport{
				ID:                   id,
				Name:                 string(name),
				WebauthnCredentialID: credentialID,
			}
			c.Transports = append(c.Transports, t)
		}
	}

	return c
}

func WebauthnCredentialFromModel(credential *models.WebauthnCredential) *webauthn.Credential {
	credId, _ := base64.RawURLEncoding.DecodeString(credential.ID)
	pKey, _ := base64.RawURLEncoding.DecodeString(credential.PublicKey)
	transport := make([]protocol.AuthenticatorTransport, len(credential.Transports))

	for i, t := range credential.Transports {
		transport[i] = protocol.AuthenticatorTransport(t.Name)
	}

	return &webauthn.Credential{
		ID:              credId,
		PublicKey:       pKey,
		AttestationType: credential.AttestationType,
		Authenticator: webauthn.Authenticator{
			AAGUID:    credential.AAGUID.Bytes(),
			SignCount: uint32(credential.SignCount),
		},
		Transport: transport,
	}
}
