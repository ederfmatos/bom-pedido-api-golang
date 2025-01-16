package entity

type MerchantPaymentGatewayConfig struct {
	MerchantID     string `json:"merchantId" bson:"merchantId"`
	PaymentGateway string `json:"paymentGateway" bson:"paymentGateway"`
	Credential     string `json:"credential" bson:"credential"`
}

func NewMerchantPaymentGatewayConfig(merchantID, paymentGateway, credential string) *MerchantPaymentGatewayConfig {
	return &MerchantPaymentGatewayConfig{MerchantID: merchantID, PaymentGateway: paymentGateway, Credential: credential}
}
