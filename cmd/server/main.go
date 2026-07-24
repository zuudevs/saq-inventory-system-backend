package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/zuudevs/saq-inventory-system-backend/internal/config"
	"github.com/zuudevs/saq-inventory-system-backend/internal/handlers"
	"github.com/zuudevs/saq-inventory-system-backend/internal/repositories"
	"github.com/zuudevs/saq-inventory-system-backend/internal/routes"
	"github.com/zuudevs/saq-inventory-system-backend/internal/schema"
	"github.com/zuudevs/saq-inventory-system-backend/internal/services"
)

func main() {
	godotenv.Load()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH is not set")
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatal(err)
	}

	db, err := config.NewDatabase(
		dbPath,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Connected to SQLite")

	// Repository
	brandRepository := repositories.NewBrandRepository(db)
	categoryRepository := repositories.NewCategoryRepository(db)
	locationRepository := repositories.NewLocationRepository(db)
	itemRepository := repositories.NewItemRepository(db)
	metadataStructureRepository := repositories.NewMetadataStructureRepository(db)
	metadataRepository := repositories.NewMetadataRepository(db)
	imageRepository := repositories.NewImageRepository(db)

	// Schema
	schemaService := schema.NewService(db)

	// Service
	brandService := &services.BrandService{
		BrandRepository: brandRepository,
	}

	categoryService := &services.CategoryService{
		CategoryRepository: categoryRepository,
	}

	locationService := &services.LocationService{
		LocationRepository: locationRepository,
	}

	itemService := &services.ItemService{
		DB:                          db,
		ItemRepository:              itemRepository,
		CategoryRepository:          categoryRepository,
		MetadataStructureRepository: metadataStructureRepository,
		MetadataRepository:          metadataRepository,
	}

	metadataStructureService := &services.MetadataStructureService{
		MetadataStructureRepository: metadataStructureRepository,
		CategoryRepository:          categoryRepository,
		SchemaService:               schemaService,
	}

	imageService := &services.ImageService{
		DB:                 db,
		StoragePath:        storagePath,
		ImageRepository:    imageRepository,
		ItemRepository:     itemRepository,
		LocationRepository: locationRepository,
	}

	exportService := &services.ExportService{
		DB:                          db,
		BrandRepository:             brandRepository,
		CategoryRepository:          categoryRepository,
		ItemRepository:              itemRepository,
		LocationRepository:          locationRepository,
		ImageRepository:             imageRepository,
		MetadataStructureRepository: metadataStructureRepository,
	}

	importService := &services.ImportService{
		DB:                 db,
		BrandRepository:    brandRepository,
		CategoryRepository: categoryRepository,
		ItemRepository:     itemRepository,
		LocationRepository: locationRepository,
		ImageRepository:    imageRepository,
	}

	// Handler
	brandHandler := handlers.NewBrandHandler(brandService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	locationHandler := handlers.NewLocationHandler(locationService)
	itemHandler := handlers.NewItemHandler(itemService)
	metadataStructureHandler := handlers.NewMetadataStructureHandler(metadataStructureService)
	imageHandler := handlers.NewImageHandler(imageService, storagePath)
	exportHandler := handlers.NewExportHandler(exportService)
	importHandler := handlers.NewImportHandler(importService)

	// Router
	r := chi.NewRouter()

	routes.HealthRoutes(r)
	routes.BrandRoutes(r, brandHandler)
	routes.CategoryRoutes(r, categoryHandler)
	routes.LocationRoutes(r, locationHandler)
	routes.ItemRoutes(r, itemHandler)
	routes.MetadataStructureRoutes(r, metadataStructureHandler)
	routes.ImageRoutes(r, imageHandler)
	routes.StorageRoutes(r, storagePath)
	routes.ExportRoutes(r, exportHandler)
	routes.ImportRoutes(r, importHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
