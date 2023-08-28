package entity

import "time"

type Operation string

const (
	SEGMENT_ADDED   Operation = "segment_added"
	SEGMENT_REMOVED Operation = "segment_removed"
)

// таблица  many to many
type UsersSegments struct {
	User    int    `db:"user"`
	Segment string `db:"segment"`
}

type UsersSegmentsStats struct {
	User       int       `db:"user"`
	Segment    string    `db:"segment"`
	Created_at time.Time `db:"creaed_at"`
	Operation  Operation `db:"operation"`
}
