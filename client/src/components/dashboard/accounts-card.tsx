import { DefinitionList } from "@/components/shared/definition-list";
import { QueryFeedback } from "@/components/shared/query-feedback";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
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

export function AccountsCard({
  accounts,
  selectedAccount,
  isLoading,
  error,
  onSelectAccount,
}: AccountsCardProps) {
  const summaryItems = selectedAccount
    ? [
        { label: "Account ID", value: selectedAccount.id },
        { label: "Status", value: selectedAccount.status },
        { label: "Currency", value: selectedAccount.currency },
        {
          label: "Balance",
          value: formatCurrency(selectedAccount.accountBalance, selectedAccount.currency),
        },
      ]
    : [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Accounts</CardTitle>
        <CardDescription>Select a trading account</CardDescription>
      </CardHeader>

      <CardContent className="space-y-4">
        <QueryFeedback
          isLoading={isLoading}
          error={error}
          loadingMessage="Loading accounts..."
          errorMessage="Failed to load accounts"
        >
          {accounts.length > 1 && selectedAccount && (
            <Select value={selectedAccount.id} onValueChange={onSelectAccount}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select account" />
              </SelectTrigger>

              <SelectContent>
                {accounts.map((account) => (
                  <SelectItem key={account.id} value={account.id}>
                    {account.name} ({account.id}) — {account.currency}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}

          <DefinitionList items={summaryItems} />
        </QueryFeedback>
      </CardContent>
    </Card>
  );
}
