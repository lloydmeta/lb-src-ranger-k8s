package k8s

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

// A struct that implements a K8 client for mocking in tests
type mockK8sClient struct {
	sync.Mutex

	// k8s client mocking funcs
	getFunc   func(ctx context.Context, key client.ObjectKey, obj runtime.Object) error
	getCalled int

	listFunc   func(ctx context.Context, list runtime.Object, opts ...client.ListOption) error
	listCalled int

	createFunc   func(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error
	createCalled int

	deleteFunc   func(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error
	deleteCalled int

	updateFunc   func(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error
	updateCalled int

	patchFunc   func(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error
	patchCalled int

	deleteAllOfFunc   func(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error
	deleteAllOfCalled int
}

func (m *mockK8sClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	m.Lock()
	defer m.Unlock()
	m.getCalled += 1
	return m.getFunc(ctx, key, obj)
}

func (m *mockK8sClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	m.Lock()
	defer m.Unlock()
	m.listCalled += 1
	return m.listFunc(ctx, list, opts...)
}

func (m *mockK8sClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	m.Lock()
	defer m.Unlock()
	m.createCalled += 1
	return m.createFunc(ctx, obj, opts...)
}

func (m *mockK8sClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	m.Lock()
	defer m.Unlock()
	m.deleteCalled += 1
	return m.deleteFunc(ctx, obj, opts...)
}

func (m *mockK8sClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	m.Lock()
	defer m.Unlock()
	m.updateCalled += 1
	return m.updateFunc(ctx, obj, opts...)
}

func (m *mockK8sClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	m.Lock()
	defer m.Unlock()
	m.patchCalled += 1
	return m.patchFunc(ctx, obj, patch, opts...)
}

func (m *mockK8sClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	m.Lock()
	defer m.Unlock()
	m.deleteAllOfCalled += 1
	return m.deleteAllOfFunc(ctx, obj, opts...)
}

func (m *mockK8sClient) Status() client.StatusWriter {
	return m
}
