package github

type ErrorResponse struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Errors           []struct {
		Message string `json:"message"`
	} `json:"errors"`
}
