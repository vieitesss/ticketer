export interface Store {
  id: string;
  name: string;
}

export interface Item {
  id: string;
  product_id: string;
  product_name: string;
  quantity: number;
  price_paid: number;
  subtotal: number;
}

export interface Receipt {
  id: string;
  store: Store;
  bought_date: string; // ISO 8601: YYYY-MM-DD
  items: Item[];
  subtotal: number;
  discounts: number;
  total_amount: number;
}

export interface ReceiptListItem {
  id: string;
  store_name: string;
  item_count: number;
  bought_date: string; // ISO 8601: YYYY-MM-DD
  total_amount: number;
}

export interface UpdateItemRequest {
  quantity: number;
  price_paid: number;
}
