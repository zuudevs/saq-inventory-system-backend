# API Documentation

This document describes the REST API endpoints provided by the SAQ Inventory System Backend.

## Response Wrapper format
All responses return JSON enclosed in a standard wrapper envelope.

### Success Response
```json
{
  "success": true,
  "message": "Resource retrieved successfully",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Detailed error message",
  "data": null
}
```

---

## 1. Brands API

### List All Brands
* **Method**: `GET`
* **Path**: `/brands`
* **Response**: List of brand objects.

### Get Brand by ID
* **Method**: `GET`
* **Path**: `/brands/{id}`
* **Response**: Brand object details.

### Create Brand
* **Method**: `POST`
* **Path**: `/brands`
* **Payload**:
  ```json
  {
    "name": "Brand Name"
  }
  ```
* **Response**: Created brand details.

### Update Brand
* **Method**: `PUT`
* **Path**: `/brands/{id}`
* **Payload**:
  ```json
  {
    "name": "Updated Brand Name"
  }
  ```
* **Response**: Updated brand details.

### Delete Brand
* **Method**: `DELETE`
* **Path**: `/brands/{id}`
* **Response**: Success envelope.

---

## 2. Categories API

### List All Categories
* **Method**: `GET`
* **Path**: `/categories`
* **Response**: List of category objects.

### Get Category by ID
* **Method**: `GET`
* **Path**: `/categories/{id}`
* **Response**: Category details.

### Create Category
* **Method**: `POST`
* **Path**: `/categories`
* **Payload**:
  ```json
  {
    "name": "Category Name",
    "description": "Category Description"
  }
  ```

### Update Category
* **Method**: `PUT`
* **Path**: `/categories/{id}`
* **Payload**:
  ```json
  {
    "name": "Updated Name",
    "description": "Updated Description"
  }
  ```

### Delete Category
* **Method**: `DELETE`
* **Path**: `/categories/{id}`

---

## 3. Locations API

### List All Locations
* **Method**: `GET`
* **Path**: `/locations`

### Get Location by ID
* **Method**: `GET`
* **Path**: `/locations/{id}`

### Create Location
* **Method**: `POST`
* **Path**: `/locations`
* **Payload**:
  ```json
  {
    "name": "Location Name",
    "room_code": "Room A-101",
    "description": "Location description"
  }
  ```

### Update Location
* **Method**: `PUT`
* **Path**: `/locations/{id}`

### Delete Location
* **Method**: `DELETE`
* **Path**: `/locations/{id}`

---

## 4. Items API

### List All Items
* **Method**: `GET`
* **Path**: `/items`
* **Response**: List of item records. If an item belongs to a category with a metadata structure defined, its metadata attributes will automatically be joined under the `"metadata"` object key.

### Get Item by ID
* **Method**: `GET`
* **Path**: `/items/{id}`

### Create Item
* **Method**: `POST`
* **Path**: `/items`
* **Payload**:
  ```json
  {
    "brand_id": 1,
    "category_id": 2,
    "location_id": 3,
    "asset_code": "AST-2026-0001",
    "name": "ThinkPad L14 Gen 4",
    "item_condition": "good",
    "item_status": "active",
    "notes": "Assigned to developer",
    "metadata": {
      "ram_gb": 16,
      "storage_gb": 512,
      "processor": "Intel Core i7"
    }
  }
  ```
  *Note: The fields inside the `"metadata"` object must conform to the defined metadata structure for category ID 2.*

### Update Item
* **Method**: `PUT`
* **Path**: `/items/{id}`
* **Payload**:
  ```json
  {
    "name": "Updated ThinkPad L14",
    "item_condition": "minor_damage",
    "item_status": "maintenance"
  }
  ```

### Delete Item
* **Method**: `DELETE`
* **Path**: `/items/{id}`

---

## 5. Metadata Structure API

### Get Metadata Structure by Category ID
* **Method**: `GET`
* **Path**: `/categories/{categoryId}/metadata-structure`

### Create Metadata Structure
* **Method**: `POST`
* **Path**: `/categories/{categoryId}/metadata-structure`
* **Payload**:
  ```json
  {
    "fields": [
      {
        "name": "processor",
        "label": "Processor Type",
        "type": "string",
        "nullable": false,
        "unique": false
      },
      {
        "name": "ram_gb",
        "label": "RAM size in GB",
        "type": "int",
        "nullable": false,
        "default": "8",
        "unique": false
      },
      {
        "name": "operating_system",
        "label": "Operating System",
        "type": "enum",
        "options": ["Windows", "Linux", "macOS"],
        "nullable": true,
        "unique": false
      }
    ]
  }
  ```

