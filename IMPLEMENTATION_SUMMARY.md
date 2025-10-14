# Receipt Manager - Implementation Summary

## Overview
A full-stack receipt management system with AI-powered OCR for processing receipt images.

## Architecture

### Backend (Go + Fiber + PostgreSQL)
- **Framework**: Fiber v3
- **Database**: PostgreSQL with normalized schema
- **AI**: Google Gemini 2.5 Flash for OCR

### Frontend (Next.js 15 + React 19)
- **Framework**: Next.js 15 with App Router
- **UI Library**: HeroUI (NextUI fork)
- **Styling**: Tailwind CSS 4

## Features Implemented

### Backend Features
1. **Receipt Upload & OCR Processing**
   - Upload receipt images (JPG, PNG)
   - AI extracts: store name, date, items, quantities, prices, discounts
   - Automatic normalization (UPPERCASE)

2. **Database Schema**
   - `stores` - Normalized store information
   - `products` - Normalized products per store
   - `receipts` - Receipt metadata with hash for duplicate detection
   - `items` - Line items linking receipts to products

3. **API Endpoints**
   - `POST /receipts/upload` - Upload and process receipt
   - `GET /receipts` - List all receipts (with pagination)
   - `GET /receipts/:id` - Get receipt details
   - `DELETE /receipts/:id` - Delete receipt
   - `PUT /items/:itemId` - Update item quantity/price

4. **Key Features**
   - UPSERT operations for stores and products
   - SHA-256 hash-based duplicate detection
   - Date-based filtering support
   - Store-based filtering support

### Frontend Features
1. **Receipts Page** (`/receipts`)
   - **Left Sidebar**: List of all receipts
     - Store name
     - Number of items
     - Purchase date
     - Total amount
     - Click to select

   - **Right Panel**: Receipt detail view
     - Full store information
     - Purchase date
     - Editable items table
     - Edit quantity and price inline
     - Total calculations
     - Delete receipt button

2. **Upload Modal**
   - Drag & drop or file select
   - Image validation (JPG, PNG)
   - Size limit (10MB)
   - Upload progress

3. **Edit Functionality**
   - Click "Edit" on any item
   - Modify quantity and price
   - Save or cancel changes
   - Automatic recalculation

## Database Schema

```sql
stores (id, name)
products (id, name, store_id) UNIQUE(name, store_id)
receipts (id, store_id, bought_date, receipt_hash, discounts)
items (id, receipt_id, product_id, quantity, price_paid)
```

## Running the Application

### Backend
```bash
cd backend
docker-compose up -d  # Start PostgreSQL
go run cmd/app/main.go
# Runs on http://localhost:8080
```

### Frontend
```bash
cd frontend
npm install
npm run dev
# Runs on http://localhost:3000
```

## Environment Variables

### Backend
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `GOOGLE_APPLICATION_CREDENTIALS` - Path to Google Cloud credentials
- `GEMINI_API_KEY` - Google Gemini API key

### Frontend
- `NEXT_PUBLIC_API_URL` - Backend API URL (default: http://localhost:8080)

## Key Technologies

### Backend
- Fiber v3 - Web framework
- pgx v5 - PostgreSQL driver
- Google Generative AI SDK - Gemini integration
- UUID - Unique identifiers

### Frontend
- Next.js 15 - React framework
- HeroUI - Component library
- Tailwind CSS 4 - Styling
- TypeScript - Type safety

## Data Flow

1. **Upload Receipt**
   ```
   User uploads image → Backend API → Gemini AI extracts data →
   → Store/Product UPSERT → Receipt created → Response sent
   ```

2. **View Receipts**
   ```
   User opens /receipts → Fetch receipt list → Display in sidebar →
   → User clicks receipt → Fetch full details → Display in detail view
   ```

3. **Edit Item**
   ```
   User clicks Edit → Inline inputs appear → User modifies →
   → Save clicked → PUT /items/:id → Refresh receipt details
   ```

4. **Delete Receipt**
   ```
   User clicks Delete → Confirmation → DELETE /receipts/:id →
   → Refresh list → Clear selection
   ```

## Future Enhancements

- Date range filtering
- Store filtering
- Export to CSV/PDF
- Receipt categories/tags
- Price history charts
- Budget tracking
- Mobile app
- Receipt sharing

## Notes

- All product names are normalized to UPPERCASE for consistency
- Receipt hash prevents duplicate uploads
- Product names are kept as-is from OCR (sizes, brands, etc.)
- Date filtering ready but not yet implemented in UI
- Pagination ready but shows all receipts currently
