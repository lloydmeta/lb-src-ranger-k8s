package k8s

import (
	lbsrcrangerv1beta1 "github.com/lloydmeta/lb-src-ranger-k8s/api/v1beta1"
	"github.com/lloydmeta/lb-src-ranger-k8s/internal/domain"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func MkLbSrcRangerReadOps(components *lbSrcRangerReconcilerComponents) domain.LbSrcRangersReadOps {
	return &lbSrcRangersReadOpsImpl{
		components: components,
	}
}

type lbSrcRangersReadOpsImpl struct {
	components *lbSrcRangerReconcilerComponents
}

type lbSrcRangersUpdateOpsImpl struct {
	components            *lbSrcRangerReconcilerComponents
	k8sLoadBalancerRanger *lbsrcrangerv1beta1.LbSrcRanger
}

func (l *lbSrcRangersReadOpsImpl) Get(id *domain.LbSrcRangerId) (domain.LbSrcRanger, *domain.LbSrcRangersReadGetErr) {
	var ranger = lbsrcrangerv1beta1.LbSrcRanger{}
	namespacedName := types.NamespacedName{
		Namespace: id.Namespace,
		Name:      id.Name,
	}
	if err := l.components.Get(l.components.ctx, namespacedName, &ranger); err != nil {
		return domain.LbSrcRanger{}, &domain.LbSrcRangersReadGetErr{
			IsNotFound: apierrs.IsNotFound(err),
			Underlying: err,
		}
	} else {
		updateOps := lbSrcRangersUpdateOpsImpl{
			components:            l.components,
			k8sLoadBalancerRanger: &ranger,
		}
		domainLbRanger := domain.LbSrcRanger{
			Id: domain.LbSrcRangerId(namespacedName),
			Spec: domain.LbSrcRangerSpec{
				TargetLabels: ranger.Spec.TargetLabels,
				UpdateEvery:  ranger.Spec.UpdateEvery.Duration,
				SrcIPUrls:    ranger.Spec.SrcIPUrls,
			},
			Status: domain.LbSrcRangerStatus{
				LastUpdatedCount: ranger.Status.LastUpdatedCount,
				LastRunAt:        ranger.Status.LastRunAt.Time,
			},
			UpdateOps: &updateOps,
		}
		return domainLbRanger, nil
	}
}

func (l *lbSrcRangersUpdateOpsImpl) UpdateStatus(status *domain.LbSrcRangerStatus) error {

	l.k8sLoadBalancerRanger.Status.LastUpdatedCount = status.LastUpdatedCount
	l.k8sLoadBalancerRanger.Status.LastRunAt = v1.NewTime(status.LastRunAt)

	if err := l.components.Status().Update(l.components.ctx, l.k8sLoadBalancerRanger); err != nil {
		return err
	} else {
		return nil
	}
}
