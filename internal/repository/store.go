package repository

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"greenpark/marketing/internal/domain"
)

// state is the full persisted snapshot of the dashboard data plus the user
// accounts. It is the on-disk JSON shape; the derived Summary is never stored.
type state struct {
	Context     domain.Context        `json:"context"`
	Funnel      []domain.FunnelStage  `json:"funnel"`
	Spend       float64               `json:"spend"`
	KPIs        []domain.KPI          `json:"kpis"`
	LeadQuality domain.LeadQuality    `json:"leadQuality"`
	Handover    []domain.HandoverItem `json:"handover"`
	Channels    []domain.Channel      `json:"channels"`
	Projects    []domain.Project      `json:"projects"`
	Assets      []domain.Asset        `json:"assets"`
	IGAccounts  []domain.IGAccount    `json:"igAccounts"`
	Content     domain.Content        `json:"content"`
	Commands    []domain.Command      `json:"commands"`
	Alerts      domain.Alerts         `json:"alerts"`
	ReasonCodes []domain.ReasonCode   `json:"reasonCodes"`
	Users       []storeUser           `json:"users"`
}

// storeUser is the persisted user shape. Unlike domain.User (which hides
// password material from API responses via json:"-"), this type DOES serialise
// the salt and hash so accounts survive a restart. It never leaves the store.
type storeUser struct {
	ID           string      `json:"id"`
	Username     string      `json:"username"`
	Name         string      `json:"name"`
	Role         domain.Role `json:"role"`
	PasswordHash string      `json:"passwordHash"`
	Salt         string      `json:"salt"`
}

func (u storeUser) toDomain() domain.User {
	return domain.User{
		ID:           u.ID,
		Username:     u.Username,
		Name:         u.Name,
		Role:         u.Role,
		PasswordHash: u.PasswordHash,
		Salt:         u.Salt,
	}
}

// seedState builds the default data set (seed figures + admin/viewer accounts)
// and assigns a stable synthetic _id to every collection row.
func seedState() *state {
	s := &state{
		Context:     seedContext(),
		Funnel:      seedFunnel(),
		Spend:       seedSpend(),
		KPIs:        seedKPIs(),
		LeadQuality: seedLeadQuality(),
		Handover:    seedHandover(),
		Channels:    seedChannels(),
		Projects:    seedProjects(),
		Assets:      seedAssets(),
		IGAccounts:  seedIGAccounts(),
		Content:     seedContent(),
		Commands:    seedCommands(),
		Alerts:      seedAlerts(),
		ReasonCodes: seedReasonCodes(),
		Users:       seedUsers(),
	}
	assignSeedIDs(s)
	return s
}

// assignSeedIDs gives each seeded row a readable, deterministic _id derived from
// its natural key or position.
func assignSeedIDs(s *state) {
	for i := range s.Funnel {
		s.Funnel[i].EntID = "fnl-" + strings.ToLower(strings.ReplaceAll(s.Funnel[i].Key, " ", "-"))
	}
	for i := range s.KPIs {
		s.KPIs[i].EntID = "kpi-" + s.KPIs[i].ID
	}
	for i := range s.Handover {
		s.Handover[i].EntID = fmt.Sprintf("ho-%02d", i+1)
	}
	for i := range s.Channels {
		s.Channels[i].EntID = fmt.Sprintf("ch-%02d", i+1)
	}
	for i := range s.Projects {
		s.Projects[i].EntID = "prj-" + s.Projects[i].IG
	}
	for i := range s.Assets {
		s.Assets[i].EntID = fmt.Sprintf("ast-%02d", i+1)
	}
	for i := range s.IGAccounts {
		s.IGAccounts[i].EntID = "ig-" + s.IGAccounts[i].Handle
	}
	for i := range s.Commands {
		s.Commands[i].EntID = fmt.Sprintf("cmd-%02d", i+1)
	}
	for i := range s.ReasonCodes {
		s.ReasonCodes[i].EntID = "rc-" + strings.ToLower(s.ReasonCodes[i].Code)
	}
}

// load reads the state from disk; if the file is missing it seeds a fresh state
// and writes it. An empty path means in-memory only (used by tests).
func load(path string) (*state, error) {
	if path == "" {
		return seedState(), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			s := seedState()
			if err := save(path, s); err != nil {
				return nil, err
			}
			return s, nil
		}
		return nil, err
	}
	s := &state{}
	if err := json.Unmarshal(b, s); err != nil {
		return nil, err
	}
	return s, nil
}

// save atomically writes the state to disk (write temp + rename).
func save(path string, s *state) error {
	if path == "" {
		return nil
	}
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// persister abstracts where the whole-state snapshot lives (JSON file or DB row).
// The repository logic is identical regardless of backend.
type persister interface {
	load() (*state, error)
	save(*state) error
}

// filePersister stores the state as an atomic JSON file on disk.
type filePersister struct{ path string }

func (f filePersister) load() (*state, error) { return load(f.path) }
func (f filePersister) save(st *state) error  { return save(f.path, st) }

// newID returns a short, collision-resistant identifier with the given prefix.
func newID(prefix string) string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return prefix + "-0000000000"
	}
	return prefix + "-" + hex.EncodeToString(b)
}

/* ---------- generic CRUD helpers over []T where T implements domain.Entity ---------- */

// upsertEntity replaces the element whose id matches, otherwise appends it.
func upsertEntity[T domain.Entity](xs []T, item T) []T {
	for i, x := range xs {
		if x.GetID() == item.GetID() {
			xs[i] = item
			return xs
		}
	}
	return append(xs, item)
}

// deleteEntity removes the element with the given id, reporting whether it existed.
func deleteEntity[T domain.Entity](xs []T, id string) ([]T, bool) {
	for i, x := range xs {
		if x.GetID() == id {
			return append(xs[:i:i], xs[i+1:]...), true
		}
	}
	return xs, false
}

// clone returns a shallow copy of a slice (so reads never alias the stored slice).
func clone[T any](xs []T) []T {
	if xs == nil {
		return nil
	}
	out := make([]T, len(xs))
	copy(out, xs)
	return out
}
