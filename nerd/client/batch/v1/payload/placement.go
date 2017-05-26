package v1payload

//PlaceProjectInput is input for queue creation
type PlaceProjectInput struct {
	ProjectID string `json:"project_id" valid:"required"`
	Host      string `json:"host" valid:"required"`
	Token     string `json:"token"`
	CAPem     string `json:"ca_pem"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Insecure  bool   `json:"insecure"`
}

//PlaceProjectOutput is output for queue creation
type PlaceProjectOutput struct {
}

//ExpelProjectInput is input for queue creation
type ExpelProjectInput struct {
	ProjectID string `json:"project_id" valid:"required"`
}

//ExpelProjectOutput is output for queue creation
type ExpelProjectOutput struct{}
