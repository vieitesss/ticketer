export interface Item {
  name: string;
  quantity: number;
  price: number;
}

export interface Receipt {
  store_name: string;
  items: Item[];
  subtotal: number;
  discounts: number;
  total_amount: number;
}
