// Package service holds the business logic of the Marketing dashboard. It
// composes the repository data and computes the derived executive summary,
// keeping transport handlers thin. Write use-cases delegate to the repository
// (which persists), so master-data edits flow straight back into the dashboard
// read.
package service

import (
	"math"

	"greenpark/marketing/internal/domain"
	"greenpark/marketing/internal/repository"
)

// MarketingService exposes the read and write use-cases of the dashboard.
type MarketingService interface {
	// reads
	Dashboard() domain.Dashboard
	Summary() domain.Summary
	Context() domain.Context
	Funnel() []domain.FunnelStage
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

	// singleton / whole-value writes
	SetContext(domain.Context) error
	SetSpend(float64) error
	SetLeadQuality(domain.LeadQuality) error
	SetContent(domain.Content) error
	SetAlerts(domain.Alerts) error

	// collection writes
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
}

type marketingService struct {
	repo repository.MarketingRepository
}

// New returns a MarketingService backed by the given repository.
func New(repo repository.MarketingRepository) MarketingService {
	return &marketingService{repo: repo}
}

// Dashboard assembles the full payload including the derived summary.
func (s *marketingService) Dashboard() domain.Dashboard {
	return domain.Dashboard{
		Context:     s.repo.Context(),
		Funnel:      s.repo.Funnel(),
		Spend:       s.repo.Spend(),
		KPIs:        s.repo.KPIs(),
		LeadQuality: s.repo.LeadQuality(),
		Handover:    s.repo.Handover(),
		Channels:    s.repo.Channels(),
		Projects:    s.repo.Projects(),
		Assets:      s.repo.Assets(),
		IGAccounts:  s.repo.IGAccounts(),
		Content:     s.repo.Content(),
		Commands:    s.repo.Commands(),
		Alerts:      s.repo.Alerts(),
		ReasonCodes: s.repo.ReasonCodes(),
		Summary:     s.Summary(),
	}
}

// Summary computes the executive figures from the context, funnel, channels,
// commands and alerts.
func (s *marketingService) Summary() domain.Summary {
	ctx := s.repo.Context()
	channels := s.repo.Channels()
	funnel := s.repo.Funnel()
	commands := s.repo.Commands()
	alerts := s.repo.Alerts()

	var totalLeads, totalMQL int
	var totalSpend float64
	for _, c := range channels {
		totalLeads += c.Leads
		totalMQL += c.MQL
		totalSpend += c.Spend
	}

	totalBooking := 0
	for _, f := range funnel {
		if f.Key == "Booking" {
			totalBooking = f.Value
		}
	}

	openCommands := 0
	for _, c := range commands {
		if c.Status != "done" {
			openCommands++
		}
	}

	progress := 0.0
	if ctx.Goal > 0 {
		progress = math.Round(float64(ctx.BookingYTD)/float64(ctx.Goal)*1000) / 10
	}

	return domain.Summary{
		Goal:         ctx.Goal,
		BookingYTD:   ctx.BookingYTD,
		Progress:     progress,
		TotalLeads:   totalLeads,
		TotalMQL:     totalMQL,
		TotalSpend:   totalSpend,
		TotalBooking: totalBooking,
		RedAlerts:    len(alerts.Red),
		OpenCommands: openCommands,
	}
}

/* ---- reads ---- */

func (s *marketingService) Context() domain.Context           { return s.repo.Context() }
func (s *marketingService) Funnel() []domain.FunnelStage       { return s.repo.Funnel() }
func (s *marketingService) KPIs() []domain.KPI                 { return s.repo.KPIs() }
func (s *marketingService) LeadQuality() domain.LeadQuality    { return s.repo.LeadQuality() }
func (s *marketingService) Handover() []domain.HandoverItem    { return s.repo.Handover() }
func (s *marketingService) Channels() []domain.Channel         { return s.repo.Channels() }
func (s *marketingService) Projects() []domain.Project         { return s.repo.Projects() }
func (s *marketingService) Assets() []domain.Asset             { return s.repo.Assets() }
func (s *marketingService) IGAccounts() []domain.IGAccount     { return s.repo.IGAccounts() }
func (s *marketingService) Content() domain.Content            { return s.repo.Content() }
func (s *marketingService) Commands() []domain.Command         { return s.repo.Commands() }
func (s *marketingService) Alerts() domain.Alerts              { return s.repo.Alerts() }
func (s *marketingService) ReasonCodes() []domain.ReasonCode   { return s.repo.ReasonCodes() }

func (s *marketingService) ProjectByName(name string) (domain.Project, error) {
	return s.repo.ProjectByName(name)
}

/* ---- singleton / whole-value writes ---- */

func (s *marketingService) SetContext(c domain.Context) error        { return s.repo.SetContext(c) }
func (s *marketingService) SetSpend(v float64) error                 { return s.repo.SetSpend(v) }
func (s *marketingService) SetLeadQuality(lq domain.LeadQuality) error {
	return s.repo.SetLeadQuality(lq)
}
func (s *marketingService) SetContent(c domain.Content) error { return s.repo.SetContent(c) }
func (s *marketingService) SetAlerts(a domain.Alerts) error   { return s.repo.SetAlerts(a) }

/* ---- collection writes ---- */

func (s *marketingService) SaveFunnelStage(v domain.FunnelStage) (domain.FunnelStage, error) {
	return s.repo.SaveFunnelStage(v)
}
func (s *marketingService) DeleteFunnelStage(id string) (bool, error) {
	return s.repo.DeleteFunnelStage(id)
}
func (s *marketingService) SaveKPI(v domain.KPI) (domain.KPI, error) { return s.repo.SaveKPI(v) }
func (s *marketingService) DeleteKPI(id string) (bool, error)        { return s.repo.DeleteKPI(id) }
func (s *marketingService) SaveHandoverItem(v domain.HandoverItem) (domain.HandoverItem, error) {
	return s.repo.SaveHandoverItem(v)
}
func (s *marketingService) DeleteHandoverItem(id string) (bool, error) {
	return s.repo.DeleteHandoverItem(id)
}
func (s *marketingService) SaveChannel(v domain.Channel) (domain.Channel, error) {
	return s.repo.SaveChannel(v)
}
func (s *marketingService) DeleteChannel(id string) (bool, error) { return s.repo.DeleteChannel(id) }
func (s *marketingService) SaveProject(v domain.Project) (domain.Project, error) {
	return s.repo.SaveProject(v)
}
func (s *marketingService) DeleteProject(id string) (bool, error) { return s.repo.DeleteProject(id) }
func (s *marketingService) SaveAsset(v domain.Asset) (domain.Asset, error) {
	return s.repo.SaveAsset(v)
}
func (s *marketingService) DeleteAsset(id string) (bool, error) { return s.repo.DeleteAsset(id) }
func (s *marketingService) SaveIGAccount(v domain.IGAccount) (domain.IGAccount, error) {
	return s.repo.SaveIGAccount(v)
}
func (s *marketingService) DeleteIGAccount(id string) (bool, error) {
	return s.repo.DeleteIGAccount(id)
}
func (s *marketingService) SaveCommand(v domain.Command) (domain.Command, error) {
	return s.repo.SaveCommand(v)
}
func (s *marketingService) DeleteCommand(id string) (bool, error) { return s.repo.DeleteCommand(id) }
func (s *marketingService) SaveReasonCode(v domain.ReasonCode) (domain.ReasonCode, error) {
	return s.repo.SaveReasonCode(v)
}
func (s *marketingService) DeleteReasonCode(id string) (bool, error) {
	return s.repo.DeleteReasonCode(id)
}
