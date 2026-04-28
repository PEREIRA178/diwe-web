package handlers

import (
	"strings"

	"diwe-web/internal/pocketbase"
	"diwe-web/internal/services"
	"diwe-web/internal/views"
	"diwe-web/internal/views/pages"
	"github.com/gofiber/fiber/v2"
)

type WebHandler struct {
	store *services.Store
	pb    *pocketbase.Client
}

func NewWebHandler(store *services.Store, pb *pocketbase.Client) *WebHandler {
	return &WebHandler{store: store, pb: pb}
}

func (h *WebHandler) Register(app *fiber.App) {
	app.Get("/", h.Home)
	app.Get("/soluciones", h.Solutions)
	app.Get("/soluciones/:slug", h.SolutionDetail)
	app.Get("/casos", h.Cases)
	app.Get("/casos/:industry", h.CasesByIndustry)
	app.Get("/diagnostico", h.Diagnostic)
	app.Post("/diagnostico", h.SubmitDiagnostic)
	app.Get("/p/:slug", h.ProposalPage)

	app.Get("/admin", h.AdminHome)
	app.Get("/admin/prospectos", h.AdminProspects)
	app.Post("/admin/prospectos/:id/status", h.AdminProspectStatus)
	app.Get("/admin/propuestas", h.AdminProposals)
}

func (h *WebHandler) Home(c *fiber.Ctx) error {
	return views.Render(c, pages.HomePage(h.store.Solutions(), h.store.CaseStudies("")))
}

func (h *WebHandler) Solutions(c *fiber.Ctx) error {
	return views.Render(c, pages.SolutionsPage(h.store.Solutions()))
}

func (h *WebHandler) SolutionDetail(c *fiber.Ctx) error {
	s, err := h.store.SolutionBySlug(c.Params("slug"))
	if err != nil {
		return fiber.ErrNotFound
	}
	return views.Render(c, pages.SolutionDetailPage(s))
}

func (h *WebHandler) Cases(c *fiber.Ctx) error {
	industry := strings.TrimSpace(c.Query("industry"))
	list := h.store.CaseStudies(industry)
	if c.Get("HX-Request") == "true" {
		return views.Render(c, pages.CasesList(list))
	}
	return views.Render(c, pages.CasesPage(list, h.store.Industries(), industry))
}

func (h *WebHandler) CasesByIndustry(c *fiber.Ctx) error {
	industry := services.NormalizeSlug(c.Params("industry"))
	list := h.store.CaseStudiesByIndustry(industry)
	if len(list) == 0 {
		return fiber.ErrNotFound
	}
	return views.Render(c, pages.CasesPage(list, h.store.Industries(), industry))
}

func (h *WebHandler) Diagnostic(c *fiber.Ctx) error {
	return views.Render(c, pages.DiagnosticPage())
}

func (h *WebHandler) SubmitDiagnostic(c *fiber.Ctx) error {
	prospect := services.Prospect{
		Name: c.FormValue("name"), Company: c.FormValue("company"), Industry: c.FormValue("industry"),
		Email: c.FormValue("email"), Phone: c.FormValue("phone"), MainProblem: c.FormValue("main_problem"),
		ApproxBudget: c.FormValue("approx_budget"), Urgency: c.FormValue("urgency"), Source: "web-diagnostico",
	}
	if prospect.Name == "" || prospect.Email == "" || prospect.Company == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Faltan campos obligatorios")
	}
	saved := h.store.CreateProspect(prospect)
	_ = h.pb.CreateRecord(c.UserContext(), "prospects", map[string]any{
		"name": saved.Name, "company": saved.Company, "industry": saved.Industry, "email": saved.Email,
		"phone": saved.Phone, "main_problem": saved.MainProblem, "approx_budget": saved.ApproxBudget, "urgency": saved.Urgency,
	})
	return views.Render(c, pages.DiagnosticResult(saved.Name))
}

func (h *WebHandler) ProposalPage(c *fiber.Ctx) error {
	proposal, err := h.store.TrackProposalView(c.Params("slug"))
	if err != nil {
		return fiber.ErrNotFound
	}
	_ = h.pb.CreateRecord(c.UserContext(), "events", map[string]any{"type": "proposal_view", "proposal_slug": proposal.Slug})
	_ = h.pb.CreateRecord(c.UserContext(), "proposals", map[string]any{"slug": proposal.Slug, "opened_count": proposal.OpenedCount})
	return views.Render(c, pages.ProposalPage(proposal))
}

func (h *WebHandler) AdminHome(c *fiber.Ctx) error {
	return views.Render(c, pages.AdminHomePage())
}

func (h *WebHandler) AdminProspects(c *fiber.Ctx) error {
	return views.Render(c, pages.ProspectsPage(h.store.Prospects()))
}

func (h *WebHandler) AdminProspectStatus(c *fiber.Ctx) error {
	item, err := h.store.UpdateProspectStatus(c.Params("id"), c.FormValue("status"))
	if err != nil {
		return fiber.ErrNotFound
	}
	_ = h.pb.CreateRecord(c.UserContext(), "events", map[string]any{"type": "prospect_status_updated", "prospect_id": item.ID, "status": item.Status})
	return views.Render(c, pages.ProspectRow(item))
}

func (h *WebHandler) AdminProposals(c *fiber.Ctx) error {
	return views.Render(c, pages.ProposalsPage(h.store.Proposals()))
}
