package firewalld

type Rule struct {
	Protocol    string
	Port        string
	SourceIps   []string
	Direction   string
	Description string
}