### Update Metadata Structure
* **Method**: `PUT`
* **Path**: `/categories/{categoryId}/metadata-structure`
* **Payload**: Same format as Create Metadata Structure. This endpoint dynamically alters the underlying table columns (adding, dropping, or modifying fields) without losing data in unchanged fields.
* **Response**: Updated metadata structure details.

### Delete Metadata Structure
* **Method**: `DELETE`
* **Path**: `/categories/{categoryId}/metadata-structure`
* **Response**: Success envelope.

---

## 6. Image Management API

### Upload Image Binary
Uploads a file to the server. Use `multipart/form-data` with a single file parameter named `file`.
* **Method**: `POST`
* **Path**: `/images/upload`
* **Payload**: `multipart/form-data`
* **Response**:
  ```json
  {
    "success": true,
    "message": "file uploaded successfully",
    "data": {
      "image_path": "images/4e9089f2-2b63-47a2-9b2f-5f07df52084c.png"
    }
  }
  ```

### Create Image Link
Associates an uploaded image path to an item or location.
* **Method**: `POST`
* **Path**: `/images`
* **Payload**:
  ```json
  {
    "item_id": 1,
    "location_id": null,
    "image_path": "images/4e9089f2-2b63-47a2-9b2f-5f07df52084c.png",
    "is_primary": 1
  }
  ```

### List All Images
* **Method**: `GET`
* **Path**: `/images`
* **Query Parameters**:
  - `?item_id={id}`: Filter images by item ID.
  - `?location_id={id}`: Filter images by location ID.

### Get Image Detail by ID
* **Method**: `GET`
* **Path**: `/images/{id}`

### Update Image Link
* **Method**: `PUT`
* **Path**: `/images/{id}`
* **Payload**:
  ```json
  {
    "is_primary": 0
  }
  ```

### Delete Image
Removes the database image record and deletes the physical file from the disk.
* **Method**: `DELETE`
* **Path**: `/images/{id}`

---

## 7. Static Asset Serving
Serves uploaded images directly.
* **Method**: `GET`
* **Path**: `/storage/*` (e.g. `/storage/images/4e9089f2-2b63-47a2-9b2f-5f07df52084c.png`)

---

## 8. Data Export API

### Export All Resources to CSV (ZIP Archive)
Exports all system resources (`Brands`, `Categories`, `Locations`, `Items`, `Images`, and `Metadata Structures`) as separate CSV files (`brands.csv`, `categories.csv`, `locations.csv`, `items.csv`, `images.csv`, `metadata_structures.csv`) bundled into a single ZIP archive.
* **Method**: `GET`
* **Path**: `/exports/csv`
* **Response Headers**:
  - `Content-Type`: `application/zip`
  - `Content-Disposition`: `attachment; filename=exports.zip`

### Export All Resources to XLSX (Excel Workbook)
Exports all system resources (`Brands`, `Categories`, `Locations`, `Items`, `Images`, and `Metadata Structures`) into a single Excel workbook containing separate worksheets for each resource (`Brands`, `Categories`, `Locations`, `Items`, `Images`, `Metadata Structures`).
* **Method**: `GET`
* **Path**: `/exports/xlsx`
* **Response Headers**:
  - `Content-Type`: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
  - `Content-Disposition`: `attachment; filename=exports.xlsx`

---

## 9. Data Import API

### Import All Resources from XLSX
Imports system resources from an Excel workbook (`.xlsx`) containing sheets for `Brands`, `Categories`, `Locations`, `Items`, and `Images`.
- **Validation Rules**:
  1. **Sheet Names**: Workbook must contain exact required sheet names (`Brands`, `Categories`, `Locations`, `Items`, `Images`).
  2. **Column Headers**: Row 1 of each sheet must match expected column headers (`ID`, `Name`, `Asset Code`, etc.).
  3. **Data Types**: Row cells are strictly validated for positive integers, non-empty text, valid booleans, and domain constraints.
- **Transaction**: Performs import within a single database transaction; if any row or constraint fails, the entire import rolls back.

* **Method**: `POST`
* **Path**: `/imports/xlsx`
* **Payload**: `multipart/form-data` with field `file`
* **Success Response**:
  ```json
  {
    "success": true,
    "message": "import completed successfully",
    "data": {
      "brands_imported": 1,
      "categories_imported": 1,
      "locations_imported": 1,
      "items_imported": 1,
      "images_imported": 1,
      "total_imported": 5
    }
  }
  ```


