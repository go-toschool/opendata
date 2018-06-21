package sanctuary

type ExtractionRequest struct {
	UserToken string `json:"user_token,omitempty"`
}

type AuthRequest struct {
	PartnerToken string `json:"partner_token,omitempty"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
}

type Response struct {
}

type Balance struct {
}

type Account struct {
}
