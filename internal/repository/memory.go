package repository

import (
	"strings"
	"sync"

	"greenpark/marketing/internal/domain"
)

// fileRepository is a mutex-guarded MarketingRepository. The full state is held
// in memory for fast reads and flushed to its persister (file or DB) on every
// write.
type fileRepository struct {
	mu sync.RWMutex
	p  persister
	st *state
}

// NewRepository returns a MarketingRepository persisted to the given JSON file
// path. An empty path keeps everything in memory only (handy for tests).
func NewRepository(path string) (MarketingRepository, error) {
	return newRepository(filePersister{path: path})
}

// NewPostgresRepository returns a MarketingRepository that persists the
// whole-state snapshot to a single PostgreSQL row.
func NewPostgresRepository(dsn string) (MarketingRepository, error) {
	p, err := newPGPersister(dsn)
	if err != nil {
		return nil, err
	}
	return newRepository(p)
}

func newRepository(p persister) (MarketingRepository, error) {
	st, err := p.load()
	if err != nil {
		return nil, err
	}
	return &fileRepository{p: p, st: st}, nil
}

// persist flushes the current state. Callers must hold the write lock.
func (r *fileRepository) persist() error { return r.p.save(r.st) }

/* ---------------------------- reads ---------------------------- */

func (r *fileRepository) Context() domain.Context {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.st.Context
}

func (r *fileRepository) Funnel() []domain.FunnelStage {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.Funnel)
}

func (r *fileRepository) Spend() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.st.Spend
}

func (r *fileRepository) KPIs() []domain.KPI {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.KPIs)
}

func (r *fileRepository) LeadQuality() domain.LeadQuality {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.st.LeadQuality
}

func (r *fileRepository) Handover() []domain.HandoverItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.Handover)
}

func (r *fileRepository) Channels() []domain.Channel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.Channels)
}

func (r *fileRepository) Projects() []domain.Project {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.Projects)
}

func (r *fileRepository) ProjectByName(name string) (domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.st.Projects {
		if strings.EqualFold(p.Name, name) {
			return p, nil
		}
	}
	return domain.Project{}, ErrNotFound
}

func (r *fileRepository) Assets() []domain.Asset {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.Assets)
}

func (r *fileRepository) IGAccounts() []domain.IGAccount {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.IGAccounts)
}

func (r *fileRepository) Content() domain.Content {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.st.Content
}

func (r *fileRepository) Commands() []domain.Command {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.Commands)
}

func (r *fileRepository) Alerts() domain.Alerts {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.st.Alerts
}

func (r *fileRepository) ReasonCodes() []domain.ReasonCode {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return clone(r.st.ReasonCodes)
}

/* ---------------------------- singleton / whole-value writes ---------------------------- */

func (r *fileRepository) SetContext(c domain.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.st.Context = c
	return r.persist()
}

func (r *fileRepository) SetSpend(v float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.st.Spend = v
	return r.persist()
}

func (r *fileRepository) SetLeadQuality(lq domain.LeadQuality) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.st.LeadQuality = lq
	return r.persist()
}

func (r *fileRepository) SetContent(c domain.Content) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.st.Content = c
	return r.persist()
}

func (r *fileRepository) SetAlerts(a domain.Alerts) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.st.Alerts = a
	return r.persist()
}

/* ---------------------------- collection writes ---------------------------- */

func (r *fileRepository) SaveFunnelStage(v domain.FunnelStage) (domain.FunnelStage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("fnl")
	}
	r.st.Funnel = upsertEntity(r.st.Funnel, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteFunnelStage(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.Funnel, id)
	r.st.Funnel = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

func (r *fileRepository) SaveKPI(v domain.KPI) (domain.KPI, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("kpi")
	}
	r.st.KPIs = upsertEntity(r.st.KPIs, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteKPI(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.KPIs, id)
	r.st.KPIs = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

func (r *fileRepository) SaveHandoverItem(v domain.HandoverItem) (domain.HandoverItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("ho")
	}
	r.st.Handover = upsertEntity(r.st.Handover, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteHandoverItem(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.Handover, id)
	r.st.Handover = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

func (r *fileRepository) SaveChannel(v domain.Channel) (domain.Channel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("ch")
	}
	r.st.Channels = upsertEntity(r.st.Channels, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteChannel(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.Channels, id)
	r.st.Channels = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

func (r *fileRepository) SaveProject(v domain.Project) (domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("prj")
	}
	r.st.Projects = upsertEntity(r.st.Projects, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteProject(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.Projects, id)
	r.st.Projects = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

func (r *fileRepository) SaveAsset(v domain.Asset) (domain.Asset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("ast")
	}
	r.st.Assets = upsertEntity(r.st.Assets, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteAsset(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.Assets, id)
	r.st.Assets = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

func (r *fileRepository) SaveIGAccount(v domain.IGAccount) (domain.IGAccount, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("ig")
	}
	r.st.IGAccounts = upsertEntity(r.st.IGAccounts, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteIGAccount(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.IGAccounts, id)
	r.st.IGAccounts = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

func (r *fileRepository) SaveCommand(v domain.Command) (domain.Command, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("cmd")
	}
	r.st.Commands = upsertEntity(r.st.Commands, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteCommand(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.Commands, id)
	r.st.Commands = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

func (r *fileRepository) SaveReasonCode(v domain.ReasonCode) (domain.ReasonCode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.EntID == "" {
		v.EntID = newID("rc")
	}
	r.st.ReasonCodes = upsertEntity(r.st.ReasonCodes, v)
	return v, r.persist()
}

func (r *fileRepository) DeleteReasonCode(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next, ok := deleteEntity(r.st.ReasonCodes, id)
	r.st.ReasonCodes = next
	if !ok {
		return false, nil
	}
	return true, r.persist()
}

/* ---------------------------- users ---------------------------- */

func (r *fileRepository) Users() []domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.User, len(r.st.Users))
	for i, u := range r.st.Users {
		out[i] = u.toDomain()
	}
	return out
}

func (r *fileRepository) UserByUsername(username string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	username = strings.ToLower(strings.TrimSpace(username))
	for _, u := range r.st.Users {
		if strings.ToLower(u.Username) == username {
			return u.toDomain(), nil
		}
	}
	return domain.User{}, ErrNotFound
}

func (r *fileRepository) UserByID(id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.st.Users {
		if u.ID == id {
			return u.toDomain(), nil
		}
	}
	return domain.User{}, ErrNotFound
}
