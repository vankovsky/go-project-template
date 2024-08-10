package dns

// import (
// 	"testing"

// 	"go-project-template/internal/pkg/models"

// 	"github.com/stretchr/testify/assert"
// )

// func TestDNSRecord_IsEqual(t *testing.T) {
// 	type fields struct {
// 		IP       models.IP
// 		Priority int
// 		TTL      int
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		record *Record
// 		want   bool
// 	}{
// 		{
// 			name: "non-matched: IP",
// 			fields: fields{
// 				IP:       "1.2.3.4",
// 				Priority: 5,
// 				TTL:      6,
// 			},
// 			record: &Record{
// 				IP:       "1.2.3.5",
// 				Priority: 5,
// 				TTL:      6,
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "non-matched: TTL",
// 			fields: fields{
// 				IP:       "1.2.3.4",
// 				Priority: 5,
// 				TTL:      6,
// 			},
// 			record: &Record{
// 				IP:       "1.2.3.4",
// 				Priority: 5,
// 				TTL:      5,
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "non-matched: Priority",
// 			fields: fields{
// 				IP:       "1.2.3.4",
// 				Priority: 5,
// 				TTL:      6,
// 			},
// 			record: &Record{
// 				IP:       "1.2.3.4",
// 				Priority: 10,
// 				TTL:      6,
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "matched",
// 			fields: fields{
// 				IP:       "1.2.3.4",
// 				Priority: 5,
// 				TTL:      6,
// 			},
// 			record: &Record{
// 				IP:       "1.2.3.4",
// 				Priority: 5,
// 				TTL:      6,
// 			},
// 			want: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			d := &Record{
// 				IP:       tt.fields.IP,
// 				Priority: tt.fields.Priority,
// 				TTL:      tt.fields.TTL,
// 			}
// 			if got := d.IsEqual(tt.record); got != tt.want {
// 				t.Errorf("IsEqual() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestDeltaLog(t *testing.T) {
// 	t.Run("no old", func(t *testing.T) {
// 		oldRecords := map[models.IP]Record{}

// 		newRecords := map[models.IP]Record{
// 			"127.0.0.3": {IP: "127.0.0.3", TTL: 25, Priority: 35},
// 		}

// 		result := DeltaLog(oldRecords, newRecords)
// 		assert.Len(t, result, 1)
// 		r := result[0]

// 		assert.Equal(t, r.Op, OpAdd)
// 		assert.Nil(t, r.PrevRecord)
// 		assert.Equal(t, *r.NewRecord, Record{
// 			IP:       "127.0.0.3",
// 			TTL:      25,
// 			Priority: 35,
// 		})
// 	})

// 	t.Run("no new", func(t *testing.T) {
// 		oldRecords := map[models.IP]Record{
// 			"127.0.0.1": {IP: "127.0.0.1", TTL: 5, Priority: 10},
// 		}

// 		newRecords := map[models.IP]Record{}

// 		result := DeltaLog(oldRecords, newRecords)
// 		assert.Len(t, result, 1)
// 		r := result[0]

// 		assert.Equal(t, r.Op, OpDelete)
// 		assert.Nil(t, r.PrevRecord)
// 		assert.Equal(t, *r.NewRecord, Record{
// 			IP:       "127.0.0.1",
// 			TTL:      5,
// 			Priority: 10,
// 		})
// 	})

// 	t.Run("add", func(t *testing.T) {
// 		oldRecords := map[models.IP]Record{
// 			"127.0.0.1": {IP: "127.0.0.1", TTL: 5, Priority: 10},
// 			"127.0.0.2": {IP: "127.0.0.2", TTL: 15, Priority: 15},
// 		}

// 		newRecords := map[models.IP]Record{
// 			"127.0.0.1": {IP: "127.0.0.1", TTL: 5, Priority: 10},
// 			"127.0.0.2": {IP: "127.0.0.2", TTL: 15, Priority: 15},
// 			"127.0.0.3": {IP: "127.0.0.3", TTL: 25, Priority: 35},
// 		}

// 		result := DeltaLog(oldRecords, newRecords)
// 		assert.Len(t, result, 1)
// 		r := result[0]

// 		assert.Equal(t, r.Op, OpAdd)
// 		assert.Nil(t, r.PrevRecord)
// 		assert.Equal(t, *r.NewRecord, Record{
// 			IP:       "127.0.0.3",
// 			TTL:      25,
// 			Priority: 35,
// 		})
// 	})

// 	t.Run("delete", func(t *testing.T) {
// 		oldRecords := map[models.IP]Record{
// 			"127.0.0.1": {IP: "127.0.0.1", TTL: 5, Priority: 10},
// 			"127.0.0.2": {IP: "127.0.0.2", TTL: 15, Priority: 15},
// 			"127.0.0.3": {IP: "127.0.0.3", TTL: 25, Priority: 35},
// 		}

// 		newRecords := map[models.IP]Record{
// 			"127.0.0.2": {IP: "127.0.0.2", TTL: 15, Priority: 15},
// 			"127.0.0.3": {IP: "127.0.0.3", TTL: 25, Priority: 35},
// 		}

// 		result := DeltaLog(oldRecords, newRecords)
// 		assert.Len(t, result, 1)
// 		r := result[0]

// 		assert.Equal(t, r.Op, OpDelete)
// 		assert.Nil(t, r.NewRecord)
// 		assert.Equal(t, *r.PrevRecord, Record{
// 			IP:       "127.0.0.1",
// 			TTL:      5,
// 			Priority: 10,
// 		})
// 	})

// 	t.Run("update", func(t *testing.T) {
// 		oldRecords := map[models.IP]Record{
// 			"127.0.0.1": {IP: "127.0.0.1", TTL: 5, Priority: 10},
// 			"127.0.0.2": {IP: "127.0.0.2", TTL: 15, Priority: 15},
// 			"127.0.0.3": {IP: "127.0.0.3", TTL: 25, Priority: 35},
// 		}

// 		newRecords := map[models.IP]Record{
// 			"127.0.0.1": {IP: "127.0.0.1", TTL: 5, Priority: 45},
// 			"127.0.0.2": {IP: "127.0.0.2", TTL: 15, Priority: 15},
// 			"127.0.0.3": {IP: "127.0.0.3", TTL: 25, Priority: 35},
// 		}

// 		result := DeltaLog(oldRecords, newRecords)
// 		assert.Len(t, result, 1)
// 		r := result[0]

// 		assert.Equal(t, r.Op, OpUpdate)
// 		assert.Equal(t, *r.PrevRecord, Record{
// 			IP:       "127.0.0.1",
// 			TTL:      5,
// 			Priority: 10,
// 		})
// 		assert.Equal(t, *r.NewRecord, Record{
// 			IP:       "127.0.0.1",
// 			TTL:      5,
// 			Priority: 45,
// 		})
// 	})
// }
