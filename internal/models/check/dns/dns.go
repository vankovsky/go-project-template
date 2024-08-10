package dns

import "net"

type Host struct {
	ID   int
	Name string
}

type DNSTypeAddr struct {
	IPs  []net.IPAddr
	Type string
}

type CheckDNSDBHander interface {
	Hosts(limit int) ([]Host, error)
	SaveHostRecords(DNSTypeAddr) error
}

// import "go-project-template/internal/pkg/models"

// const (
// 	OpAdd    = "add"
// 	OpDelete = "delete"
// 	OpUpdate = "update"
// )

// type Record struct {
// 	IP       models.IP
// 	Priority int
// 	TTL      int
// }

// type RecordLog struct {
// 	NewRecord  *Record
// 	PrevRecord *Record
// 	Op         string `json:"op"`
// }

// func (d *Record) IsEqual(record *Record) bool {
// 	return d.IP == record.IP &&
// 		d.Priority == record.Priority &&
// 		d.TTL == record.TTL
// }

// func DeltaLog(old, new map[models.IP]Record) []RecordLog {
// 	var diff []RecordLog

// 	if len(old) == 0 {
// 		for _, r := range new {
// 			diff = append(diff, RecordLog{
// 				NewRecord: &r,
// 				Op:        OpAdd,
// 			})
// 		}
// 		return diff
// 	}

// 	if len(new) == 0 {
// 		for _, r := range old {
// 			diff = append(diff, RecordLog{
// 				NewRecord: &r,
// 				Op:        OpDelete,
// 			})
// 		}

// 		return diff
// 	}

// 	for newIP, newR := range new {
// 		if oldR, ok := old[newIP]; ok {

// 			if !newR.IsEqual(&oldR) {
// 				newR := newR
// 				diff = append(diff, RecordLog{
// 					NewRecord:  &newR,
// 					PrevRecord: &oldR,
// 					Op:         OpUpdate,
// 				})
// 			}
// 			delete(new, newIP)
// 			delete(old, newIP)
// 		}
// 	}

// 	for _, r := range new {
// 		diff = append(diff, RecordLog{
// 			NewRecord: &r,
// 			Op:        OpAdd,
// 		})
// 	}

// 	for _, r := range old {
// 		diff = append(diff, RecordLog{
// 			PrevRecord: &r,
// 			Op:         OpDelete,
// 		})
// 	}

// 	return diff
// }
