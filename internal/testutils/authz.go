package testutils

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/openkcm/cmk/internal/constants"
	"github.com/openkcm/cmk/internal/model"
	"github.com/openkcm/cmk/internal/repo"
	"github.com/openkcm/common-sdk/pkg/auth"
)

// AuthClientData contains a group and an identifier associated with an AuthClient
type AuthClientData struct {
	Group      *model.Group
	GroupID    string // For convenience. Just a string version of the Group.ID
	Identifier string
}

// ClientMapOpt are options which can be used, for example, when retrieving the
// ClientData from an AuthClient
type ClientMapOpt func(*auth.ClientData)

// GetClientMap gets the ClientMap from the AuthClient. This can be used to authenticate
func (cd AuthClientData) GetClientMap(opts ...ClientMapOpt) map[any]any {
	clientData := getClientData(cd.Identifier, []string{cd.Group.IAMIdentifier})

	for _, o := range opts {
		o(clientData)
	}

	return map[any]any{constants.ClientData: clientData}
}

// WithAdditionalGroup is an option when getting a ClientMap from an AuthClient.
// It adds an additional group to the ClientData Groups
func WithAdditionalGroup(groupName string) ClientMapOpt { // TODO use AddGroup
	return func(cd *auth.ClientData) {
		cd.Groups = append(cd.Groups, groupName)
	}
}

// WithOverriddenIdentifier is an option when getting a ClientMap from an AuthClient.
// It overrides the AuthClient Identifier. This can be used, for example,
// when testing for other users in (or not in) the AuthClient Group
func WithOverriddenIdentifier(identifier string) ClientMapOpt {
	return func(cd *auth.ClientData) {
		cd.Identifier = identifier
	}
}

// WithOverriddenGroup is an option when getting a ClientMap from an AuthClient.
// It overrides the AuthClient Groups. This can be used, for example,
// when testing for invalid groups for a given AuthClient identifier
func WithOverriddenGroup(numGroups int) ClientMapOpt { // TODO package level functions
	return func(cd *auth.ClientData) {
		cd.Groups = make([]string, numGroups)
		for i := range numGroups {
			cd.Groups[i] = uuid.NewString()
		}
	}
}

// WithAuthClientDataKC provides an option for the NewKeyConfig function
// This option will initialise the KeyConfig with the AuthClient Group
func WithAuthClientDataKC(authClient AuthClientData) KeyConfigOpt {
	return func(kc *model.KeyConfiguration) {
		kc.AdminGroup = *authClient.Group
		kc.AdminGroupID = authClient.Group.ID
	}
}

// AuthClientOpt are options which can be used with NewAuthClient
type AuthClientOpt func(*AuthClientData)

// NewAuthClient creates an AuthClient using random strings for values and creates
// in database the group
func NewAuthClient(ctx context.Context, tb testing.TB, r repo.Repo, opts ...AuthClientOpt) AuthClientData {
	authClient := newAuthClient(opts...)
	CreateTestEntities(ctx, tb, r, authClient.Group)
	return authClient
}

// GetAuthClientMap does the same as the NewAuthClient, except it returns the ClientMap directly.
// It can be used for simple tests when a separate AuthClient is not required
func GetAuthClientMap(ctx context.Context, tb testing.TB, r repo.Repo, opts ...AuthClientOpt) map[any]any {
	authClient := newAuthClient(opts...)
	CreateTestEntities(ctx, tb, r, authClient.Group)
	return authClient.GetClientMap()
}

// WithAuditorRole is an option when getting an AuthClient with NewAuthClient, or the ClientMap
// with GetAuthClientMap. It specifies TenantAuditorRole for the group
func WithAuditorRole() AuthClientOpt {
	return func(acd *AuthClientData) {
		acd.Group.Role = constants.TenantAuditorRole
	}
}

// WithKeyAdminRole is an option when getting an AuthClient with NewAuthClient, or the ClientMap
// with GetAuthClientMap. It specifies KeyAdminRole for the group
func WithKeyAdminRole() AuthClientOpt {
	return func(acd *AuthClientData) {
		acd.Group.Role = constants.KeyAdminRole
	}
}

// WithTenantAdminRole is an option when getting an AuthClient with NewAuthClient, or the ClientMap
// with GetAuthClientMap. It specifies TenantAdminRole for the group
func WithTenantAdminRole() AuthClientOpt {
	return func(acd *AuthClientData) {
		acd.Group.Role = constants.TenantAdminRole
	}
}

// WithIdentifier is an option when getting an AuthClient with NewAuthClient, or the ClientMap
// with GetAuthClientMap. It allows the default random value for the AuthClient Identifier to be
// overridden
func WithIdentifier(identifier string) AuthClientOpt {
	return func(acd *AuthClientData) {
		acd.Identifier = identifier
	}
}

// GetClientMap returns a client map created with the provided identifier and group names
// It does not create anything in the database
func GetClientMap(identifier string, groupNames []string) map[any]any {
	return map[any]any{constants.ClientData: getClientData(identifier, groupNames)}
}

// GetGrouplessClientMap returns a client map with a random identifier and no groupnames
// It does not create anything in the database
func GetGrouplessClientMap() map[any]any {
	return map[any]any{constants.ClientData: getClientData(uuid.NewString(), []string{})}
}

// GetInvalidClientMap returns a client map with random identifier and random groupnames
// It does not create anything in the database
func GetInvalidClientMap(opts ...ClientMapOpt) map[any]any {
	clientData := getClientData(uuid.NewString(), []string{uuid.NewString(), uuid.NewString()})
	return map[any]any{constants.ClientData: clientData}
}

func newAuthClient(opts ...AuthClientOpt) AuthClientData {
	group := NewGroup(func(g *model.Group) {
		g.ID = uuid.New()
		g.Name = uuid.NewString()
		g.IAMIdentifier = uuid.NewString()
		g.Role = constants.TenantAuditorRole
	})

	authClientData := AuthClientData{
		Group:      group,
		GroupID:    group.ID.String(),
		Identifier: uuid.NewString(),
	}

	for _, o := range opts {
		o(&authClientData)
	}

	return authClientData
}

func getClientData(identifier string, groupNames []string) *auth.ClientData {
	return &auth.ClientData{
		Identifier: identifier,
		Groups:     groupNames,
	}
}
