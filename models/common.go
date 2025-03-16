// Package models provides data models for the Apple Music API.
package models

// Resource represents a resource in the Apple Music API.
type Resource struct {
	// The type of the resource.
	Type string `json:"type"`

	// The unique identifier for the resource.
	ID string `json:"id"`

	// The URL for the resource.
	HREF string `json:"href,omitempty"`
}

// Artwork represents artwork for a resource.
type Artwork struct {
	// The width of the artwork in pixels.
	Width int `json:"width"`

	// The height of the artwork in pixels.
	Height int `json:"height"`

	// The URL for the artwork.
	URL string `json:"url"`

	// The background color to use when displaying the artwork.
	BgColor string `json:"bgColor,omitempty"`

	// The text color to use when displaying the artwork.
	TextColor1 string `json:"textColor1,omitempty"`

	// The secondary text color to use when displaying the artwork.
	TextColor2 string `json:"textColor2,omitempty"`

	// The tertiary text color to use when displaying the artwork.
	TextColor3 string `json:"textColor3,omitempty"`

	// The quaternary text color to use when displaying the artwork.
	TextColor4 string `json:"textColor4,omitempty"`
}

// PlayParameters represents play parameters for a resource.
type PlayParameters struct {
	// The unique identifier for the resource.
	ID string `json:"id"`

	// The type of the resource.
	Kind string `json:"kind"`

	// Whether the resource is a preview.
	IsLibrary bool `json:"isLibrary,omitempty"`

	// The preview URL for the resource.
	PreviewURL string `json:"previewURL,omitempty"`

	// The catalog ID for the resource.
	CatalogID string `json:"catalogId,omitempty"`
}

// EditorialNotes represents editorial notes for a resource.
type EditorialNotes struct {
	// The standard editorial notes.
	Standard string `json:"standard,omitempty"`

	// The short editorial notes.
	Short string `json:"short,omitempty"`
}

// Preview represents a preview for a resource.
type Preview struct {
	// The URL for the preview.
	URL string `json:"url"`

	// Whether the preview is playable.
	Playable bool `json:"playable,omitempty"`
}

// Storefront represents a storefront.
type Storefront struct {
	// The storefront ID.
	ID string `json:"id"`

	// The storefront default language tag.
	DefaultLanguageTag string `json:"defaultLanguageTag"`

	// The storefront name.
	Name string `json:"name"`

	// The storefront supported language tags.
	SupportedLanguageTags []string `json:"supportedLanguageTags"`
}

// Pagination represents pagination information.
type Pagination struct {
	// The offset for the next page.
	Next string `json:"next,omitempty"`

	// The total number of resources.
	Total int `json:"total,omitempty"`

	// The limit for the current page.
	Limit int `json:"limit,omitempty"`

	// The offset for the current page.
	Offset int `json:"offset,omitempty"`
}

// Chart represents a chart.
type Chart struct {
	// The name of the chart.
	Name string `json:"name"`

	// The chart order.
	Order string `json:"order"`

	// The chart data.
	Data []Resource `json:"data"`
}

// Relationship represents a relationship between resources.
type Relationship struct {
	// The relationship data.
	Data []Resource `json:"data"`

	// The relationship href.
	HREF string `json:"href,omitempty"`

	// The relationship next href.
	Next string `json:"next,omitempty"`
}

// QueryParameters represents query parameters for the Apple Music API.
type QueryParameters struct {
	// The number of resources to fetch.
	Limit int `json:"limit,omitempty"`

	// The offset for the resources to fetch.
	Offset int `json:"offset,omitempty"`

	// The fields to include in the response.
	Include []string `json:"include,omitempty"`

	// The fields to exclude from the response.
	Exclude []string `json:"exclude,omitempty"`

	// The language tag for the response.
	LanguageTag string `json:"l,omitempty"`

	// The storefront for the response.
	Storefront string `json:"storefront,omitempty"`
}

// Response represents a response from the Apple Music API.
type Response struct {
	// The response data.
	Data []interface{} `json:"data,omitempty"`

	// The response errors.
	Errors []interface{} `json:"errors,omitempty"`

	// The response meta.
	Meta map[string]interface{} `json:"meta,omitempty"`

	// The response results.
	Results map[string]interface{} `json:"results,omitempty"`
}

