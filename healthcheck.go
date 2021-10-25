package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/Financial-Times/gourmet/log"
)

type HealthChecker struct {
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

func (hc *HealthChecker) check(addr string) (bool, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return false, fmt.Errorf("could not parse url %s: %w", addr, err)
	}
	u.Path = path.Join(u.Path, "/__health")
	resp, err := http.Get(u.String())
	if err != nil {
		return false, fmt.Errorf("could not make http request to check healthiness: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("service returned unexpected http status code %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("could not read service response: %w", err)
	}
	type healthStatus struct {
		OK bool `json:"ok"`
	}
	hs := &healthStatus{}
	err = json.Unmarshal(body, hs)
	if err != nil {
		return false, fmt.Errorf("could not parse service response: %w", err)
	}

	if !hs.OK {
		return false, nil
	}

	return true, nil
}

type HealthCheckStatus struct {
	Service   Service `json:"service"`
	IsHealthy bool    `json:"is_healthy"`
	Err       error   `json:"error"`
}

type HealthStore struct {
	data map[string]HealthCheckStatus
	mtx  sync.Mutex
}

func NewHealthStore() *HealthStore {
	return &HealthStore{
		data: make(map[string]HealthCheckStatus),
	}
}

func (hs *HealthStore) Set(s Service, isHealthy bool, err error) {
	hs.mtx.Lock()
	defer hs.mtx.Unlock()
	hs.data[s.GetID()] = HealthCheckStatus{
		Service:   s,
		IsHealthy: isHealthy,
		Err:       err,
	}
}

func (hs *HealthStore) Get(s Service) HealthCheckStatus {
	hs.mtx.Lock()
	defer hs.mtx.Unlock()
	return hs.data[s.GetID()]
}

func (hs *HealthStore) GetAll() []HealthCheckStatus {
	hs.mtx.Lock()
	defer hs.mtx.Unlock()
	var retval []HealthCheckStatus
	for _, i := range hs.data {
		retval = append(retval, i)
	}
	return retval
}

type HealthMonitor struct {
	checker     *HealthChecker
	healthStore *HealthStore
	services    []Service
	logger      *log.StructuredLogger
}

func (m *HealthMonitor) Check() {
	m.logger.Debug("tick")
	for _, s := range m.services {
		isHealthy, err := m.checker.check(s.URL)
		if err != nil {
			m.logger.Info("error getting health status",
				log.WithField("url", s.URL),
				log.WithField("name", s.Name),
				log.WithError(err),
			)
		}
		m.healthStore.Set(s, isHealthy, err)
	}
}

type HealthCheckLifecycle struct {
	monitor         *HealthMonitor
	refreshInterval int
	doneChan        chan bool
}

func NewHealthCheckLifescycle(m *HealthMonitor, r int) *HealthCheckLifecycle {
	return &HealthCheckLifecycle{
		monitor:         m,
		refreshInterval: r,
		doneChan:        make(chan bool),
	}
}

func (l *HealthCheckLifecycle) OnStart(_ context.Context) error {
	ticker := time.NewTicker(time.Duration(l.refreshInterval) * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				l.monitor.Check()
			case <-l.doneChan:
				return
			}
		}
	}()
	l.monitor.Check()
	return nil
}

func (l *HealthCheckLifecycle) OnStop(_ context.Context) error {
	l.doneChan <- true
	return nil
}
