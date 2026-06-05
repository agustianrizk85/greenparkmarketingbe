package repository

import (
	"greenpark/marketing/internal/domain"
	"greenpark/marketing/internal/passwd"
)

// This file holds the representative seed data for the Marketing "Qualified
// Demand Control Tower" dashboard. It mirrors the figures shown to the CEO
// war-room and is intended to be replaced by real data sources (ad platforms,
// CRM, analytics) behind the same MarketingRepository interface.
//
// Funnel/channel volumes are raw counts; monetary values are in Rupiah. KPI
// value/target/gap fields are pre-formatted display strings.

func seedContext() domain.Context {
	return domain.Context{
		Period:       "Juni 2026",
		Updated:      "02 Jun 2026 · 08:40 WIB",
		Goal:         500,
		BookingYTD:   312,
		Completeness: 88,
	}
}

func seedFunnel() []domain.FunnelStage {
	return []domain.FunnelStage{
		{Key: "Impression", Value: 4820000, Owner: "Marketing"},
		{Key: "Reach", Value: 1640000, Owner: "Marketing"},
		{Key: "Engagement", Value: 214000, Owner: "Marketing"},
		{Key: "Link Click", Value: 38500, Owner: "Marketing"},
		{Key: "LP View", Value: 24300, Owner: "Marketing / AI"},
		{Key: "Leads", Value: 3420, Owner: "Marketing"},
		{Key: "MQL", Value: 1180, Owner: "Marketing"},
		{Key: "SAL", Value: 940, Owner: "Sales"},
		{Key: "CV", Value: 612, Owner: "Sales"},
		{Key: "PV", Value: 388, Owner: "Sales"},
		{Key: "Booking", Value: 96, Owner: "Sales + Finance"},
		{Key: "KPR", Value: 71, Owner: "Keuangan / KPR"},
		{Key: "Akad", Value: 52, Owner: "Keuangan + Sales"},
		{Key: "Cash-In", Value: 48, Owner: "Finance"},
	}
}

func seedSpend() float64 { return 1920000000 } // Rp 1.92 M

func seedKPIs() []domain.KPI {
	return []domain.KPI{
		{ID: "leads", Label: "Total Leads", Value: "3.420", Target: "3.000", Gap: "+14%", Trend: []float64{2600, 2750, 2900, 3100, 3250, 3420}, Status: "good", Note: "Demand di atas target, momentum naik."},
		{ID: "mql", Label: "MQL / Qualified", Value: "1.180", Target: "1.050", Gap: "+12%", Trend: []float64{820, 870, 940, 1020, 1100, 1180}, Status: "good", Note: "Volume qualified sehat."},
		{ID: "mqlrate", Label: "MQL Rate", Value: "34.5%", Target: "35%", Gap: "-0.5pp", Trend: []float64{31.5, 32.1, 32.4, 32.9, 33.8, 34.5}, Status: "warn", Note: "Sedikit di bawah target, targeting perlu dijaga."},
		{ID: "salrate", Label: "SAL Rate", Value: "79.7%", Target: "78%", Gap: "+1.7pp", Trend: []float64{72, 74, 75, 77, 78.5, 79.7}, Status: "good", Note: "Sales menerima mayoritas MQL."},
		{ID: "spend", Label: "Total Spend", Value: "Rp 1.92 M", Target: "Rp 2 M", Gap: "-4%", Trend: []float64{1.6, 1.68, 1.74, 1.82, 1.88, 1.92}, Status: "good", Note: "Pembakaran budget terkendali."},
		{ID: "cpl", Label: "CPL", Value: "Rp 561rb", Target: "Rp 600rb", Gap: "-7%", Trend: []float64{640, 620, 605, 590, 575, 561}, Status: "good", Note: "Biaya per lead efisien."},
		{ID: "cpql", Label: "CPQL", Value: "Rp 1.6 jt", Target: "Rp 1,6 jt", Gap: "+2%", Trend: []float64{1820, 1760, 1720, 1690, 1655, 1627}, Status: "warn", Note: "Mendekati target, pantau channel boros."},
		{ID: "lead_cv", Label: "Lead → CV", Value: "17.9%", Target: "18%", Gap: "-0.1pp", Trend: []float64{15.2, 15.8, 16.4, 17.0, 17.5, 17.9}, Status: "warn", Note: "Appointment cukup, ruang perbaikan."},
		{ID: "cv_pv", Label: "CV → PV", Value: "63.4%", Target: "65%", Gap: "-1.6pp", Trend: []float64{58, 59.5, 60.8, 61.7, 62.6, 63.4}, Status: "warn", Note: "Sebagian janji visit tidak hadir."},
		{ID: "pv_book", Label: "PV → Booking", Value: "24.7%", Target: "26%", Gap: "-1.3pp", Trend: []float64{21, 22, 22.8, 23.5, 24.1, 24.7}, Status: "warn", Note: "Closing & trust di lokasi perlu diperkuat."},
		{ID: "cpb", Label: "Cost / Booking", Value: "Rp 20 jt", Target: "Rp 19 jt", Gap: "+5%", Trend: []float64{24, 23, 22, 21.5, 20.6, 20}, Status: "warn", Note: "Biaya per transaksi sedikit di atas target."},
		{ID: "roi", Label: "Marketing ROI", Value: "4.2×", Target: "4.0×", Gap: "+0.2×", Trend: []float64{3.4, 3.6, 3.8, 3.9, 4.05, 4.2}, Status: "good", Note: "Kontribusi margin positif."},
		{ID: "winning", Label: "Winning Campaign", Value: "3", Target: "≥ 3", Gap: "0", Trend: []float64{1, 1, 2, 2, 3, 3}, Status: "good", Note: "3 campaign lolos syarat winning."},
		{ID: "asset", Label: "Digital Asset Health", Value: "72", Suffix: "/100", Target: "80", Gap: "-8", Trend: []float64{61, 64, 66, 68, 70, 72}, Status: "warn", Note: "Beberapa akun IG & GBP melemahkan skor."},
	}
}

