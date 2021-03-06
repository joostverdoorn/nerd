package v1payload

// CreateBillingPackageInput is the input for assigning a billing package to a project.
// This results in the creation of a quota in the right namespace.
type CreateBillingPackageInput struct {
	BillingPackageID string `json:"billing_package_id"`
	RequestsCPU      string `json:"requests_cpu"`
	RequestsMemory   string `json:"requests_memory"`
}

// CreateBillingPackageOutput is the output from assigning a billing package to a project.
type CreateBillingPackageOutput struct {
	ProjectID        string `json:"project_id" valid:"required"`
	BillingPackageID string `json:"billing_package_id" valid:"required"`
	RequestsCPU      string `json:"requests_cpu"`
	RequestsMemory   string `json:"requests_memory"`
}

// UpdateBillingPackageInput is the input for updating the billing package capacity
type UpdateBillingPackageInput struct {
	OnDemand       bool   `json:"on_demand"`
	RequestsCPU    string `json:"requests_cpu"`
	RequestsMemory string `json:"requests_memory"`
}

// UpdateBillingPackageOutput is the output for updating the billing package capacity
type UpdateBillingPackageOutput struct {
	ProjectID        string `json:"project_id" valid:"required"`
	BillingPackageID string `json:"billing_package_id" valid:"required"`
	RequestsCPU      string `json:"requests_cpu"`
	RequestsMemory   string `json:"requests_memory"`
}

// RemoveBillingPackageInput is the input for removing a billing package from a project
type RemoveBillingPackageInput struct {
}

// RemoveBillingPackageOutput is the output from removing a billing package from a project
type RemoveBillingPackageOutput struct {
}

// DeleteBillingPackageInput is the input for deleting a billing package
type DeleteBillingPackageInput struct {
}

// DeleteBillingPackageOutput is the output from deleting a billing package
type DeleteBillingPackageOutput struct {
}

//BillingPackageSummary is summary of a billing package
type BillingPackageSummary struct {
	RequestsCPU      string `json:"requests_cpu" valid:"required"`
	RequestsMemory   string `json:"requests_memory" valid:"required`
	BillingPackageID string `json:"billing_package_id" valid:"required"`
}

// ListBillingPackagesInput is the input for listing billing packages.
type ListBillingPackagesInput struct {
}

// ListBillingPackagesOutput is the output from listing billing packages of a project
type ListBillingPackagesOutput struct {
	ProjectID       string                   `json:"project_id" valid:"required"`
	BillingPackages []*BillingPackageSummary `json:"billing_packages" valid:"required"`
	Total           *Resource
	Used            *Resource
}

// Resource is a general struct that will be used in our list payloads.
type Resource struct {
	RequestsCPU    string `json:"requests_cpu" valid:"required"`
	RequestsMemory string `json:"requests_memory" valid:"required`
	LimitsCPU      string `json:"limits_cpu" valid:"required"`
	LimitsMemory   string `json:"limits_memory" valid:"required"`
}
