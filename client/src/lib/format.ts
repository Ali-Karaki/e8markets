// accountBalance -> Account Balance
export function formatFieldLabel(key: string) {
  return key
    .replace(/([A-Z])/g, " $1")
    .replace(/^./, (s) => s.toUpperCase())
    .trim();
}

export function formatCurrency(amount: number | undefined, currency: string = "USD") {
  if (!amount) return "—";

  return amount.toLocaleString(undefined, {
    style: "currency",
    currency,
  });
}

export function formatNumber(value: number) {
  return Number.isFinite(value) ? value.toLocaleString() : value;
}