func seedLeadQuality() domain.LeadQuality {
	return domain.LeadQuality{
		Breakdown: []domain.LeadBreakdown{
			{Label: "Hot MQL (80–100)", Value: 408, Color: "hot"},
			{Label: "Warm MQL (60–79)", Value: 522, Color: "warm"},
			{Label: "Nurture (40–59)", Value: 250, Color: "nurture"},
			{Label: "Low / Noise (<40)", Value: 240, Color: "low"},
		},
		Stats: []domain.LeadStat{
			{Label: "MQL Rate", Value: "34.5%"},
			{Label: "Low Quality Rate", Value: "7.0%"},
			{Label: "Duplicate Rate", Value: "3.4%"},
		},
		TopSource:     "Meta Ads — KPR Tanpa Panik",
		BottomSource:  "TikTok Ads — Promo Umum",
		TopProject:    "Zhauz Limo",
		BottomProject: "Le Hauz Cibubur",
	}
}

func seedHandover() []domain.HandoverItem {
	return []domain.HandoverItem{
		{Label: "SAL Acceptance Rate", Value: "79.7%", Status: "good"},
		{Label: "Rejected Lead Rate", Value: "20.3%", Status: "warn"},
		{Label: "Avg Handover Time", Value: "1.8 jam", Status: "good"},
		{Label: "Avg First Response", Value: "12 mnt", Status: "good"},
		{Label: "Leads Without Owner", Value: "14", Status: "warn"},
		{Label: "SLA Breach Rate", Value: "6.1%", Status: "warn"},
	}
}

func seedChannels() []domain.Channel {
	return []domain.Channel{
		{Name: "Meta Ads", Group: "Paid", Spend: 720000000, Leads: 1180, MQL: 470, CPL: 610000, CPQL: 1532000, ROI: "4.8×", Status: "scale"},
		{Name: "TikTok Ads", Group: "Paid", Spend: 410000000, Leads: 880, MQL: 228, CPL: 466000, CPQL: 1798000, ROI: "2.9×", Status: "optimize"},
		{Name: "Google Ads", Group: "Paid", Spend: 360000000, Leads: 540, MQL: 232, CPL: 667000, CPQL: 1552000, ROI: "4.1×", Status: "scale"},
		{Name: "YouTube Ads", Group: "Paid", Spend: 120000000, Leads: 150, MQL: 41, CPL: 800000, CPQL: 2927000, ROI: "1.6×", Status: "pause"},
		{Name: "Website / LP", Group: "Owned", Spend: 0, Leads: 280, MQL: 118, CPL: 0, CPQL: 0, ROI: "—", Status: "scale"},
		{Name: "Organic IG", Group: "Owned", Spend: 0, Leads: 190, MQL: 61, CPL: 0, CPQL: 0, ROI: "—", Status: "optimize"},
		{Name: "Google Business", Group: "Trust", Spend: 0, Leads: 120, MQL: 22, CPL: 0, CPQL: 0, ROI: "—", Status: "optimize"},
		{Name: "Event / Open House", Group: "Offline", Spend: 110000000, Leads: 60, MQL: 8, CPL: 1833000, CPQL: 0, ROI: "0.9×", Status: "test"},
		{Name: "Referral / Agent", Group: "Offline", Spend: 0, Leads: 20, MQL: 0, CPL: 0, CPQL: 0, ROI: "—", Status: "test"},
	}
}

