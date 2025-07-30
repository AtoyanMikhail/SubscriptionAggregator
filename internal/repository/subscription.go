package repository

import "time"

type Subscription struct {
    ID          int       `db:"id" json:"id"`
    Price       int       `db:"price" json:"price"`
    UserID      string    `db:"user_id" json:"user_id"`
    ServiceName string    `db:"service_name" json:"service_name"`
    StartDate   time.Time `db:"start_date" json:"start_date"`
    EndDate     *time.Time `db:"end_date" json:"end_date,omitempty"` // Nullable
}
