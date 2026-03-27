package constants

const (
	//nolint:gosec
	AuthTypeSecret      = "AUTH_TYPE_SECRET"
	AuthTypeCertificate = "AUTH_TYPE_CERTIFICATE"

	TenantAdminGroup   string = "TenantAdministrator"
	TenantAuditorGroup string = "TenantAuditor"

	KeyAdminRole      BusinessRole = "KEY_ADMINISTRATOR"
	TenantAdminRole   BusinessRole = "TENANT_ADMINISTRATOR"
	TenantAuditorRole BusinessRole = "TENANT_AUDITOR"

	InternalTenantProvisioningRole InternalRole = "INTERNAL_TENANT_PROVISIONING"
)

type BusinessRole string
type InternalRole string
