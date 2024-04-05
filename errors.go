package linkpreview

import "fmt"

type ErrNonHTTPScheme struct {
	url string
}

func (e *ErrNonHTTPScheme) Error() string {
	return fmt.Sprintf("non-HTTP/HTTPS scheme: %s", e.url) 
}

type ErrNonSupportedResponse struct {
	respCode int
}

func (e *ErrNonSupportedResponse) Error() string {
	return fmt.Sprintf("non-supported response code: %d", e.respCode)
}

type ErrNonSupportedContentType struct {
	contentType string
}

func (e *ErrNonSupportedContentType) Error() string {
	return fmt.Sprintf("non-supported content type: %s", e.contentType)
}