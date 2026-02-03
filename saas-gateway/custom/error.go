package custom

import "fmt"

type Error struct {
	Code        string `json:"code"`
	InternalKey string `json:"internalKey"`
	Message     string `json:"message"`
	ParentError error
}

func (b *Error) Error() string {
	return fmt.Sprintf("Code: %s, Internal Key: %s, Message: %s, Parent Error: %v", b.Code, b.InternalKey, b.Message, b.ParentError)
}
