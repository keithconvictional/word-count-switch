package models

type UpdateProduct struct {
	Active bool `json:"active"`
	BodyHTML              string                `json:"bodyHtml"`
	Images []Images `json:"images"`
	Tags                  []string              `json:"tags"`
	Title                 string                `json:"title"`
	Vendor                string                `json:"vendor"`
	Variants              []Variant            `json:"variants"`
	Options               []Options             `json:"options"`
	GoogleProductCategory GoogleProductCategory `json:"googleProductCategory"`
	Attributes            map[string]interface{}           `json:"attributes"`
}
