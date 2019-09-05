package k8s

import (
	"context"
	"errors"
	"github.com/lloydmeta/lb-src-ranger-k8s/internal/domain"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildListQuery(t *testing.T) {
	ranger := domain.LbSrcRanger{
		Id: domain.LbSrcRangerId{
			Namespace: "namespace",
		},
		Spec: domain.LbSrcRangerSpec{
			TargetLabels: map[string]string{
				"label": "value",
			},
		},
		Status:    domain.LbSrcRangerStatus{},
		UpdateOps: nil,
	}
	filter := buildListQuery(&ranger)
	expected := []client.ListOption{
		client.InNamespace("namespace"),
		client.MatchingLabels(map[string]string{
			"label": "value",
		}),
	}
	assert.Equal(t, expected, filter)
}

func TestLbServicesReadOpsImpl_FilterFor_ok(t *testing.T) {
	k8sClient := mockK8sClient{}
	components := lbSrcRangerReconcilerComponents{Client: &k8sClient}
	readOps := MkLbServicesReadOps(&components)
	service := v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name: "my-service",
		},
		Spec: v1.ServiceSpec{
			LoadBalancerSourceRanges: []string{"hello"},
		},
	}
	listReturn := &v1.ServiceList{
		Items: []v1.Service{
			service,
		},
	}
	k8sClient.listFunc = func(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
		listPointer := list.(*v1.ServiceList)
		listReturn.DeepCopyInto(listPointer)
		return nil
	}
	rangerToFilterFor := domain.LbSrcRanger{}
	l, _ := readOps.FilterFor(&rangerToFilterFor)
	assert.Equal(t, 1, k8sClient.listCalled)
	assert.Equal(t, 1, len(l))
	assert.Equal(t, "my-service", l[0].Name)
	assert.Equal(t, []domain.Cidr{domain.Cidr(service.Spec.LoadBalancerSourceRanges[0])}, l[0].LbSrcRanges)
}

func TestLbServicesReadOpsImpl_FilterFor_err(t *testing.T) {
	k8sClient := mockK8sClient{}
	components := lbSrcRangerReconcilerComponents{Client: &k8sClient}
	readOps := MkLbServicesReadOps(&components)
	k8sClient.listFunc = func(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
		return errors.New("crud")
	}
	rangerToFilterFor := domain.LbSrcRanger{}
	_, err := readOps.FilterFor(&rangerToFilterFor)
	assert.Equal(t, 1, k8sClient.listCalled)
	assert.Equal(t, "crud", err.Error())
}

func TestLbServicesUpdateOpsImpl_UpdateCidrs_ok(t *testing.T) {
	k8sClient := mockK8sClient{}
	k8sClient.updateFunc = func(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
		return nil
	}
	components := lbSrcRangerReconcilerComponents{
		Client: &k8sClient,
	}
	updateOps := lbServicesUpdateOpsImpl{
		components: &components,
		k8sService: &v1.Service{},
	}
	err := updateOps.UpdateCidrs(&[]domain.Cidr{"1.1.1.0/25"})
	assert.Equal(t, 1, k8sClient.updateCalled)
	assert.Nil(t, err)
}

func TestLbServicesUpdateOpsImpl_UpdateCidrs_err(t *testing.T) {
	k8sClient := mockK8sClient{}
	k8sClient.updateFunc = func(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
		return errors.New("ugh")
	}
	components := lbSrcRangerReconcilerComponents{
		Client: &k8sClient,
	}
	updateOps := lbServicesUpdateOpsImpl{
		components: &components,
		k8sService: &v1.Service{},
	}
	err := updateOps.UpdateCidrs(&[]domain.Cidr{})
	assert.Equal(t, 1, k8sClient.updateCalled)
	assert.Equal(t, "ugh", err.Error())
}
