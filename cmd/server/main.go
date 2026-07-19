package main

import (
	"log"
	"net/http"
	"os"

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

	db, err := config.NewDatabase(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Connected to MySQL")

	// Repository
	brandRepository := repositories.NewBrandRepository(db)
	categoryRepository := repositories.NewCategoryRepository(db)
	locationRepository := repositories.NewLocationRepository(db)
	itemRepository := repositories.NewItemRepository(db)
	metadataStructureRepository := repositories.NewMetadataStructureRepository(db)
	metadataRepository := repositories.NewMetadataRepository(db)

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

	// Handler
	brandHandler := handlers.NewBrandHandler(brandService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	locationHandler := handlers.NewLocationHandler(locationService)
	itemHandler := handlers.NewItemHandler(itemService)
	metadataStructureHandler := handlers.NewMetadataStructureHandler(metadataStructureService)

	// Router
	r := chi.NewRouter()

	routes.BrandRoutes(r, brandHandler)
	routes.CategoryRoutes(r, categoryHandler)
	routes.LocationRoutes(r, locationHandler)
	routes.ItemRoutes(r, itemHandler)
	routes.MetadataStructureRoutes(r, metadataStructureHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
