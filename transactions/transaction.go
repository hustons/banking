package transactions

import (
  "time"
)

type transaction interface {
  Output()
  GetCompletedDate() time.Time
}

