package services

import "time"

type Industry struct {
	ID          string
	Slug        string
	Name        string
	Description string
}

type Solution struct {
	ID             string
	Slug           string
	Name           string
	Problem        string
	Includes       []string
	PriceFrom      string
	EstimatedTime  string
	CTA            string
	TargetIndustry []string
}

type CaseStudy struct {
	ID               string
	Slug             string
	IndustrySlug     string
	IndustryName     string
	InitialProblem   string
	Sold             string
	Implemented      []string
	Result           string
	CTA              string
	HighlightedOffer string
}

type Prospect struct {
	ID                string
	Name              string
	Company           string
	Industry          string
	Email             string
	Phone             string
	MainProblem       string
	ApproxBudget      string
	Urgency           string
	Status            string
	CreatedAt         time.Time
	Source            string
	LastContactedNote string
}

type Proposal struct {
	ID               string
	Slug             string
	CompanyName      string
	ProblemSummary   string
	ProposedSolution string
	Items            []string
	Price            string
	Deadline         string
	OpenedCount      int
	LastOpenedAt     *time.Time
	Status           string
}

type Event struct {
	ID         string
	Type       string
	RelatedID  string
	Payload    map[string]any
	CreatedAt  time.Time
	OccurredAt time.Time
}
