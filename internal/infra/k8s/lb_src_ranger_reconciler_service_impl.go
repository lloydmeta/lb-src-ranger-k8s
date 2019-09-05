package k8s

import (
	"context"
	"github.com/go-logr/logr"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/lloydmeta/lb-src-ranger-k8s/internal/domain"
	"github.com/lloydmeta/lb-src-ranger-k8s/internal/infra"
	"time"
)

func MkReconcilerService(
	id domain.LbSrcRangerId,
	client client.Client,
	log logr.Logger,
	httpClient *http.Client,
	ctx context.Context,
	now func() time.Time) domain.LbSrcRangerReconcilerService {
	r := lbSrcRangerReconcilerComponents{
		Client:     client,
		log:        log,
		httpClient: httpClient,
		ctx:        ctx,
	}
	return &lbSrcReconcilerServiceImpl{
		id:         id,
		components: &r,
		now:        now,
	}

}

type lbSrcReconcilerServiceImpl struct {
	id         domain.LbSrcRangerId
	components *lbSrcRangerReconcilerComponents
	now        func() time.Time
}

func (l *lbSrcReconcilerServiceImpl) MkCidrsFetcher() domain.CidrsFetcher {
	return domain.MkCidrsFetcher(infra.MkHttpClientUrlReader(l.components.httpClient))
}

func (l *lbSrcReconcilerServiceImpl) Now() time.Time {
	return l.now()
}

func (l *lbSrcReconcilerServiceImpl) Id() domain.LbSrcRangerId {
	return l.id
}

func (l *lbSrcReconcilerServiceImpl) MkLogger() logr.Logger {
	return l.components.log.WithValues("loadbalancerranger", l.id)
}

func (l *lbSrcReconcilerServiceImpl) MkLbRangerReadOps() domain.LbSrcRangersReadOps {
	return MkLbSrcRangerReadOps(l.components)
}

func (l *lbSrcReconcilerServiceImpl) MkLbServicesReadOps() domain.LbServicesReadOps {
	return MkLbServicesReadOps(l.components)
}
