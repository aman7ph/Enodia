package firewall

import (
	"log"
	"runtime"
	"sync"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// Manager handles the COM thread and exposes thread-safe methods
type Manager struct {
	jobs    chan func(*ole.IDispatch)
	stop    chan struct{}
	wg      sync.WaitGroup
	running bool
}

// NewManager initializes the background COM worker
func NewManager() *Manager {
	m := &Manager{
		jobs: make(chan func(*ole.IDispatch)),
		stop: make(chan struct{}),
	}
	m.running = true
	m.wg.Add(1)
	go m.worker()
	return m
}

// Close cleans up the worker thread
func (m *Manager) Close() {
	if m.running {
		close(m.stop)
		m.wg.Wait()
		m.running = false
	}
}

// worker is the dedicated OS thread for all COM interaction
func (m *Manager) worker() {
	defer m.wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := ole.CoInitialize(0); err != nil {
		log.Printf("[Enodia] FATAL: COM Init failed: %v", err)
		return
	}
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject(PROGID_POLICY2)
	if err != nil {
		log.Printf("[Enodia] FATAL: Failed to create Policy2: %v", err)
		return
	}
	defer unknown.Release()

	policy, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Printf("[Enodia] FATAL: Failed to query Policy: %v", err)
		return
	}
	defer policy.Release()

	rulesRaw, err := oleutil.GetProperty(policy, "Rules")
	if err != nil {
		log.Printf("[Enodia] FATAL: Failed to get Rules: %v", err)
		return
	}
	rules := rulesRaw.ToIDispatch()
	defer rules.Release()

	log.Println("[Enodia] Firewall Manager Ready.")

	for {
		select {
		case job := <-m.jobs:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[Enodia] Panic in job: %v", r)
					}
				}()
				job(rules)
			}()
		case <-m.stop:
			return
		}
	}
}
