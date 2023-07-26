package bms

type UnexpectedResponseError struct {
	Message   string
	Want, Got any
}

func (u UnexpectedResponseError) Error() string {
	if u.Message == "" {
		return "unexpected response"
	}
	return u.Message
}
