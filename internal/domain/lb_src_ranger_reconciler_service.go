package domain

import (
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"time"
)

type LbSrcRangerReconcilerService interface {
	Id() LbSrcRangerId
	MkLbRangerReadOps() LbSrcRangersReadOps
	MkLbServicesReadOps() LbServicesReadOps
	MkLogger() logr.Logger
	MkCidrsFetcher() CidrsFetcher
	Now() time.Time
}

type ReconcilerResult struct {
	Requeue      bool
	RequeueAfter time.Duration
}

// This is our main business logic, the reward for rigourously following
// domain modeling is to be able to write the logic without dealing with
// k8s directly.
func Reconcile(l LbSrcRangerReconcilerService) (ReconcilerResult, error) {
	lbId := l.Id()
	lbRangerReadOps := l.MkLbRangerReadOps()
	lbServicesReadOps := l.MkLbServicesReadOps()
	logger := l.MkLogger().WithValues("service", "LbSrcRangerReconcilerService")

	ranger, getErr := lbRangerReadOps.Get(&lbId)
	if getErr != nil {
		logger.Error(getErr, "Unable to Fetch LoadBalancerRanger")
		if getErr.IsNotFound {
			return ReconcilerResult{}, nil
		} else {
			return ReconcilerResult{Requeue: true}, getErr
		}
	}

	services, listErr := lbServicesReadOps.FilterFor(&ranger)
	if listErr != nil {
		return ReconcilerResult{Requeue: true}, listErr
	}

	if len(services) > 0 {
		cidrsFetcher := l.MkCidrsFetcher()
		if cidrs, err := cidrsFetcher.Fetch(&ranger.Spec.SrcIPUrls); err != nil {
			return ReconcilerResult{Requeue: true}, err
		} else {
			updatedCount := 0
			for _, service := range services {
				// Only update if needed
				if !cmp.Equal(service.LbSrcRanges, cidrs) {
					service.LbSrcRanges = cidrs
					if err := service.UpdateOps.UpdateCidrs(&cidrs); err != nil {
						return ReconcilerResult{Requeue: true}, err
					} else {
						updatedCount += 1
					}
				}
			}

			status := LbSrcRangerStatus{
				LastUpdatedCount: updatedCount,
				LastRunAt:        l.Now(),
			}
			if err := ranger.UpdateOps.UpdateStatus(&status); err != nil {
				logger.Error(err, "Unable to update Ranger status")
				return ReconcilerResult{Requeue: true}, err
			}
		}
	}
	logger.Info("Reconciled successfully", "id", lbId)
	return ReconcilerResult{RequeueAfter: buildRequeueAfter(ranger.Spec.UpdateEvery, logger)}, nil
}

func buildRequeueAfter(updateEvery time.Duration, logger logr.Logger) time.Duration {
	if updateEvery.Nanoseconds() <= 0 {
		logger.Info("Update every duration was < 0, so using reasonable default of 1 minute", "update_every", updateEvery)
		return time.Minute
	} else {
		return updateEvery
	}
}
