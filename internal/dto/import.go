package dto

type ImportSummary struct {
	BrandsImported     int `json:"brands_imported"`
	CategoriesImported int `json:"categories_imported"`
	LocationsImported  int `json:"locations_imported"`
	ItemsImported      int `json:"items_imported"`
	ImagesImported     int `json:"images_imported"`
	TotalImported      int `json:"total_imported"`
}
