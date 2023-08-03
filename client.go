package mercadopago

type client struct {
	token string
}

func NewClient(token string) *client {
	return &client{
		token: token,
	}
}
