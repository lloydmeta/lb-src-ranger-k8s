package controllers

import (
	"context"
	"github.com/lloydmeta/lb-src-ranger-k8s/api/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Context("Inside of a new namespace", func() {
	ctx := context.TODO()
	ns := setupTest(ctx)

	Describe("updating services", func() {

		It("should update existing services with the right labels", func() {
			createService := &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dummy-svc2",
					Namespace: ns.Name,
					Labels:    map[string]string{"src": "ranger"},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port:       80,
							TargetPort: intstr.FromString("http"),
						},
					},
					Type: v1.ServiceTypeLoadBalancer,
				},
			}
			err := k8sClient.Create(ctx, createService)
			Expect(err).NotTo(HaveOccurred(), "failed to create service")

			createRanger := &v1beta1.LbSrcRanger{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "lbsrcranger-test2",
					Namespace: ns.Name,
				},
				Spec: v1beta1.LbSrcRangerSpec{
					TargetLabels: createService.ObjectMeta.Labels,
					UpdateEvery:  metav1.Duration{Duration: time.Millisecond * 100},
					SrcIPUrls:    []string{mockSrcServerUrl},
				},
			}
			err = k8sClient.Create(ctx, createRanger)
			Expect(err).NotTo(HaveOccurred(), "failed to create test LbSrcRanger resource")

			retrievedRanger := &v1beta1.LbSrcRanger{}
			Eventually(
				getResourceFunc(ctx, client.ObjectKey{Name: createRanger.Name, Namespace: createRanger.Namespace}, retrievedRanger),
				time.Second*5, time.Millisecond*100).Should(BeNil())

			retrievedService := &v1.Service{}
			expected := buildCidrStrs(&mockSrcCidrs)
			Eventually(func() []string {
				if err := getResourceFunc(ctx, client.ObjectKey{Name: createService.Name, Namespace: createService.Namespace}, retrievedService)(); err != nil {
					return nil
				} else {
					return retrievedService.Spec.LoadBalancerSourceRanges
				}
			}, time.Second*10, time.Millisecond*100).Should(Equal(expected))
		})
	})
})

func getResourceFunc(ctx context.Context, key client.ObjectKey, obj runtime.Object) func() error {
	return func() error {
		return k8sClient.Get(ctx, key, obj)
	}
}
