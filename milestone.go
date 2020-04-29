package github

type Milestone struct {
	URL          string `json:"url"`
	HtmlURL      string `json:"html_url"`
	LabelsURL    string `json:"labels_url"`
	ID           int64  `json:"id"`
	NodeID       string `json:"node_id"`
	Number       int64  `json:"number"`
	State        string `json:"state"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Creator      User   `json:"creator"`
	OpenIssues   int64  `json:"open_issues"`
	ClosedIssues int64  `json:"closed_issues"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	ClosedAt     string `json:"closed_at"`
	DueOn        string `json:"due_on"`
}
