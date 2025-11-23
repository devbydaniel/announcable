package shared

// BaseTemplateData provides common data structure for all page templates
type BaseTemplateData struct {
	Title                 string
	HasActiveSubscription bool
}

// PageData interface allows type-safe access to base template data
// Implement this interface for page-specific data structs that embed BaseTemplateData
type PageData interface {
	GetBaseData() *BaseTemplateData
}
