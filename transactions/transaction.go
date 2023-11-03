package transactions

import (
  "time"
)

type transaction interface {
  Output() string
  GetCompletedDate() time.Time
}

