package couchdb

// ClientOption defines a function the can modify the client parameters.
type ClientOption func(*Client) error

// WithUsername returns an option that sets the client username.
func WithUsername(value string) ClientOption {
	return func(c *Client) error {
		c.username = value
		return nil
	}
}

// WithPassword returns an option that sets the client password.
func WithPassword(value string) ClientOption {
	return func(c *Client) error {
		c.password = value
		return nil
	}
}