func seedProjects() []domain.Project {
	return []domain.Project{
		{Name: "The Hauz Premiere", IG: "thehauzpremiere", Demand: 84, Readiness: 88, Leads: 520, MQL: 196, Booking: 18},
		{Name: "Zhauz Limo", IG: "zhauzlimo", Demand: 90, Readiness: 78, Leads: 610, MQL: 240, Booking: 21},
		{Name: "Vertihome Serua", IG: "vertihomeserua", Demand: 58, Readiness: 82, Leads: 280, MQL: 88, Booking: 7},
		{Name: "Vertihauz Sawangan", IG: "vertihauzsawangan", Demand: 72, Readiness: 60, Leads: 360, MQL: 124, Booking: 11},
		{Name: "Vertihauz Limo 3", IG: "vertihauzlimo3", Demand: 66, Readiness: 54, Leads: 300, MQL: 96, Booking: 8},
		{Name: "Vertihauz Cibubur", IG: "vertihauzcibubur", Demand: 78, Readiness: 70, Leads: 410, MQL: 142, Booking: 13},
		{Name: "The Hauz Pancoran Mas", IG: "thehauzpancoranmas", Demand: 62, Readiness: 86, Leads: 250, MQL: 80, Booking: 6},
		{Name: "The Hauz Cilodong", IG: "thehauzcilodong", Demand: 44, Readiness: 48, Leads: 150, MQL: 38, Booking: 2},
		{Name: "Le Hauz Limo", IG: "lehauzlimo", Demand: 70, Readiness: 40, Leads: 330, MQL: 104, Booking: 6},
		{Name: "Le Hauz Cibubur", IG: "lehauzcibubur", Demand: 38, Readiness: 66, Leads: 120, MQL: 28, Booking: 1},
	}
}

func seedAssets() []domain.Asset {
	return []domain.Asset{
		{Type: "Website", Handle: "greenparkgroup.co.id", Health: 81, Active: true, Note: "GA4 ✓ · Pixel ✓ · CVR 1.4%"},
		{Type: "TikTok", Handle: "@greenparkgroup", Health: 64, Active: true, Note: "Posting rutin · views naik"},
		{Type: "YouTube", Handle: "@greenparkgroup916", Health: 52, Active: true, Note: "Frekuensi rendah · CTR lemah"},
		{Type: "Google Business", Handle: "GBP — multi project", Health: 58, Active: true, Note: "3 project belum verified"},
	}
}

func seedIGAccounts() []domain.IGAccount {
	return []domain.IGAccount{
		{Handle: "thehauzpremiere", Health: 84, Active: true, Days: 1},
		{Handle: "zhauzlimo", Health: 88, Active: true, Days: 0},
		{Handle: "vertihomeserua", Health: 70, Active: true, Days: 2},
		{Handle: "vertihauzsawangan", Health: 66, Active: true, Days: 3},
		{Handle: "vertihauzlimo3", Health: 48, Active: false, Days: 9},
		{Handle: "vertihauzcibubur", Health: 74, Active: true, Days: 1},
		{Handle: "thehauzpancoranmas", Health: 61, Active: true, Days: 4},
		{Handle: "thehauzcilodong", Health: 40, Active: false, Days: 12},
		{Handle: "lehauzlimo", Health: 69, Active: true, Days: 2},
		{Handle: "lehauzcibubur", Health: 44, Active: false, Days: 8},
	}
}

