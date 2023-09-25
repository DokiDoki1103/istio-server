package graph

import (
	"crypto/md5"
	"fmt"
	"github.com/prometheus/common/model"
	mo "istio-server/model"
)

func AddNode(value model.Value, nodeMap map[string]*Node) {
	switch t := value.Type(); t {
	case model.ValVector:
		vector := value.(model.Vector)
		for _, s := range vector {
			m := s.Metric
			key := string(m["service_istio_io_canonical_name"])
			n := &Node{
				Workload:     string(m["service_istio_io_canonical_name"]),
				NodeType:     "",
				Version:      string(m["version"]),
				TenantID:     string(m["tenant_id"]),
				Service:      string(m["destination_service"]),
				ServiceID:    string(m["service_id"]),
				ServiceAlias: string(m["service_alias"]),
				RainBongApp:  string(m["rainbond_app"]),
				Instance:     string(m["instance"]),
			}

			var tenantService mo.TenantService
			mo.DB.Model(&tenantService).Select("service_cname").Where("service_id = ?", n.ServiceID).Take(&tenantService)
			n.Name = tenantService.ServiceCName
			nodeMap[key] = n
		}

	}
}

func AddEdge(value model.Value, edges []*Edge, nodeMap map[string]*Node) []*Edge {
	switch t := value.Type(); t {
	case model.ValVector:
		vector := value.(model.Vector)
		idMap := make(map[string]bool)
		for _, s := range vector {
			m := s.Metric
			source := string(m["source_workload"])
			dist := string(m["destination_workload"])
			id := fmt.Sprintf("%x", md5.Sum([]byte(source+dist)))
			if "unknown" == source || "unknown" == dist || idMap[id] {
				continue
			}
			idMap[id] = true
			e := &Edge{
				ID:       id,
				Protocol: string(m["request_protocol"]),
				Source:   source,
				Dist:     dist,
			}
			_, exists := nodeMap[e.Source]
			_, exists2 := nodeMap[e.Dist]

			if exists2 && exists && e.Source != e.Dist {
				edges = append(edges, e)
				nodeMap[e.Dist].Protocol = e.Protocol
			}
		}
	}
	return edges
}
