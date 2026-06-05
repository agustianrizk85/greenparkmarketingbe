// Package domain holds the core business entities of the Marketing "Qualified
// Demand Control Tower" dashboard. These types are the single source of truth
// for the data shape and carry no dependency on transport or storage concerns.
//
// The dashboard tracks the full demand funnel — from Impression to Cash-In —
// plus lead quality, channel/project performance, digital assets, content,
// the CEO command panel and the alert system.
//
// Collection entities carry a synthetic _id (EntID) — the stable handle the
// master-data admin uses to update/delete a row, independent of any editable
// business field.
package domain

// Status is the common traffic-light health indicator used across KPIs and the
// handover panel. It matches the front-end status tokens.
type Status string

const (
	StatusGood Status = "good" // healthy
	StatusWarn Status = "warn" // watch
	StatusBad  Status = "bad"  // alert
)

// Context holds the top-line dashboard header figures.
type Context struct {
	Period       string `json:"period"`
	Updated      string `json:"updated"`
	Goal         int    `json:"goal"`         // booking goal for the year (unit)
	BookingYTD   int    `json:"bookingYTD"`   // booking achieved year-to-date (unit)
	Completeness int    `json:"completeness"` // data completeness %
}

// FunnelStage is one stage of the full demand funnel (Impression → Cash-In).
type FunnelStage struct {
	EntID string `json:"_id"`
	Key   string `json:"key"`
	Value int    `json:"value"`
	Owner string `json:"owner"` // accountable department — key into the owner colour map
}

// KPI is one North Star indicator on the ribbon. Value/target/gap are
// pre-formatted display strings; trend feeds the sparkline.
type KPI struct {
	EntID  string    `json:"_id"`
	ID     string    `json:"id"`
	Label  string    `json:"label"`
	Value  string    `json:"value"`
	Suffix string    `json:"suffix,omitempty"`
	Target string    `json:"target"`
	Gap    string    `json:"gap"`
	Trend  []float64 `json:"trend"`
	Status string    `json:"status"` // good | warn | bad
	Note   string    `json:"note"`
}

// LeadBreakdown is one quality band of the MQL scoring distribution.
type LeadBreakdown struct {
	Label string `json:"label"`
	Value int    `json:"value"`
	Color string `json:"color"` // hot | warm | nurture | low
}

// LeadStat is a single headline lead-quality statistic.
type LeadStat struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// LeadQuality is the lead quality & MQL scoring panel (edited as a singleton).
type LeadQuality struct {
	Breakdown     []LeadBreakdown `json:"breakdown"`
	Stats         []LeadStat      `json:"stats"`
	TopSource     string          `json:"topSource"`
	BottomSource  string          `json:"bottomSource"`
	TopProject    string          `json:"topProject"`
	BottomProject string          `json:"bottomProject"`
}

