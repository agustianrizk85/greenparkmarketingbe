package http

import "net/http"

// NewRouter wires all routes to the handler and applies global middleware.
//
// Access tiers:
//   - public: health check + login
//   - authenticated (any logged-in user): all dashboard reads + me/logout
//   - admin only: every master-data write
func NewRouter(h *Handler, allowOrigin string) http.Handler {
	mux := http.NewServeMux()

	// ---- public ----
	mux.HandleFunc("GET /api/health", h.health)
	mux.HandleFunc("POST /api/auth/login", h.login)

	// ---- authenticated session ----
	mux.HandleFunc("GET /api/auth/me", h.requireAuth(h.me))
	mux.HandleFunc("POST /api/auth/logout", h.requireAuth(h.logout))

	// ---- reads (authenticated) ----
	mux.HandleFunc("GET /api/dashboard", h.requireAuth(h.dashboard))
	mux.HandleFunc("GET /api/summary", h.requireAuth(h.summary))
	mux.HandleFunc("GET /api/context", h.requireAuth(h.context))
	mux.HandleFunc("GET /api/funnel", h.requireAuth(h.funnel))
	mux.HandleFunc("GET /api/kpis", h.requireAuth(h.kpis))
	mux.HandleFunc("GET /api/lead-quality", h.requireAuth(h.leadQuality))
	mux.HandleFunc("GET /api/handover", h.requireAuth(h.handover))
	mux.HandleFunc("GET /api/channels", h.requireAuth(h.channels))
	mux.HandleFunc("GET /api/projects", h.requireAuth(h.projects))
	mux.HandleFunc("GET /api/projects/{name}", h.requireAuth(h.projectByName))
	mux.HandleFunc("GET /api/assets", h.requireAuth(h.assets))
	mux.HandleFunc("GET /api/ig-accounts", h.requireAuth(h.igAccounts))
	mux.HandleFunc("GET /api/content", h.requireAuth(h.content))
	mux.HandleFunc("GET /api/commands", h.requireAuth(h.commands))
	mux.HandleFunc("GET /api/alerts", h.requireAuth(h.alerts))
	mux.HandleFunc("GET /api/reason-codes", h.requireAuth(h.reasonCodes))

	// ---- singleton / whole-value writes (admin) ----
	mux.HandleFunc("PUT /api/context", h.requireAdmin(h.setContext))
	mux.HandleFunc("PUT /api/spend", h.requireAdmin(h.setSpend))
	mux.HandleFunc("PUT /api/lead-quality", h.requireAdmin(h.setLeadQuality))
	mux.HandleFunc("PUT /api/content", h.requireAdmin(h.setContent))
	mux.HandleFunc("PUT /api/alerts", h.requireAdmin(h.setAlerts))

	// ---- collection writes (admin) ----
	mux.HandleFunc("POST /api/funnel", h.requireAdmin(h.saveFunnelStage))
	mux.HandleFunc("DELETE /api/funnel/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteFunnelStage)))
	mux.HandleFunc("POST /api/kpis", h.requireAdmin(h.saveKPI))
	mux.HandleFunc("DELETE /api/kpis/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteKPI)))
	mux.HandleFunc("POST /api/handover", h.requireAdmin(h.saveHandoverItem))
	mux.HandleFunc("DELETE /api/handover/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteHandoverItem)))
	mux.HandleFunc("POST /api/channels", h.requireAdmin(h.saveChannel))
	mux.HandleFunc("DELETE /api/channels/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteChannel)))
	mux.HandleFunc("POST /api/projects", h.requireAdmin(h.saveProject))
	mux.HandleFunc("DELETE /api/projects/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteProject)))
	mux.HandleFunc("POST /api/assets", h.requireAdmin(h.saveAsset))
	mux.HandleFunc("DELETE /api/assets/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteAsset)))
	mux.HandleFunc("POST /api/ig-accounts", h.requireAdmin(h.saveIGAccount))
	mux.HandleFunc("DELETE /api/ig-accounts/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteIGAccount)))
	mux.HandleFunc("POST /api/commands", h.requireAdmin(h.saveCommand))
	mux.HandleFunc("DELETE /api/commands/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteCommand)))
	mux.HandleFunc("POST /api/reason-codes", h.requireAdmin(h.saveReasonCode))
	mux.HandleFunc("DELETE /api/reason-codes/{id}", h.requireAdmin(h.deleteHandler(h.svc.DeleteReasonCode)))

	// ---- maintenance (admin) ----
	mux.HandleFunc("POST /api/admin/seed", h.requireAdmin(h.reseed))
	mux.HandleFunc("POST /api/admin/clear", h.requireAdmin(h.clear))

	return chain(mux, logger, cors(allowOrigin))
}
