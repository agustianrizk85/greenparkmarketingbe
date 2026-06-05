package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"greenpark/marketing/internal/auth"
	"greenpark/marketing/internal/domain"
	"greenpark/marketing/internal/repository"
	"greenpark/marketing/internal/service"
)

// Handler holds the dependencies for the HTTP handlers.
type Handler struct {
	svc  service.MarketingService
	auth *auth.Service
}

// NewHandler creates a Handler bound to the service and auth service.
func NewHandler(svc service.MarketingService, authSvc *auth.Service) *Handler {
	return &Handler{svc: svc, auth: authSvc}
}

/* ---------------------------- auth plumbing ---------------------------- */

type ctxKey int

const userCtxKey ctxKey = 0

func bearer(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimSpace(h[len("Bearer "):])
	}
	return ""
}

// requireAuth wraps a handler, rejecting requests without a valid session.
func (h *Handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := h.auth.Validate(bearer(r))
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), userCtxKey, u)))
	}
}

// requireAdmin wraps a handler, requiring a valid session with the admin role.
func (h *Handler) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return h.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if u, ok := r.Context().Value(userCtxKey).(domain.User); !ok || u.Role != domain.RoleAdmin {
			writeError(w, http.StatusForbidden, "butuh akses admin")
			return
		}
		next(w, r)
	})
}

// decode reads the JSON request body into a value of type T.
func decode[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		writeError(w, http.StatusBadRequest, "body JSON tidak valid: "+err.Error())
		return v, false
	}
	return v, true
}

/* ---------------------------- auth handlers ---------------------------- */

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	req, ok := decode[loginReq](w, r)
	if !ok {
		return
	}
	token, user, err := h.auth.Login(req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	h.auth.Logout(bearer(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	u, _ := r.Context().Value(userCtxKey).(domain.User)
	writeJSON(w, http.StatusOK, u)
}

/* ---------------------------- read handlers ---------------------------- */

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "marketing"})
}

func (h *Handler) dashboard(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Dashboard())
}

func (h *Handler) summary(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Summary())
}

func (h *Handler) context(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Context())
}

func (h *Handler) funnel(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Funnel())
}

func (h *Handler) kpis(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.KPIs())
}

func (h *Handler) leadQuality(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.LeadQuality())
}

func (h *Handler) handover(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Handover())
}

func (h *Handler) channels(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Channels())
}

func (h *Handler) projects(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Projects())
}

func (h *Handler) projectByName(w http.ResponseWriter, r *http.Request) {
	project, err := h.svc.ProjectByName(r.PathValue("name"))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "project not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load project")
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *Handler) assets(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Assets())
}

func (h *Handler) igAccounts(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.IGAccounts())
}

func (h *Handler) content(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Content())
}

func (h *Handler) commands(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Commands())
}

func (h *Handler) alerts(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Alerts())
}

func (h *Handler) reasonCodes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.ReasonCodes())
}

/* ---------------------------- singleton / whole-value write handlers ---------------------------- */

func (h *Handler) setContext(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.Context](w, r)
	if !ok {
		return
	}
	if err := h.svc.SetContext(v); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, v)
}

type spendReq struct {
	Spend float64 `json:"spend"`
}

func (h *Handler) setSpend(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[spendReq](w, r)
	if !ok {
		return
	}
	if err := h.svc.SetSpend(v.Spend); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) setLeadQuality(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.LeadQuality](w, r)
	if !ok {
		return
	}
	if err := h.svc.SetLeadQuality(v); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) setContent(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.Content](w, r)
	if !ok {
		return
	}
	if err := h.svc.SetContent(v); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) setAlerts(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.Alerts](w, r)
	if !ok {
		return
	}
	if err := h.svc.SetAlerts(v); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, v)
}

/* ---------------------------- collection write handlers ---------------------------- */

func (h *Handler) saveFunnelStage(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.FunnelStage](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveFunnelStage(v)
	respondSave(w, saved, err)
}

func (h *Handler) saveKPI(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.KPI](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveKPI(v)
	respondSave(w, saved, err)
}

func (h *Handler) saveHandoverItem(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.HandoverItem](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveHandoverItem(v)
	respondSave(w, saved, err)
}

func (h *Handler) saveChannel(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.Channel](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveChannel(v)
	respondSave(w, saved, err)
}

func (h *Handler) saveProject(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.Project](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveProject(v)
	respondSave(w, saved, err)
}

func (h *Handler) saveAsset(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.Asset](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveAsset(v)
	respondSave(w, saved, err)
}

func (h *Handler) saveIGAccount(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.IGAccount](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveIGAccount(v)
	respondSave(w, saved, err)
}

func (h *Handler) saveCommand(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.Command](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveCommand(v)
	respondSave(w, saved, err)
}

func (h *Handler) saveReasonCode(w http.ResponseWriter, r *http.Request) {
	v, ok := decode[domain.ReasonCode](w, r)
	if !ok {
		return
	}
	saved, err := h.svc.SaveReasonCode(v)
	respondSave(w, saved, err)
}

func respondSave[T any](w http.ResponseWriter, saved T, err error) {
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, saved)
}

// deleteHandler adapts a repository delete to an HTTP handler keyed on {id}.
func (h *Handler) deleteHandler(del func(id string) (bool, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, err := del(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}
