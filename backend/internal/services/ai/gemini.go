package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/vieitesss/ticketer/internal/models"
	"github.com/vieitesss/ticketer/pkg/logger"
	"google.golang.org/genai"
)

type GeminiService struct {
	client *genai.Client
}

func NewGeminiService(ctx context.Context) (*GeminiService, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &GeminiService{client: client}, nil
}

func (s *GeminiService) getMimeType(imagePath string) string {
	ext := strings.ToLower(filepath.Ext(imagePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "image/jpeg" // default fallback
	}
}

func (s *GeminiService) identifyStore(ctx context.Context, imageData []byte, mimeType string) (string, error) {
	log.Info("Identifying store from receipt")

	prompt := `## ROLE

You are a store identifier for Spanish shopping receipts.

## INSTRUCTION

Identify the store name by looking at the top of the receipt.

## STEPS

1. Look for the store name in the first few lines of the receipt
2. Identify known Spanish brands: ALDI, Carrefour, Mercadona, Lidl, etc.
3. Extract the exact name as it appears

## EXPECTATION

Answer JUST with the store name in UPPERCASE letters.

## NARROWING

- ONLY the store name, in UPPERCASE
- If you cannot identify it, respond "UNKNOWN"
- Do not include any extra text or punctuation`

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: imageData, MIMEType: mimeType}},
	}

	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		[]*genai.Content{{Parts: parts}},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to identify store: %w", err)
	}

	storeName := strings.TrimSpace(result.Text())
	storeName = strings.ToUpper(storeName)

	log.Info("Store identified", "store", storeName)
	return storeName, nil
}

func (s *GeminiService) getStorePrompt(storeName string) string {
	switch storeName {
	case "ALDI":
		log.Debug("Using ALDI-specific prompt template")
		return `## INSTRUCTION

Process the ALDI receipt line by line and output the information you found in JSON format.
Think step by step.

## STEPS

1. Go through each line of the receipt from top to bottom until the line with "-----":
2. You will find three types of lines:
  - Lines indicating quantity (Type A)
  - Lines indicating product details (Type B)
3. The quantity lines (Type A) always come BEFORE the product lines (Type B).

   Type A - Line with quantity and price details format "QUANTITY|WEIGHT [unit] x PRICE €[/unit]"
   - It may not exist for every product.
   - If it exists, it is ALWAYS the line BEFORE the product line (Type B).
   - Examples: "2 x 0,92 €" or "0,508 kg x 7,85 €/kg"
   - The first number BEFORE "x" indicates the quantity or weight of the next product, and is stored as "quantity".
   - The number AFTER "x" and BEFORE "€" is the price per unit or per kg, and is stored as "price".
   - This values correspond to the product in the NEXT line (Type B)

   Type B - Line with product details, format "NAME PRICE € CODE"
   - NAME is the product "name".
   - If there is a PREVIOUS line (Type A): use the "quantity" and "price" from that line for this product.
   - If there is NO PREVIOUS line (Type A): use "quantity" = 1 and "price" = number BEFORE "€" in this line.

4. Save the product "{name, quantity, price}" with the correct values after following the rules above.
5. Repeat for all products in the receipt.

## EXPECTATION

Extract all products with their correct "name", "quantity" and "price".

## NARROWING

- Anything between brackets [] is optional
- "quantity" can be decimal ("0,508", "0,67", etc.)
- The quantity or weight of each product is the one in the line BEFORE with "x", if it exists
- The quantity of each product is 1, if there is NO line BEFORE with "x"
- "price" is:
  - For products with line "Type A" BEFORE: the number BETWEEN "x" and "€" in the line BEFORE
  - For products WITHOUT line "Type A" BEFORE: the number BEFORE "€" in the product line
- NO invent products or quantities
- THINK step by step`

	case "CARREFOUR EXPRESS", "CARREFOUR":
		log.Debug("Using CARREFOUR-specific prompt template")
		return `## INSTRUCTION

Process the CARREFOUR [EXPRESS] receipt line by line and output the information you found in JSON format.
Think step by step.

## STEPS

1. Go through each line of the receipt from top to bottom:
2. You will find three types of lines:
   - Lines indicating product details (Type A)
   - Lines indicating quantity and price (Type B)
   - Lines indicating discounts (Type C)
3. The quantity lines (Type B) always come AFTER the product lines (Type A).

   Type A - Line with product details, format "NAME PRICE" or "NAME CODE"
   - The format is "product name" and "price", or "product name" and "code"
   - Extract: "name" = product name, "price" = number at the end of the line (if it's a PRICE)
   - If the last element is a CODE (letters and numbers), the PRICE and QUANTITY appear in the NEXT line
   - Move to the NEXT line to verify "quantity" and "price"
   - If the product name starts with "DESCUENTO" (discount):
     - Add the ABSOLUTE VALUE to "product_discount" (it will always be positive)
     - DO NOT add this product to the products list

   Type B - Line with "x" format "X x ( N )        Y"
   - Examples: "2 x ( 0,92 )        1,84" or "0,67 x ( 7,85 )        5,26"
   - X = quantity or weight of the PREVIOUS product (from the line above)
   - N = price per unit
   - Y = total price to pay (normally X * N = Y)
   - IMPORTANT: Store for the product in the PREVIOUS line:
     - "quantity" = X (first number before "x")
     - "price" = N (number in parentheses)

   If you find another product line (Type A) without having found line Type B for the previous product:
   - The previous product is a single item (without multiple units)
   - "quantity" = 1
   - "price" = the PRICE from the product line (Type A)
   - Register it and continue with the new product

4. At the end of the receipt: if there is a product left without quantity:
   - "quantity" = 1
   - "price" = the PRICE from the product line (Type A)

5. Get the applied discounts:
   - Look for the line "DESCUENTOS", below the "A PAGAR" line and in BOLD, and extract its amount
   - If there are no discounts, do not set it
   - Store the value as "discounts"

## EXPECTATION

Extract all products with their correct "name", "quantity" and "price", and the "discounts" to apply.

## NARROWING

- "quantity" can be decimal ("0,67", "0,508", etc.)
- For products with line "Type B" AFTER: "price" is N and quantity is X, in "X x ( N )        Y"
- For products "Type A" with PRICE on the same line: "price" is the number before € and "quantity" = 1
- "discounts" is the number in the "DESCUENTOS" line in BOLD, or null if it doesn't exist
- DO NOT invent products or quantities
- THINK step by step`

	default:
		log.Debug("Using generic prompt template")
		return `## INSTRUCTION

Extract products from the receipt by identifying columns or data structure.
Think step by step.

## STEPS

1. Look for column headers in the first lines of the receipt:
   - Typical columns: "Cant", "Cantidad", "Uds", "Precio", "Importe", "Total"
   - If you find headers: use that structure for all products

2. If there ARE identified columns:
   - For each product line, extract:
     - "name" = text in name/description column
     - "quantity" = number in "Cant", "Cantidad" or "Unidades" column (if it exists)
     - "price" = price per unit, as float, in "Precio", "Unidad" or similar column (if it exists)
   - If there is no quantity column: "quantity" = 1

3. If there are NO clear columns:
   - For each line that looks like a product:
     - "name" = product text
     - "quantity" = number on the same line (if it exists), or 1 if not
     - "price" = float number before the € symbol (or the last number on the line)

5. Look for applied discounts:
   - Lines with words "DESCUENTO", "AHORRO"
   - Extract the number as "discounts", or 0 if there are no discounts

## EXPECTATION

Extract all products with "{name, quantity, price}" and the "discounts" of the receipt.

## NARROWING

- "quantity" can be decimal ("0,67", "1,5", etc.) or integer
- If there is no explicit quantity, always use "quantity" = 1
- "price" is the number before the € symbol (or the last number on the line)
- DO NOT invent products or quantities that you don't see on the receipt
- THINK step by step`
	}
}

