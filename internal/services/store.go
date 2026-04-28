package services

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	mu          sync.RWMutex
	industries  []Industry
	solutions   []Solution
	caseStudies []CaseStudy
	prospects   []Prospect
	proposals   []Proposal
	events      []Event
}

func NewStoreWithSeed() *Store {
	now := time.Now()
	industries := []Industry{
		{ID: "ind-agr", Slug: "agricola", Name: "Agrícola", Description: "Empresas agrícolas que venden a distribuidores y necesitan operación en terreno."},
		{ID: "ind-centro", Slug: "centro-comercial", Name: "Centro comercial", Description: "Centros comerciales con alto flujo de visitas y necesidad de señalética digital."},
		{ID: "ind-edu", Slug: "educacion", Name: "Educación", Description: "Instituciones que necesitan captar y orientar mejor a sus estudiantes y familias."},
		{ID: "ind-legal", Slug: "legal-profesional", Name: "Legal y profesional", Description: "Estudios jurídicos y servicios profesionales que venden confianza y rapidez."},
	}

	solutions := []Solution{
		{ID: "sol-web-b2b", Slug: "sitio-web-b2b", Name: "Sitio web comercial B2B", Problem: "Tu empresa depende de referencias y no convierte tráfico en reuniones.", Includes: []string{"Estructura comercial por servicios", "Casos por industria", "Formulario de diagnóstico", "Integración WhatsApp"}, PriceFrom: "USD 1,200", EstimatedTime: "2 a 4 semanas", CTA: "Quiero un sitio que venda", TargetIndustry: []string{"agricola", "legal-profesional"}},
		{ID: "sol-erp-simple", Slug: "erp-simple-operaciones", Name: "ERP simple de operaciones", Problem: "Tu equipo usa planillas sueltas y no tiene visibilidad diaria.", Includes: []string{"Módulos de pedidos y tareas", "Roles por área", "Panel de seguimiento", "Capacitación inicial"}, PriceFrom: "USD 2,500", EstimatedTime: "4 a 8 semanas", CTA: "Necesito ordenar operación", TargetIndustry: []string{"agricola"}},
		{ID: "sol-totem-directorio", Slug: "totem-directorio-pantallas", Name: "Directorio, tótem y pantallas", Problem: "Tus visitantes no encuentran locales ni promociones a tiempo.", Includes: []string{"Directorio web", "Interfaz táctil para tótem", "Gestión de pantallas", "Actualización remota"}, PriceFrom: "USD 3,000", EstimatedTime: "5 a 8 semanas", CTA: "Quiero mejorar experiencia en piso", TargetIndustry: []string{"centro-comercial", "educacion"}},
		{ID: "sol-whatsapp", Slug: "whatsapp-automatizado", Name: "WhatsApp automatizado comercial", Problem: "Recibes consultas pero el equipo responde tarde o sin guión.", Includes: []string{"Flujos por tipo de consulta", "Plantillas de respuesta", "Derivación a ejecutivo", "Reporte básico de conversaciones"}, PriceFrom: "USD 900", EstimatedTime: "1 a 2 semanas", CTA: "Quiero responder más rápido", TargetIndustry: []string{"educacion", "legal-profesional"}},
	}

	caseStudies := []CaseStudy{
		{ID: "case-agr-01", Slug: "agricola-web-erp-tablets", IndustrySlug: "agricola", IndustryName: "Agrícola", InitialProblem: "El equipo de campo reportaba por WhatsApp y la gerencia consolidaba datos a mano.", Sold: "Sitio web comercial + ERP simple + tablets para supervisores", Implemented: []string{"Sitio orientado a distribuidores", "ERP simple para pedidos y stock", "Formularios móviles en tablets"}, Result: "En 90 días redujeron 35% el tiempo de cierre de pedidos y aumentaron solicitudes web calificadas.", CTA: "Quiero algo parecido", HighlightedOffer: "Operación de campo ordenada y ventas con trazabilidad."},
		{ID: "case-centro-01", Slug: "centro-comercial-directorio-totem-pantallas", IndustrySlug: "centro-comercial", IndustryName: "Centro comercial", InitialProblem: "Los visitantes no encontraban tiendas y administración no tenía canal central para campañas.", Sold: "Directorio web + tótem interactivo + red de pantallas", Implemented: []string{"Mapa de tiendas por categoría", "Tótem con búsqueda rápida", "Gestor de campañas para pantallas"}, Result: "Subió 22% el uso de promociones internas y bajaron consultas presenciales repetidas.", CTA: "Quiero algo parecido", HighlightedOffer: "Más visibilidad para locales y mejor experiencia de visita."},
		{ID: "case-edu-01", Slug: "educacion-web-totem-pantallas-whatsapp-test", IndustrySlug: "educacion", IndustryName: "Educación", InitialProblem: "Admisiones tardaba en responder y campus tenía señalización débil.", Sold: "Web + tótem + pantallas + WhatsApp + test vocacional", Implemented: []string{"Landing por carrera", "Tótem de orientación en campus", "WhatsApp automatizado para admisiones", "Test con recomendación inicial"}, Result: "Duplicaron contactos calificados en admisión y disminuyeron llamadas repetidas de orientación.", CTA: "Quiero algo parecido", HighlightedOffer: "Captación y orientación con menos carga administrativa."},
		{ID: "case-legal-01", Slug: "legal-web-agendamiento-whatsapp", IndustrySlug: "legal-profesional", IndustryName: "Legal y profesional", InitialProblem: "El estudio perdía consultas por respuestas tardías y agenda manual.", Sold: "Sitio web + agendamiento + WhatsApp", Implemented: []string{"Servicios claros por especialidad", "Agenda de reuniones", "Botón directo a WhatsApp"}, Result: "Aumentaron 40% las reuniones agendadas desde web en 60 días.", CTA: "Quiero algo parecido", HighlightedOffer: "Más reuniones con clientes correctos sin perseguir prospectos."},
	}

	nowCopy := now
	proposals := []Proposal{
		{ID: "prop-01", Slug: "agro-norte-fase1", CompanyName: "Agro Norte", ProblemSummary: "Dependen de mensajes sueltos para pedidos y pierden seguimiento comercial.", ProposedSolution: "Sitio web comercial + módulo simple de pedidos + panel de seguimiento.", Items: []string{"Sitio comercial por líneas de servicio", "Formulario diagnóstico conectado a prospectos", "Panel simple de pedidos"}, Price: "USD 4,800", Deadline: "6 semanas", OpenedCount: 1, LastOpenedAt: &nowCopy, Status: "enviada"},
	}

	return &Store{industries: industries, solutions: solutions, caseStudies: caseStudies, prospects: []Prospect{}, proposals: proposals, events: []Event{}}
}

