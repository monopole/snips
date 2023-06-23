package myjira

type issueSearchRequest struct {
	Jql        string   `json:"jql,omitempty" :"jql" :"jql"`
	StartAt    int      `json:"startAt",omitempty`
	MaxResults int      `json:"maxResults,omitempty" :"maxResults" :"maxResults"`
	Fields     []string `json:"fields,omitempty" :"fields"`
	Expand     []string `json:"expand,omitempty" :"expand"`
}

type issueSearchResponse struct {
	Expand    string        `json:"expand,omitempty"`
	StartAt   int           `json:"startAt,omitempty"`
	MxResults int           `json:"maxResults,omitempty"`
	Total     int           `json:"total,omitempty"`
	Issues    []issueRecord `json:"issues,omitempty"`
}

type project struct {
	Self string `json:"self,omitempty"`
	Id   string `json:"id,omitempty"`
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

type user struct {
	Self         string `json:"self,omitempty"`
	Name         string `json:"name,omitempty"`
	Key          string `json:"key,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
}

type issueDetails struct {
	Summary     string   `json:"summary,omitempty"`
	Creator     user     `json:"creator"`
	Description string   `json:"description,omitempty"`
	Project     project  `json:"project"`
	Reporter    user     `json:"reporter"`
	Assignee    user     `json:"assignee"`
	Update      string   `json:"update,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

type issueRecord struct {
	Expand string       `json:"expand,omitempty"`
	Id     string       `json:"id,omitempty"`
	Self   string       `json:"self,omitempty"`
	Key    string       `json:"key,omitempty"`
	Fields issueDetails `json:"fields"`
}
