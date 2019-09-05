package k8s

import (
	"context"
	"github.com/go-logr/logr"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"net/http"
)

// This holds k8-related components as well as base "infra" std-lib Go
// IO-related components like an http client
type lbSrcRangerReconcilerComponents struct {
	client.Client
	log        logr.Logger
	httpClient *http.Client
	ctx        context.Context
}