func (s *Store) Industries() []Industry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := append([]Industry(nil), s.industries...)
	return result
}

func (s *Store) Solutions() []Solution {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := append([]Solution(nil), s.solutions...)
	return result
}

func (s *Store) SolutionBySlug(slug string) (Solution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.solutions {
		if item.Slug == slug {
			return item, nil
		}
	}
	return Solution{}, errors.New("solution not found")
}

func (s *Store) CaseStudies(industrySlug string) []CaseStudy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	filtered := make([]CaseStudy, 0)
	for _, item := range s.caseStudies {
		if industrySlug == "" || item.IndustrySlug == industrySlug {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (s *Store) CaseStudiesByIndustry(industrySlug string) []CaseStudy {
	return s.CaseStudies(industrySlug)
}

func (s *Store) CreateProspect(p Prospect) Prospect {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ID = uuid.NewString()
	p.CreatedAt = time.Now()
	if p.Status == "" {
		p.Status = "nuevo"
	}
	s.prospects = append(s.prospects, p)
	s.events = append(s.events, Event{ID: uuid.NewString(), Type: "prospect_created", RelatedID: p.ID, CreatedAt: time.Now(), OccurredAt: time.Now(), Payload: map[string]any{"industry": p.Industry}})
	return p
}

func (s *Store) Prospects() []Prospect {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := append([]Prospect(nil), s.prospects...)
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result
}

func (s *Store) UpdateProspectStatus(id, status string) (Prospect, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.prospects {
		if s.prospects[i].ID == id {
			s.prospects[i].Status = status
			s.events = append(s.events, Event{ID: uuid.NewString(), Type: "prospect_status_updated", RelatedID: id, CreatedAt: time.Now(), OccurredAt: time.Now(), Payload: map[string]any{"status": status}})
			return s.prospects[i], nil
		}
	}
	return Prospect{}, fmt.Errorf("prospect not found")
}

func (s *Store) Proposals() []Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := append([]Proposal(nil), s.proposals...)
	return result
}

func (s *Store) ProposalBySlug(slug string) (Proposal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.proposals {
		if p.Slug == slug {
			return p, nil
		}
	}
	return Proposal{}, errors.New("proposal not found")
}

func (s *Store) TrackProposalView(slug string) (Proposal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.proposals {
		if s.proposals[i].Slug == slug {
			now := time.Now()
			s.proposals[i].OpenedCount++
			s.proposals[i].LastOpenedAt = &now
			s.events = append(s.events, Event{ID: uuid.NewString(), Type: "proposal_view", RelatedID: s.proposals[i].ID, CreatedAt: now, OccurredAt: now, Payload: map[string]any{"slug": slug}})
			return s.proposals[i], nil
		}
	}
	return Proposal{}, errors.New("proposal not found")
}

func (s *Store) Events() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := append([]Event(nil), s.events...)
	return result
}

func NormalizeSlug(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}
