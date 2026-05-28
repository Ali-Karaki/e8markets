import { DefinitionList } from "@/components/shared/definition-list";
import { QueryFeedback } from "@/components/shared/query-feedback";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import type { PositionSnapshot } from "@/lib/api";
import { formatNumber } from "@/lib/format";

type PositionHistoryCardProps = {
  snapshots: PositionSnapshot[];
  instrumentName?: string;
  isLoading: boolean;
  error: unknown;
};

function snapshotItems(snapshot: PositionSnapshot, instrumentName?: string) {
  const { position } = snapshot;
  return [
    { label: "Snapshot time", value: new Date(snapshot.syncedAt).toLocaleString() },
    { label: "Sync run", value: snapshot.syncRunId },
    ...(instrumentName ? [{ label: "Symbol", value: instrumentName }] : []),
    { label: "Position ID", value: snapshot.positionId },
    ...(position.side ? [{ label: "Side", value: position.side }] : []),
    ...(position.qty !== undefined ? [{ label: "Qty", value: formatNumber(position.qty) }] : []),
    ...(position.avgPrice !== undefined
      ? [{ label: "Avg price", value: formatNumber(position.avgPrice) }]
      : []),
  ];
}

export function PositionHistoryCard({
  snapshots,
  instrumentName,
  isLoading,
  error,
}: PositionHistoryCardProps) {
  const filterNote = instrumentName
    ? `Filtered to ${instrumentName}.`
    : "All symbols for this account.";

  return (
    <Card>
      <CardHeader>
        <CardTitle>Synced snapshots</CardTitle>
        <CardDescription>
          Point-in-time copies saved when you click Sync positions. Open trades appear here
          too—this is not closed-trade history. Persists in Postgres after reload. {filterNote}
        </CardDescription>
      </CardHeader>

      <CardContent>
        <QueryFeedback
          isLoading={isLoading}
          error={error}
          isEmpty={!isLoading && !error && snapshots.length === 0}
          loadingMessage="Loading snapshots..."
          errorMessage="Failed to load snapshots"
          emptyMessage="No snapshots yet. Use Sync positions above to save what is open right now."
        >
          <ul className="space-y-4">
            {snapshots.map((snapshot) => (
              <li key={snapshot.id} className="rounded-md border p-3">
                <DefinitionList
                  items={snapshotItems(snapshot, instrumentName)}
                  columns={2}
                />
              </li>
            ))}
          </ul>
        </QueryFeedback>
      </CardContent>
    </Card>
  );
}