func (s *GeminiService) ProcessReceipt(ctx context.Context, imagePath string) (*models.Receipt, error) {
	log.Info("Starting receipt processing", "path", imagePath)

	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	// Determine MIME type from file extension
	mimeType := s.getMimeType(imagePath)
	log.Debug("Detected image format", "mimeType", mimeType)

	// Step 1: Identify store
	storeName, err := s.identifyStore(ctx, imageData, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to identify store: %w", err)
	}

	// Step 2: Get store-specific prompt
	storePrompt := s.getStorePrompt(storeName)

	fullPrompt := fmt.Sprintf(`## ROLE

You are a specialized processor for supermarket receipts, in this case from the supermarket %s.

%s

- Generate a JSON with the following format:
{
  "store_name": string,
  "items": [{"name": string, "quantity": float, "price": float}],
  "discounts": float | null
}`, storeName, storePrompt)

	log.Info("Sending receipt to model for extraction", "store", storeName)

	parts := []*genai.Part{
		{Text: fullPrompt},
		{InlineData: &genai.Blob{Data: imageData, MIMEType: mimeType}},
	}

	// Define response schema for structured output
	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"store_name": {
				Type:     genai.TypeString,
				Nullable: genai.Ptr(true),
			},
			"items": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name":     {Type: genai.TypeString},
						"quantity": {Type: genai.TypeNumber},
						"price":    {Type: genai.TypeNumber},
					},
					PropertyOrdering: []string{"name", "quantity", "price"},
				},
			},
			"discounts":    {Type: genai.TypeNumber, Nullable: genai.Ptr(true)},
		},
		PropertyOrdering: []string{
			"store_name",
			"items",
			"discounts",
		},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   schema,
	}

	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		[]*genai.Content{{Parts: parts}},
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	responseText := result.Text()
	log.Info("Received response from model")
	log.Debug("Raw model response", "response", logger.FormatJSON(responseText))

	var receipt models.Receipt
	if err := json.Unmarshal([]byte(responseText), &receipt); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w\nResponse: %s", err, responseText)
	}

	log.Info("Successfully parsed receipt",
		"store_name", receipt.StoreName,
		"items", len(receipt.Items),
		"discounts", receipt.Discounts)
	log.Info("Receipt processing completed")

	return &receipt, nil
}
