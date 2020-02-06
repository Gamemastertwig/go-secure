package web

import (
	"html/template"
	"net/http"
)

// ContentTemplate is a struct of variables that will be used to display content
type ContentTemplate struct {
	RproxyBody   string
	FirewallBody string
	LBalBody     string
	IDSBody      string
	LogBody      string
}

// Content generates HTML template to display desired content
func Content(w http.ResponseWriter, r *http.Request) {
	var content ContentTemplate

	content.RproxyBody = "Reverse Proxy content get loaded here..."
	content.FirewallBody = "Firewall content get loaded here..."
	content.LBalBody = "[LB] Comming Soon..."
	content.IDSBody = "[IDS] Comming Soon..."
	content.LogBody = "Logger content gets loaded here..."

	template := template.Must(template.ParseFiles("web/index.html"))
	template.Execute(w, content)
}
