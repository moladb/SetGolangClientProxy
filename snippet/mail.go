// mail info
package snippet

type EmailList struct {
	EmailList []string `json:"emailList"`
}

type SendMailResponse struct {
	Result     bool      `json:"result"`
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
	Info       EmailList `json:"info"`
}
