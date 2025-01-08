package merchant

type PaymentGatewayConfig struct {
	MerchantID     string `json:"merchantId" bson:"merchantId"`
	PaymentGateway string `json:"paymentGateway" bson:"paymentGateway"`
	Credential     string `json:"credential" bson:"credential"`
}

func NewPaymentGatewayConfig(merchantID string, paymentGateway string, credential string) *PaymentGatewayConfig {
	return &PaymentGatewayConfig{MerchantID: merchantID, PaymentGateway: paymentGateway, Credential: credential}
}
