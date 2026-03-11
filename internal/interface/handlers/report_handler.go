package handlers

import (
	"net/http"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/dto/responses"
	"github.com/zuxt268/berry/internal/interface/middleware"
	"github.com/zuxt268/berry/internal/usecase"
)

type ReportHandler struct {
	ga4UseCase       usecase.GA4ReportUseCase
	gscUseCase       usecase.GSCReportUseCase
	gbpUseCase       usecase.GBPReportUseCase
	instagramUseCase usecase.InstagramReportUseCase
	lineUseCase      usecase.LineReportUseCase
}

func NewReportHandler(
	ga4 usecase.GA4ReportUseCase,
	gsc usecase.GSCReportUseCase,
	gbp usecase.GBPReportUseCase,
	instagram usecase.InstagramReportUseCase,
	line usecase.LineReportUseCase,
) *ReportHandler {
	return &ReportHandler{
		ga4UseCase:       ga4,
		gscUseCase:       gsc,
		gbpUseCase:       gbp,
		instagramUseCase: instagram,
		lineUseCase:      line,
	}
}

func (h *ReportHandler) parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	if fromStr == "" || toStr == "" {
		return time.Time{}, time.Time{}, domain.ErrInvalidArgument
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidArgument
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidArgument
	}
	return from, to, nil
}

func (h *ReportHandler) getUserID(r *http.Request) (uint64, error) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		return 0, domain.ErrUnauthorized
	}
	return user.ID, nil
}

// GA4Reports GET /api/users/ga4/reports?from=2006-01-02&to=2006-01-02
func (h *ReportHandler) GA4Reports(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		HandleError(w, err)
		return
	}
	from, to, err := h.parseDateRange(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	reports, err := h.ga4UseCase.GetReports(r.Context(), userID, from, to)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp := make([]responses.GA4DailyReportResponse, len(reports))
	for i, rpt := range reports {
		resp[i] = responses.ToGA4DailyReportResponse(rpt)
	}
	respondJSON(w, http.StatusOK, map[string]any{"reports": resp})
}

// GSCReports GET /api/users/gsc/reports?from=2006-01-02&to=2006-01-02
func (h *ReportHandler) GSCReports(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		HandleError(w, err)
		return
	}
	from, to, err := h.parseDateRange(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	reports, err := h.gscUseCase.GetReports(r.Context(), userID, from, to)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp := make([]responses.GSCDailyReportResponse, len(reports))
	for i, rpt := range reports {
		resp[i] = responses.ToGSCDailyReportResponse(rpt)
	}
	respondJSON(w, http.StatusOK, map[string]any{"reports": resp})
}

// GBPReports GET /api/users/gbp/reports?from=2006-01-02&to=2006-01-02
func (h *ReportHandler) GBPReports(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		HandleError(w, err)
		return
	}
	from, to, err := h.parseDateRange(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	reports, err := h.gbpUseCase.GetReports(r.Context(), userID, from, to)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp := make([]responses.GBPDailyReportResponse, len(reports))
	for i, rpt := range reports {
		resp[i] = responses.ToGBPDailyReportResponse(rpt)
	}
	respondJSON(w, http.StatusOK, map[string]any{"reports": resp})
}

// InstagramReports GET /api/users/instagram/reports?from=2006-01-02&to=2006-01-02
func (h *ReportHandler) InstagramReports(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		HandleError(w, err)
		return
	}
	from, to, err := h.parseDateRange(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	reports, err := h.instagramUseCase.GetReports(r.Context(), userID, from, to)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp := make([]responses.InstagramDailyReportResponse, len(reports))
	for i, rpt := range reports {
		resp[i] = responses.ToInstagramDailyReportResponse(rpt)
	}
	respondJSON(w, http.StatusOK, map[string]any{"reports": resp})
}

// LineReports GET /api/users/line/reports?from=2006-01-02&to=2006-01-02
func (h *ReportHandler) LineReports(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		HandleError(w, err)
		return
	}
	from, to, err := h.parseDateRange(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	reports, err := h.lineUseCase.GetReports(r.Context(), userID, from, to)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp := make([]responses.LineDailyReportResponse, len(reports))
	for i, rpt := range reports {
		resp[i] = responses.ToLineDailyReportResponse(rpt)
	}
	respondJSON(w, http.StatusOK, map[string]any{"reports": resp})
}