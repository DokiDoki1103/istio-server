package graph

type Node struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	NodeType     string `json:"nodeType"`
	Workload     string `json:"workload"`
	Version      string `json:"version"`
	Service      string `json:"service"`
	TenantID     string `json:"tenant_id"`
	ServiceID    string `json:"service_id"`
	ServiceAlias string `json:"service_alias"`
	RainBongApp  string `json:"rainbond_app"`
	Instance     string `json:"instance"`
	Protocol     string `json:"protocol"`
	RequestTime  string `json:"request_time"`
	RequestRate  string `json:"request_rate"`
	//RequestErrorRate string `json:"request_error_rate"`
}

type Edge struct {
	ID          string `json:"id"`
	Source      string `json:"source"`
	Dist        string `json:"dist"`
	Protocol    string `json:"protocol"`
	RequestTime string `json:"request_time"`
	RequestRate string `json:"request_rate"`
}

type Graph struct {
	Nodes map[string]*Node `json:"nodes"`
	Edges []*Edge          `json:"edges"`
}
