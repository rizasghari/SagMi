package models

type Field struct {
	Short bool   `json:"short"`
	Title string `json:"title"`
	Value string `json:"value"`
}

type Action struct {
	Text string `json:"text"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Attachment struct {
	Color              string   `json:"color"`
	HealthCheckService string   `json:"service"`
	Fields             []Field  `json:"fields"`
	Actions            []Action `json:"actions,omitempty"`
}

type MessageBody struct {
	Attachments []Attachment `json:"attachments"`
	Text        string       `json:"text"`
	Markdown    bool         `json:"mrkdwn"`
	UserName    string       `json:"username"`
}

type AlarmData struct {
	Environment        string
	AppName            string
	Service            string
	ServiceURL         string
	Response           string
	HealthCheckService string
}