func seedContent() domain.Content {
	return domain.Content{
		Winning: []domain.WinningCampaign{
			{Name: "KPR Tanpa Panik", Project: "Zhauz Limo", Channel: "Meta", Criteria: 7, CPL: "Rp 540rb", MQL: "41%", Booking: 11},
			{Name: "Rumah Siap Akad", Project: "The Hauz Premiere", Channel: "Google", Criteria: 6, CPL: "Rp 610rb", MQL: "38%", Booking: 9},
			{Name: "Open House Sawangan", Project: "Vertihauz Sawangan", Channel: "Meta", Criteria: 5, CPL: "Rp 580rb", MQL: "36%", Booking: 7},
		},
		Best:   domain.ContentHighlight{Name: "Reels: Walkthrough Show Unit", Account: "@zhauzlimo", Metric: "188rb reach · 41 leads"},
		Worst:  domain.ContentHighlight{Name: "Carousel: Promo Umum", Account: "@lehauzcibubur", Metric: "9rb reach · 1 lead"},
		Rework: 4,
		Pause:  2,
	}
}

func seedCommands() []domain.Command {
	return []domain.Command{
		{Issue: "CPL Meta naik 30%", Cause: "Creative fatigue", Impact: "Leads mahal", Command: "Ganti 3 creative", PIC: "Digital Marketing", Deadline: "H+3", Expected: "CPL −15%", Status: "open"},
		{Issue: "MQL tinggi, CV rendah", Cause: "SAL lambat", Impact: "Leads bocor", Command: "Audit response time", PIC: "Kadep Sales", Deadline: "H+2", Expected: "CV naik", Status: "open"},
		{Issue: "Project A: LP view tinggi, leads rendah", Cause: "LP tidak meyakinkan", Impact: "Traffic bocor", Command: "Rework landing page", PIC: "Marketing + AI", Deadline: "H+5", Expected: "LP CVR naik", Status: "open"},
		{Issue: "PV tinggi, booking rendah", Cause: "Trust / site lemah", Impact: "Booking bocor", Command: "Bedah PV loss", PIC: "Sales + Teknik", Deadline: "H+7", Expected: "Booking naik", Status: "progress"},
		{Issue: "GBP project belum verified", Cause: "Trust rendah", Impact: "Direction click rendah", Command: "Verifikasi GBP", PIC: "Marketing", Deadline: "H+5", Expected: "Visit naik", Status: "open"},
	}
}

func seedAlerts() domain.Alerts {
	return domain.Alerts{
		Red: []string{
			"YouTube Ads: spend tinggi, ROI 1.6× — kandidat pause",
			"3 akun IG project tidak aktif > 7 hari",
			"3 project running ads tanpa GBP verified",
		},
		Yellow: []string{
			"PV → Booking di bawah target 1.3pp",
			"MQL Rate 0.5pp di bawah target",
			"Le Hauz Cibubur: demand lemah, konten perlu rework",
		},
		Green: []string{
			"CPL Rp 561rb — efisien",
			"SAL Rate 79.7% — sehat",
			"Marketing ROI 4.2× — positif",
		},
	}
}

func seedReasonCodes() []domain.ReasonCode {
	return []domain.ReasonCode{
		{Code: "UNR", Layer: "Leads→CV", Label: "Unreachable", Count: 96},
		{Code: "ENG", Layer: "Leads→CV", Label: "No Schedule Locked", Count: 142},
		{Code: "COM", Layer: "CV→PV", Label: "Weak Commitment", Count: 88},
		{Code: "REM", Layer: "CV→PV", Label: "Reminder Failure", Count: 54},
		{Code: "FIN", Layer: "PV→Booking", Label: "Financially Infeasible", Count: 73},
		{Code: "NST", Layer: "PV→Booking", Label: "No Next Step", Count: 61},
	}
}

// seedUsers creates the default accounts. Change these immediately in any real
// deployment. Default credentials: admin/admin123 and viewer/viewer123.
func seedUsers() []storeUser {
	mk := func(id, username, name string, role domain.Role, password string) storeUser {
		salt := passwd.NewSalt()
		return storeUser{
			ID:           id,
			Username:     username,
			Name:         name,
			Role:         role,
			Salt:         salt,
			PasswordHash: passwd.Hash(password, salt),
		}
	}
	return []storeUser{
		mk("usr-admin", "admin", "Administrator Marketing", domain.RoleAdmin, "admin123"),
		mk("usr-viewer", "viewer", "Viewer", domain.RoleViewer, "viewer123"),
	}
}
