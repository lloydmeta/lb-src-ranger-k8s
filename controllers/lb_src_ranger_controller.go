/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/lloydmeta/lb-src-ranger-k8s/internal/domain"

	"github.com/lloydmeta/lb-src-ranger-k8s/internal/infra/k8s"

	"k8s.io/utils/clock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"net/http"

	lbsrcrangerv1beta1 "github.com/lloydmeta/lb-src-ranger-k8s/api/v1beta1"
)

// LbSrcRangerController holds dependencies that will be used for reconciling
// a LbSrcRanger object. Honestly, just here to fix the compiler's circular dependencies
// error
type LbSrcRangerController struct {
	client     client.Client
	logger     logr.Logger
	clock      clock.Clock
	httpClient *http.Client
}

func MkLbSrcRangerController(
	client client.Client,
	logger logr.Logger,
	clock clock.Clock,
	httpClient *http.Client) LbSrcRangerController {
	return LbSrcRangerController{
		client: client, logger: logger, clock: clock, httpClient: httpClient,
	}
}

// +kubebuilder:rbac:groups=lbsrcranger.beachape.com,resources=lbsrcrangers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lbsrcranger.beachape.com,resources=lbsrcrangers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;update;patch
func (r *LbSrcRangerController) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	id := domain.LbSrcRangerId(req.NamespacedName)
	lbRangerReconcilerService := k8s.MkReconcilerService(id, r.client, r.logger, r.httpClient, ctx, r.clock.Now)
	result, err := domain.Reconcile(lbRangerReconcilerService)
	return ctrl.Result(result), err
}

func (r *LbSrcRangerController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lbsrcrangerv1beta1.LbSrcRanger{}).
		Complete(r)
}
