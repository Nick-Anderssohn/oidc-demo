package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserData struct {
	ID         string      `json:"id"`
	Email      string      `json:"email"`
	Identities []*Identity `json:"identities"`
}

type Identity struct {
	ID                 string          `json:"id"`
	IdentityProviderID string          `json:"identityProviderId"`
	ExternalID         string          `json:"externalId"`
	MostRecentIDToken  json.RawMessage `json:"mostRecentIdToken"`
}

type Service struct {
	Resolver *deps.Resolver
}

func (s *Service) GetUserData(ctx context.Context, userID pgtype.UUID) (*UserData, error) {
	userDataSlice, err := s.Resolver.Queries.GetUserData(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(userDataSlice) == 0 {
		return nil, fmt.Errorf("no user data found for user ID: %s", userID)
	}

	identities := []*Identity{}
	for _, userData := range userDataSlice {
		identities = append(identities, &Identity{
			ID:                 userData.IdentityID.String(),
			IdentityProviderID: userData.IdentityProviderID.String,
			ExternalID:         userData.ExternalID.String,
			MostRecentIDToken:  userData.MostRecentIDToken,
		})
	}

	return &UserData{
		ID:         userID.String(),
		Email:      userDataSlice[0].UserEmail,
		Identities: identities,
	}, nil
}
