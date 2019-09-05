package k8s

import (
	"github.com/lloydmeta/lb-src-ranger-k8s/internal/domain"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func MkLbServicesReadOps(components *lbSrcRangerReconcilerComponents) domain.LbServicesReadOps {
	return &lbServicesReadOpsImpl{
		components: components,
	}
}

type lbServicesReadOpsImpl struct {
	components *lbSrcRangerReconcilerComponents
}

type lbServicesUpdateOpsImpl struct {
	components *lbSrcRangerReconcilerComponents
	k8sService *v1.Service
}

func (l *lbServicesUpdateOpsImpl) UpdateCidrs(cidrs *[]domain.Cidr) error {
	strCidrs := make([]string, len(*cidrs))
	for i, cidr := range *cidrs {
		strCidrs[i] = string(cidr)
	}
	l.k8sService.Spec.LoadBalancerSourceRanges = strCidrs
	return l.components.Update(l.components.ctx, l.k8sService)
}

func (l2 *lbServicesReadOpsImpl) FilterFor(l *domain.LbSrcRanger) ([]domain.LbService, error) {
	var loadBalancerServices = &v1.ServiceList{}
	loadBalancerServicesQuery := buildListQuery(l)
	if err := l2.components.List(l2.components.ctx, loadBalancerServices, loadBalancerServicesQuery...); err != nil {
		return nil, err
	} else {
		list := make([]domain.LbService, len(loadBalancerServices.Items))
		for i, s := range loadBalancerServices.Items {
			updateOps := lbServicesUpdateOpsImpl{
				components: l2.components,
				k8sService: &s,
			}
			lbSrcRngs := make([]domain.Cidr, len(s.Spec.LoadBalancerSourceRanges))
			for i, r := range s.Spec.LoadBalancerSourceRanges {
				lbSrcRngs[i] = domain.Cidr(r)
			}
			lbService := domain.LbService{
				Name:        s.Name,
				LbSrcRanges: lbSrcRngs,
				UpdateOps:   &updateOps,
			}
			list[i] = lbService
		}
		return list, nil
	}
}

func buildListQuery(l *domain.LbSrcRanger) []client.ListOption {
	return []client.ListOption{
		client.InNamespace(l.Id.Namespace),
		client.MatchingLabels(l.Spec.TargetLabels),
	}
}
