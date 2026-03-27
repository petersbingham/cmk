package authz

import "github.com/openkcm/cmk/internal/constants"

var RepoInternalPolicies = make(map[constants.InternalRole][]BasePolicy[constants.InternalRole,
	RepoResourceTypeName, RepoAction])

type internalRepoPolicies struct {
	Roles    []constants.InternalRole
	Policies []BasePolicy[constants.InternalRole, RepoResourceTypeName, RepoAction]
}

var InternalRepoPolicyData = internalRepoPolicies{
	Roles: []constants.InternalRole{
		constants.InternalTenantProvisioningRole,
	},
	Policies: []BasePolicy[constants.InternalRole, RepoResourceTypeName, RepoAction]{
		NewPolicy(
			"InternalTenantProvisioning",
			constants.InternalTenantProvisioningRole,
			[]BaseResourceType[RepoResourceTypeName, RepoAction]{
				{
					ID: RepoResourceTypeTenant,
					Actions: []RepoAction{
						RepoActionCreate,
					},
				},
				{
					ID: RepoResourceTypeGroup,
					Actions: []RepoAction{
						RepoActionCreate,
					},
				},
			},
		),
	},
}

func init() {
	// Index policies by role for fast lookup
	RepoInternalPolicies = make(map[constants.InternalRole][]BasePolicy[
		constants.InternalRole, RepoResourceTypeName, RepoAction])
	for _, policy := range InternalRepoPolicyData.Policies {
		RepoInternalPolicies[policy.Role] = append(RepoInternalPolicies[policy.Role], policy)
	}
}
