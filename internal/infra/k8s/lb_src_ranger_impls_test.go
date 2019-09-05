package k8s

import (
	"context"
	"github.com/lloydmeta/lb-src-ranger-k8s/api/v1beta1"
	"github.com/lloydmeta/lb-src-ranger-k8s/internal/domain"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestLbRangersReadOpsImpl_Get_ok(t *testing.T) {
	k8sClient := mockK8sClient{}
	components := lbSrcRangerReconcilerComponents{Client: &k8sClient}
	readOps := MkLbSrcRangerReadOps(&components)
	ranger := v1beta1.LbSrcRanger{
		Spec: v1beta1.LbSrcRangerSpec{
			SrcIPUrls: []string{"http://somewhereoutthere/ips"},
		},
	}
	k8sClient.getFunc = func(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
		p := obj.(*v1beta1.LbSrcRanger)
		ranger.DeepCopyInto(p)
		return nil
	}
	r, err := readOps.Get(&domain.LbSrcRangerId{})
	assert.Nil(t, err)
	assert.Equal(t, 1, k8sClient.getCalled)
	assert.Equal(t, []string{"http://somewhereoutthere/ips"}, r.Spec.SrcIPUrls)
}

func TestLbRangersReadOpsImpl_Get_err_not_found(t *testing.T) {
	k8sClient := mockK8sClient{}
	components := lbSrcRangerReconcilerComponents{Client: &k8sClient}
	readOps := MkLbSrcRangerReadOps(&components)

	k8sClient.getFunc = func(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
		return &errors.StatusError{ErrStatus: v1.Status{
			Reason: v1.StatusReasonNotFound,
		}}
	}
	_, err := readOps.Get(&domain.LbSrcRangerId{})
	assert.Equal(t, 1, k8sClient.getCalled)
	assert.True(t, err.IsNotFound)
}

func TestLbRangersReadOpsImpl_Get_err_other(t *testing.T) {
	k8sClient := mockK8sClient{}
	components := lbSrcRangerReconcilerComponents{Client: &k8sClient}
	readOps := MkLbSrcRangerReadOps(&components)

	k8sClient.getFunc = func(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
		return errors.NewBadRequest("lol")
	}
	_, err := readOps.Get(&domain.LbSrcRangerId{})
	assert.Equal(t, 1, k8sClient.getCalled)
	assert.False(t, err.IsNotFound)
}

func TestLbRangersUpdateOpsImpl_UpdateStatus_ok(t *testing.T) {
	k8sClient := mockK8sClient{}
	components := lbSrcRangerReconcilerComponents{Client: &k8sClient}
	lbRanger := v1beta1.LbSrcRanger{}
	updateOps := lbSrcRangersUpdateOpsImpl{components: &components, k8sLoadBalancerRanger: &lbRanger}
	k8sClient.updateFunc = func(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
		return nil
	}
	status := domain.LbSrcRangerStatus{}
	err := updateOps.UpdateStatus(&status)
	assert.Equal(t, 1, k8sClient.updateCalled)
	assert.Nil(t, err)
}

func TestLbRangersUpdateOpsImpl_UpdateStatus_err(t *testing.T) {
	k8sClient := mockK8sClient{}
	components := lbSrcRangerReconcilerComponents{Client: &k8sClient}
	lbRanger := v1beta1.LbSrcRanger{}
	updateOps := lbSrcRangersUpdateOpsImpl{components: &components, k8sLoadBalancerRanger: &lbRanger}
	k8sClient.updateFunc = func(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
		return errors.NewBadRequest("ugh")
	}
	status := domain.LbSrcRangerStatus{}
	err := updateOps.UpdateStatus(&status)
	assert.Equal(t, 1, k8sClient.updateCalled)
	assert.Equal(t, "ugh", err.Error())
}
