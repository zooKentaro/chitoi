package helper

import "time"

// IsSameDay is XXX
func IsSameDay(a, b time.Time) bool {
    d := time.Duration(-a.Hour()) * time.Hour
    BeginningOfA := a.Truncate(time.Hour).Add(d)

    d = time.Duration(-b.Hour()) * time.Hour
    BeginningOfB := b.Truncate(time.Hour).Add(d)

    return BeginningOfA.Unix() == BeginningOfB.Unix()
}