// HandoverItem is one metric of the MQL → SAL handover accountability panel.
type HandoverItem struct {
	EntID  string `json:"_id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Status string `json:"status"` // good | warn | bad
}

// Channel is one acquisition channel in the performance matrix.
type Channel struct {
	EntID  string  `json:"_id"`
	Name   string  `json:"name"`
	Group  string  `json:"group"`  // Paid | Owned | Trust | Offline
	Spend  float64 `json:"spend"`  // Rupiah
	Leads  int     `json:"leads"`
	MQL    int     `json:"mql"`
	CPL    float64 `json:"cpl"`    // cost per lead (Rupiah)
	CPQL   float64 `json:"cpql"`   // cost per qualified lead (Rupiah)
	ROI    string  `json:"roi"`    // display string, e.g. "4.8×"
	Status string  `json:"status"` // scale | optimize | pause | test
}

// Project is one development project's demand & readiness scorecard.
type Project struct {
	EntID     string `json:"_id"`
	Name      string `json:"name"`
	IG        string `json:"ig"`        // instagram handle (without @)
	Demand    int    `json:"demand"`    // 0..100
	Readiness int    `json:"readiness"` // 0..100
	Leads     int    `json:"leads"`
	MQL       int    `json:"mql"`
	Booking   int    `json:"booking"`
}

// Asset is one entry in the digital asset registry (website, social, GBP…).
type Asset struct {
	EntID  string `json:"_id"`
	Type   string `json:"type"`
	Handle string `json:"handle"`
	Health int    `json:"health"` // 0..100
	Active bool   `json:"active"`
	Note   string `json:"note"`
}

// IGAccount is one project Instagram account in the registry grid.
type IGAccount struct {
	EntID  string `json:"_id"`
	Handle string `json:"handle"`
	Health int    `json:"health"` // 0..100
	Active bool   `json:"active"`
	Days   int    `json:"days"` // days since last post
}

// WinningCampaign is a campaign that passed the winning-campaign criteria.
type WinningCampaign struct {
	Name     string `json:"name"`
	Project  string `json:"project"`
	Channel  string `json:"channel"`
	Criteria int    `json:"criteria"` // criteria met (out of 8)
	CPL      string `json:"cpl"`
	MQL      string `json:"mql"`
	Booking  int    `json:"booking"`
}

// ContentHighlight is the best/worst creative call-out.
type ContentHighlight struct {
	Name    string `json:"name"`
	Account string `json:"account"`
	Metric  string `json:"metric"`
}

// Content is the content & winning-campaign panel (edited as a singleton).
type Content struct {
	Winning []WinningCampaign `json:"winning"`
	Best    ContentHighlight  `json:"best"`
	Worst   ContentHighlight  `json:"worst"`
	Rework  int               `json:"rework"`
	Pause   int               `json:"pause"`
}

// Command is one row of the CEO command panel: issue → command → PIC → deadline.
type Command struct {
	EntID    string `json:"_id"`
	Issue    string `json:"issue"`
	Cause    string `json:"cause"`
	Impact   string `json:"impact"`
	Command  string `json:"command"`
	PIC      string `json:"pic"`
	Deadline string `json:"deadline"`
	Expected string `json:"expected"`
	Status   string `json:"status"` // open | progress | done
}

// Alerts is the red/yellow/green alert system (edited as a singleton).
type Alerts struct {
	Red    []string `json:"red"`
	Yellow []string `json:"yellow"`
	Green  []string `json:"green"`
}

// ReasonCode is one leakage reason code mapped to a funnel layer.
type ReasonCode struct {
	EntID string `json:"_id"`
	Code  string `json:"code"`
	Layer string `json:"layer"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// Summary holds the executive figures derived from the rest of the data set.
type Summary struct {
	Goal         int     `json:"goal"`         // booking goal (unit)
	BookingYTD   int     `json:"bookingYTD"`   // booking achieved YTD (unit)
	Progress     float64 `json:"progress"`     // bookingYTD / goal %
	TotalLeads   int     `json:"totalLeads"`   // total leads across channels
	TotalMQL     int     `json:"totalMQL"`     // total qualified leads
	TotalSpend   float64 `json:"totalSpend"`   // total marketing spend (Rupiah)
	TotalBooking int     `json:"totalBooking"` // bookings from the funnel
	RedAlerts    int     `json:"redAlerts"`    // count of active red alerts
	OpenCommands int     `json:"openCommands"` // command rows not yet done
}

// Dashboard is the full payload consumed by the front-end in a single call.
type Dashboard struct {
	Context     Context       `json:"context"`
	Funnel      []FunnelStage `json:"funnel"`
	Spend       float64       `json:"spend"`
	KPIs        []KPI         `json:"kpis"`
	LeadQuality LeadQuality   `json:"leadQuality"`
	Handover    []HandoverItem `json:"handover"`
	Channels    []Channel     `json:"channels"`
	Projects    []Project     `json:"projects"`
	Assets      []Asset       `json:"assets"`
	IGAccounts  []IGAccount   `json:"igAccounts"`
	Content     Content       `json:"content"`
	Commands    []Command     `json:"commands"`
	Alerts      Alerts        `json:"alerts"`
	ReasonCodes []ReasonCode  `json:"reasonCodes"`
	Summary     Summary       `json:"summary"`
}

// Entity is implemented by every CRUD collection element. The synthetic _id is
// the stable handle the master-data admin uses to update/delete a row.
type Entity interface {
	GetID() string
}

func (f FunnelStage) GetID() string  { return f.EntID }
func (k KPI) GetID() string          { return k.EntID }
func (h HandoverItem) GetID() string { return h.EntID }
func (c Channel) GetID() string      { return c.EntID }
func (p Project) GetID() string      { return p.EntID }
func (a Asset) GetID() string        { return a.EntID }
func (i IGAccount) GetID() string    { return i.EntID }
func (c Command) GetID() string      { return c.EntID }
func (r ReasonCode) GetID() string   { return r.EntID }

// Role enumerates the access levels for a dashboard user.
type Role string

const (
	RoleAdmin  Role = "admin"  // full master-data access
	RoleViewer Role = "viewer" // read-only dashboard access
)

// User is a dashboard account. Password material is never serialised to clients
// (json:"-"); it only lives in the persisted store.
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Role         Role   `json:"role"`
	PasswordHash string `json:"-"`
	Salt         string `json:"-"`
}
