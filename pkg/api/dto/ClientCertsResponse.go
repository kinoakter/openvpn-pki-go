package dto

type ClientCertResponse struct {
	CaCrt               string `json:"caCrt"`
	ClientCrt           string `json:"clientCrt"`
	ClientKey           string `json:"clientKey"`
	TlsCryptClientV2Key string `json:"tlsCryptClientV2Key"`
	ExpiresAt           int64  `json:"expiresAt"`
}
