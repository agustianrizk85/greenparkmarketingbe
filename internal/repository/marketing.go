// Package repository defines storage access for the Marketing dashboard and
// ships a file-backed, in-memory implementation seeded with representative data.
// Writes are mutex-guarded and persisted to a JSON file so master-data edits
// survive restarts. Swapping in a database-backed store only requires satisfying
// the MarketingRepository interface.
package repository

import (
	"errors"

	"greenpark/marketing/internal/domain"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("resource not found")

// MarketingRepository is the persistence boundary for the dashboard data set.
type MarketingRepository interface {
	// ---- reads ----
	Context() domain.Context
	Funnel() []domain.FunnelStage
	Spend() float64
	KPIs() []domain.KPI
	LeadQuality() domain.LeadQuality
	Handover() []domain.HandoverItem
	Channels() []domain.Channel
	Projects() []domain.Project
	ProjectByName(name string) (domain.Project, error)
	Assets() []domain.Asset
	IGAccounts() []domain.IGAccount
	Content() domain.Content
	Commands() []domain.Command
	Alerts() domain.Alerts
	ReasonCodes() []domain.ReasonCode

	// ---- singleton / whole-value writes ----
	SetContext(domain.Context) error
	SetSpend(float64) error
	SetLeadQuality(domain.LeadQuality) error
	SetContent(domain.Content) error
	SetAlerts(domain.Alerts) error

	// ---- collection writes (Save = create when _id empty, else update) ----
	SaveFunnelStage(domain.FunnelStage) (domain.FunnelStage, error)
	DeleteFunnelStage(id string) (bool, error)
	SaveKPI(domain.KPI) (domain.KPI, error)
	DeleteKPI(id string) (bool, error)
	SaveHandoverItem(domain.HandoverItem) (domain.HandoverItem, error)
	DeleteHandoverItem(id string) (bool, error)
	SaveChannel(domain.Channel) (domain.Channel, error)
	DeleteChannel(id string) (bool, error)
	SaveProject(domain.Project) (domain.Project, error)
	DeleteProject(id string) (bool, error)
	SaveAsset(domain.Asset) (domain.Asset, error)
	DeleteAsset(id string) (bool, error)
	SaveIGAccount(domain.IGAccount) (domain.IGAccount, error)
	DeleteIGAccount(id string) (bool, error)
	SaveCommand(domain.Command) (domain.Command, error)
	DeleteCommand(id string) (bool, error)
	SaveReasonCode(domain.ReasonCode) (domain.ReasonCode, error)
	DeleteReasonCode(id string) (bool, error)

	// ---- users (auth) ----
	Users() []domain.User
	UserByUsername(username string) (domain.User, error)
	UserByID(id string) (domain.User, error)

	// ---- maintenance ----
	Reseed() error // restore the built-in example data (users kept)
	Clear() error  // wipe all dashboard data (users kept)
}
