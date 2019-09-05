package domain

import (
	"fmt"
	"time"
)

type LbSrcRangerId struct {
	Namespace string
	Name      string
}

type LbSrcRanger struct {
	Id        LbSrcRangerId
	Spec      LbSrcRangerSpec
	Status    LbSrcRangerStatus
	UpdateOps LbSrcRangersUpdateOps
}

type LbSrcRangerSpec struct {
	TargetLabels map[string]string
	UpdateEvery  time.Duration
	SrcIPUrls    []string
}

type LbSrcRangerStatus struct {
	LastUpdatedCount int
	LastRunAt        time.Time
}

type LbSrcRangersReadOps interface {
	Get(id *LbSrcRangerId) (LbSrcRanger, *LbSrcRangersReadGetErr)
}

type LbSrcRangersUpdateOps interface {
	UpdateStatus(status *LbSrcRangerStatus) error
}

type LbSrcRangersReadGetErr struct {
	IsNotFound bool
	Underlying error
}

func (l *LbSrcRangersReadGetErr) Error() string {
	return fmt.Sprintf("LbSrcRangersReadGetErr { IsNotFound: %v, underlying: %v }", l.IsNotFound, l.Underlying.Error())
}
