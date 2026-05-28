import { DefinitionList } from "@/components/shared/definition-list";
import { QueryFeedback } from "@/components/shared/query-feedback";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { Account } from "@/lib/api";
import { formatCurrency } from "@/lib/format";

type AccountsCardProps = {
  accounts: Account[];
  selectedAccount: Account | undefined;
  isLoading: boolean;
  error: unknown;
  onSelectAccount: (id: string) => void;
};

function accountDescription(count: number): string {
  if (count === 0) return "No accounts available";
  if (count === 1) return "Your trading account";
  return "Select a trading account";
}

export function AccountsCard({
  accounts,
  selectedAccount,
  isLoading,
  error,
  onSelectAccount,
}: AccountsCardProps) {
  const summaryItems = selectedAccount
    ? [
        { label: "Account name", value: selectedAccount.name },
        { label: "Account ID", value: selectedAccount.id },
        { label: "Account number", value: selectedAccount.accNum },
        { label: "Status", value: selectedAccount.status },
        { label: "Currency", value: selectedAccount.currency },
        {
          label: "Balance",
          value: formatCurrency(
            selectedAccount.accountBalance,
            selectedAccount.currency,
          ),
        },
      ]
    : [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Accounts</CardTitle>
        <CardDescription>{accountDescription(accounts.length)}</CardDescription>
      </CardHeader>

      <CardContent className="space-y-4">
        <QueryFeedback
          isLoading={isLoading}
          error={error}
          isEmpty={!isLoading && !error && accounts.length === 0}
          loadingMessage="Loading accounts..."
          errorMessage="Failed to load accounts"
          emptyMessage="No trading accounts were returned for this login."
        >
          {accounts.length === 1 && selectedAccount ? (
            <p className="text-sm text-muted-foreground">
              Single trading account (no switcher needed).
            </p>
          ) : null}

          {accounts.length > 1 && selectedAccount ? (
            <Select value={selectedAccount.id} onValueChange={onSelectAccount}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select account" />
              </SelectTrigger>

              <SelectContent>
                {accounts.map((account) => (
                  <SelectItem key={account.id} value={account.id}>
                    {account.name} — #{account.accNum} ({account.id}) —{" "}
                    {account.currency}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          ) : null}

          {accounts.length > 0 ? <DefinitionList items={summaryItems} /> : null}
        </QueryFeedback>
      </CardContent>
    </Card>
  );
}
