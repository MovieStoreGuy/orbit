package types

type Machine struct {
	Name       string            `json:"name"`
	Labels     map[string]string `json:"labels"`
	ExternalIP []string          `json:"externalIPs"`
	Zone       string            `json:"zone"`
	Status     string            `json:"status"`
}
