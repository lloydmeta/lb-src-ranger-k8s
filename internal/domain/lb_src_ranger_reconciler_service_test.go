package domain

import (
	"errors"
	"github.com/go-logr/logr"
	testing2 "github.com/go-logr/logr/testing"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

var (
	mockRangerId = LbSrcRangerId{}
)

func Test_buildRequeueAfter_ok(t *testing.T) {
	r := buildRequeueAfter(time.Minute, testing2.NullLogger{})
	assert.Equal(t, time.Minute, r)
}

func Test_buildRequeueAfter_too_small(t *testing.T) {
	r := buildRequeueAfter(time.Duration(0), testing2.NullLogger{})
	assert.Equal(t, time.Minute, r)
}

func TestReconcile_ok(t *testing.T) {
	rangerUpdateOps := mockLbRangersUpdateOps{updateFunc: func(status *LbSrcRangerStatus) error {
		return nil
	}}

	rangerReaderOps := mockLbRangersReadOps{getFunc: func(id *LbSrcRangerId) (ranger LbSrcRanger, err *LbSrcRangersReadGetErr) {
		return LbSrcRanger{
			Spec: LbSrcRangerSpec{
				UpdateEvery: time.Minute,
			},
			UpdateOps: &rangerUpdateOps,
		}, nil
	}}

	servicesUpdateOps := mockLbServicesUpdateOps{updateFunc: func(cidrs *[]Cidr) error {
		return nil
	}}
	servicesReadOps := mockLbServicesReadOps{filterFunc: func(l *LbSrcRanger) (services []LbService, e error) {
		return []LbService{
			{
				UpdateOps: &servicesUpdateOps,
			},
		}, nil
	}}
	cidrsFetcher := mockCidrsFetcher{
		cidrs: []Cidr{"1.2.3.0/25"},
	}
	service := mkMockService(mockRangerId, &rangerReaderOps, &servicesReadOps, &cidrsFetcher, time.Now())
	r, err := Reconcile(&service)
	assert.Nil(t, err)
	assert.Equal(t, 1, rangerReaderOps.getCalled)
	assert.Equal(t, 1, rangerUpdateOps.updateCalled)
	assert.Equal(t, 1, servicesReadOps.filterCalled)
	assert.Equal(t, 1, servicesUpdateOps.updateCalled)
	assert.Equal(t, 1, cidrsFetcher.called)
	assert.Equal(t, r.RequeueAfter, time.Minute)
}

func TestReconcile_ok_no_services(t *testing.T) {
	rangerUpdateOps := mockLbRangersUpdateOps{updateFunc: func(status *LbSrcRangerStatus) error {
		return nil
	}}

	rangerReaderOps := mockLbRangersReadOps{getFunc: func(id *LbSrcRangerId) (ranger LbSrcRanger, err *LbSrcRangersReadGetErr) {
		return LbSrcRanger{
			Spec: LbSrcRangerSpec{
				UpdateEvery: time.Minute,
			},
			UpdateOps: &rangerUpdateOps,
		}, nil
	}}

	servicesUpdateOps := mockLbServicesUpdateOps{updateFunc: func(cidrs *[]Cidr) error {
		return nil
	}}
	servicesReadOps := mockLbServicesReadOps{filterFunc: func(l *LbSrcRanger) (services []LbService, e error) {
		return []LbService{}, nil
	}}
	cidrsFetcher := mockCidrsFetcher{}
	service := mkMockService(mockRangerId, &rangerReaderOps, &servicesReadOps, &cidrsFetcher, time.Now())
	r, err := Reconcile(&service)
	assert.Nil(t, err)
	assert.Equal(t, 1, rangerReaderOps.getCalled)
	assert.Equal(t, 0, rangerUpdateOps.updateCalled)
	assert.Equal(t, 1, servicesReadOps.filterCalled)
	assert.Equal(t, 0, servicesUpdateOps.updateCalled)
	assert.Equal(t, 0, cidrsFetcher.called)
	assert.Equal(t, r.RequeueAfter, time.Minute)
}

func TestReconcile_rangers_reader_err_not_found(t *testing.T) {
	rangerUpdateOps := mockLbRangersUpdateOps{}
	rangerReaderOps := mockLbRangersReadOps{getFunc: func(id *LbSrcRangerId) (ranger LbSrcRanger, err *LbSrcRangersReadGetErr) {
		return LbSrcRanger{}, &LbSrcRangersReadGetErr{IsNotFound: true}
	}}

	servicesUpdateOps := mockLbServicesUpdateOps{}
	servicesReadOps := mockLbServicesReadOps{}
	cidrsFetcher := mockCidrsFetcher{}
	service := mkMockService(mockRangerId, &rangerReaderOps, &servicesReadOps, &cidrsFetcher, time.Now())
	r, err := Reconcile(&service)
	assert.Nil(t, err)
	assert.Equal(t, 1, rangerReaderOps.getCalled)
	assert.Equal(t, 0, rangerUpdateOps.updateCalled)
	assert.Equal(t, 0, servicesReadOps.filterCalled)
	assert.Equal(t, 0, servicesUpdateOps.updateCalled)
	assert.Equal(t, 0, cidrsFetcher.called)
	assert.False(t, r.Requeue)
}

func TestReconcile_rangers_reader_err(t *testing.T) {
	rangerUpdateOps := mockLbRangersUpdateOps{}
	rangerReaderOps := mockLbRangersReadOps{getFunc: func(id *LbSrcRangerId) (ranger LbSrcRanger, err *LbSrcRangersReadGetErr) {
		return LbSrcRanger{}, &LbSrcRangersReadGetErr{IsNotFound: false}
	}}

	servicesUpdateOps := mockLbServicesUpdateOps{}
	servicesReadOps := mockLbServicesReadOps{}
	cidrsFetcher := mockCidrsFetcher{}
	service := mkMockService(mockRangerId, &rangerReaderOps, &servicesReadOps, &cidrsFetcher, time.Now())
	r, err := Reconcile(&service)
	assert.NotNil(t, err)
	assert.Equal(t, 1, rangerReaderOps.getCalled)
	assert.Equal(t, 0, rangerUpdateOps.updateCalled)
	assert.Equal(t, 0, servicesReadOps.filterCalled)
	assert.Equal(t, 0, servicesUpdateOps.updateCalled)
	assert.Equal(t, 0, cidrsFetcher.called)
	assert.True(t, r.Requeue)
}

func TestReconcile_services_filter_err(t *testing.T) {
	rangerUpdateOps := mockLbRangersUpdateOps{updateFunc: func(status *LbSrcRangerStatus) error {
		return nil
	}}

	rangerReaderOps := mockLbRangersReadOps{getFunc: func(id *LbSrcRangerId) (ranger LbSrcRanger, err *LbSrcRangersReadGetErr) {
		return LbSrcRanger{
			Spec: LbSrcRangerSpec{
				UpdateEvery: time.Minute,
			},
			UpdateOps: &rangerUpdateOps,
		}, nil
	}}

	servicesReadOps := mockLbServicesReadOps{filterFunc: func(l *LbSrcRanger) (services []LbService, e error) {
		return nil, errors.New("ugh")
	}}
	cidrsFetcher := mockCidrsFetcher{}
	service := mkMockService(mockRangerId, &rangerReaderOps, &servicesReadOps, &cidrsFetcher, time.Now())
	r, err := Reconcile(&service)
	assert.Equal(t, "ugh", err.Error())
	assert.Equal(t, 1, rangerReaderOps.getCalled)
	assert.Equal(t, 0, rangerUpdateOps.updateCalled)
	assert.Equal(t, 1, servicesReadOps.filterCalled)
	assert.Equal(t, 0, cidrsFetcher.called)
	assert.True(t, r.Requeue)
}

func TestReconcile_cidr_fetch_err(t *testing.T) {
	rangerUpdateOps := mockLbRangersUpdateOps{updateFunc: func(status *LbSrcRangerStatus) error {
		return nil
	}}
	rangerReaderOps := mockLbRangersReadOps{getFunc: func(id *LbSrcRangerId) (ranger LbSrcRanger, err *LbSrcRangersReadGetErr) {
		return LbSrcRanger{
			Spec: LbSrcRangerSpec{
				UpdateEvery: time.Minute,
			},
			UpdateOps: &rangerUpdateOps,
		}, nil
	}}
	servicesUpdateOps := mockLbServicesUpdateOps{updateFunc: func(cidrs *[]Cidr) error {
		return nil
	}}
	servicesReadOps := mockLbServicesReadOps{filterFunc: func(l *LbSrcRanger) (services []LbService, e error) {
		return []LbService{
			{
				UpdateOps: &servicesUpdateOps,
			},
		}, nil
	}}
	cidrsFetcher := mockCidrsFetcher{
		err: true,
	}
	service := mkMockService(mockRangerId, &rangerReaderOps, &servicesReadOps, &cidrsFetcher, time.Now())
	r, err := Reconcile(&service)
	assert.NotNil(t, err)
	assert.Equal(t, 1, rangerReaderOps.getCalled)
	assert.Equal(t, 0, rangerUpdateOps.updateCalled)
	assert.Equal(t, 1, servicesReadOps.filterCalled)
	assert.Equal(t, 1, cidrsFetcher.called)
	assert.True(t, r.Requeue)
}

func TestReconcile_services_update_err(t *testing.T) {
	rangerUpdateOps := mockLbRangersUpdateOps{updateFunc: func(status *LbSrcRangerStatus) error {
		return nil
	}}
	rangerReaderOps := mockLbRangersReadOps{getFunc: func(id *LbSrcRangerId) (ranger LbSrcRanger, err *LbSrcRangersReadGetErr) {
		return LbSrcRanger{
			Spec: LbSrcRangerSpec{
				UpdateEvery: time.Minute,
			},
			UpdateOps: &rangerUpdateOps,
		}, nil
	}}
	servicesUpdateOps := mockLbServicesUpdateOps{updateFunc: func(cidrs *[]Cidr) error {
		return errors.New("ugh")
	}}
	servicesReadOps := mockLbServicesReadOps{filterFunc: func(l *LbSrcRanger) (services []LbService, e error) {
		return []LbService{
			{
				UpdateOps: &servicesUpdateOps,
			},
		}, nil
	}}
	cidrsFetcher := mockCidrsFetcher{
		cidrs: []Cidr{"1.1.1.0/25"},
	}
	service := mkMockService(mockRangerId, &rangerReaderOps, &servicesReadOps, &cidrsFetcher, time.Now())
	r, err := Reconcile(&service)
	assert.NotNil(t, err)
	assert.Equal(t, 1, rangerReaderOps.getCalled)
	assert.Equal(t, 0, rangerUpdateOps.updateCalled)
	assert.Equal(t, 1, servicesReadOps.filterCalled)
	assert.Equal(t, 1, cidrsFetcher.called)
	assert.True(t, r.Requeue)
}

func TestReconcile_rangers_update_err(t *testing.T) {
	rangerUpdateOps := mockLbRangersUpdateOps{updateFunc: func(status *LbSrcRangerStatus) error {
		return errors.New("ugh")
	}}
	rangerReaderOps := mockLbRangersReadOps{getFunc: func(id *LbSrcRangerId) (ranger LbSrcRanger, err *LbSrcRangersReadGetErr) {
		return LbSrcRanger{
			Spec: LbSrcRangerSpec{
				UpdateEvery: time.Minute,
			},
			UpdateOps: &rangerUpdateOps,
		}, nil
	}}
	servicesUpdateOps := mockLbServicesUpdateOps{updateFunc: func(cidrs *[]Cidr) error {
		return nil
	}}
	servicesReadOps := mockLbServicesReadOps{filterFunc: func(l *LbSrcRanger) (services []LbService, e error) {
		return []LbService{
			{
				UpdateOps: &servicesUpdateOps,
			},
		}, nil
	}}
	cidrsFetcher := mockCidrsFetcher{
		cidrs: []Cidr{"1.1.1.0/25"},
	}
	service := mkMockService(mockRangerId, &rangerReaderOps, &servicesReadOps, &cidrsFetcher, time.Now())
	r, err := Reconcile(&service)
	assert.NotNil(t, err)
	assert.Equal(t, 1, rangerReaderOps.getCalled)
	assert.Equal(t, 1, rangerUpdateOps.updateCalled)
	assert.Equal(t, 1, servicesReadOps.filterCalled)
	assert.Equal(t, 1, cidrsFetcher.called)
	assert.True(t, r.Requeue)
}

// mocks

func mkMockService(
	id LbSrcRangerId,
	lbRangersReadOps *mockLbRangersReadOps,
	lbServicesReadOps *mockLbServicesReadOps,
	cidrsFetcher *mockCidrsFetcher,
	now time.Time) mockServiceImpl {
	return mockServiceImpl{
		id:                id,
		lbRangersReadOps:  lbRangersReadOps,
		lbServicesReadOps: lbServicesReadOps,
		cidrsFetcher:      cidrsFetcher,
		now:               now,
	}
}

type mockServiceImpl struct {
	id                LbSrcRangerId
	lbRangersReadOps  *mockLbRangersReadOps
	lbServicesReadOps *mockLbServicesReadOps
	cidrsFetcher      *mockCidrsFetcher
	now               time.Time
}

func (m *mockServiceImpl) Id() LbSrcRangerId {
	return m.id
}

func (m *mockServiceImpl) MkLbRangerReadOps() LbSrcRangersReadOps {
	return m.lbRangersReadOps
}

func (m *mockServiceImpl) MkLbServicesReadOps() LbServicesReadOps {
	return m.lbServicesReadOps
}

func (m *mockServiceImpl) MkLogger() logr.Logger {
	return testing2.NullLogger{}
}

func (m *mockServiceImpl) MkCidrsFetcher() CidrsFetcher {
	return m.cidrsFetcher
}

func (m *mockServiceImpl) Now() time.Time {
	return m.now
}

type mockLbRangersReadOps struct {
	getFunc   func(id *LbSrcRangerId) (LbSrcRanger, *LbSrcRangersReadGetErr)
	getCalled int
}

func (m *mockLbRangersReadOps) Get(id *LbSrcRangerId) (LbSrcRanger, *LbSrcRangersReadGetErr) {
	m.getCalled += 1
	return m.getFunc(id)
}

type mockLbRangersUpdateOps struct {
	updateFunc   func(status *LbSrcRangerStatus) error
	updateCalled int
}

func (m *mockLbRangersUpdateOps) UpdateStatus(status *LbSrcRangerStatus) error {
	m.updateCalled += 1
	return m.updateFunc(status)
}

type mockLbServicesReadOps struct {
	filterFunc   func(l *LbSrcRanger) ([]LbService, error)
	filterCalled int
}

func (m *mockLbServicesReadOps) FilterFor(l *LbSrcRanger) ([]LbService, error) {
	m.filterCalled += 1
	return m.filterFunc(l)
}

type mockLbServicesUpdateOps struct {
	updateFunc   func(cidrs *[]Cidr) error
	updateCalled int
}

func (m *mockLbServicesUpdateOps) UpdateCidrs(cidrs *[]Cidr) error {
	m.updateCalled += 1
	return m.updateFunc(cidrs)
}

type mockCidrsFetcher struct {
	cidrs  []Cidr
	called int
	err    bool
}

func (m *mockCidrsFetcher) Fetch(srcUrls *[]string) ([]Cidr, error) {
	m.called += 1
	if m.err {
		return nil, errors.New("shiz")
	} else {
		return m.cidrs, nil
	}
}
