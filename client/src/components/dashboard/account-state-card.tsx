import { DefinitionList } from "@/components/shared/definition-list";
import { QueryFeedback } from "@/components/shared/query-feedback";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { formatFieldLabel, formatNumber } from "@/lib/format";

type AccountStateCardProps = {
  state: Record<string, number>;
  isLoading: boolean;
  error: unknown;
};

export function AccountStateCard({
  state,
  isLoading,
  error,
}: AccountStateCardProps) {
  const items = Object.entries(state).map(([key, value]) => ({
    label: formatFieldLabel(key),
    value: formatNumber(value),
  }));

  return (
    <Card>
      <CardHeader>
        <CardTitle>Account details</CardTitle>
        <CardDescription>
          Balance, equity, margin, and related fields
        </CardDescription>
      </CardHeader>

      <CardContent>
        <QueryFeedback
          isLoading={isLoading}
          error={error}
          isEmpty={items.length === 0}
          loadingMessage="Loading account state..."
          errorMessage="Failed to load account details"
          emptyMessage="No account state available."
        >
          <DefinitionList items={items} columns={3} />
        </QueryFeedback>
      </CardContent>
    </Card>
  );
}
