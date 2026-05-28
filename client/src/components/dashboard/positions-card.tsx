import { DefinitionList } from "@/components/shared/definition-list";
import { QueryFeedback } from "@/components/shared/query-feedback";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import type { Position } from "@/lib/api";
import { apiErrorMessage } from "@/lib/errors";
import { formatNumber } from "@/lib/format";

type PositionsCardProps = {
  positions: Position[];
  instrumentName?: string;
  isLoading: boolean;
  error: unknown;
  isSyncing: boolean;
  syncError: unknown;
  lastSyncMessage?: string;
  emptyMessage?: string;
  onSync: () => void;
};

function positionItems(position: Position, instrumentName?: string) {
  return [
    ...(instrumentName ? [{ label: "Symbol", value: instrumentName }] : []),
    { label: "Position ID", value: position.id },
    ...(position.tradableInstrumentId
      ? [{ label: "Instrument ID", value: position.tradableInstrumentId }]
      : []),
    ...(position.side ? [{ label: "Side", value: position.side }] : []),
    ...(position.qty !== undefined
      ? [{ label: "Qty", value: formatNumber(position.qty) }]
      : []),
    ...(position.avgPrice !== undefined
      ? [{ label: "Avg price", value: formatNumber(position.avgPrice) }]
      : []),
  ];
}

export function PositionsCard({
  positions,
  instrumentName,
  isLoading,
  error,
  isSyncing,
  syncError,
  lastSyncMessage,
  emptyMessage = "No open positions for this account.",
  onSync,
}: PositionsCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-start justify-between gap-4 space-y-0">
        <div className="space-y-1.5">
          <CardTitle>Current positions</CardTitle>
          <CardDescription>
            Live open positions from TradeLocker. Sync saves a snapshot to the
            database (shown below)—it does not close trades. Auto-sync runs
            every 3 minutes when a new open position appears (tab visible).
          </CardDescription>
        </div>
        <Button type="button" size="sm" onClick={onSync} disabled={isSyncing}>
          {isSyncing ? "Syncing..." : "Sync positions"}
        </Button>
      </CardHeader>

      <CardContent className="space-y-4">
        {syncError ? (
          <p className="text-sm text-destructive">
            {apiErrorMessage(syncError, "Failed to sync positions")}
          </p>
        ) : null}
        {lastSyncMessage ? (
          <p className="text-sm text-muted-foreground">{lastSyncMessage}</p>
        ) : null}

        <QueryFeedback
          isLoading={isLoading}
          error={error}
          isEmpty={!isLoading && !error && positions.length === 0}
          loadingMessage="Loading positions..."
          errorMessage="Failed to load positions"
          emptyMessage={emptyMessage}
        >
          <ul className="space-y-4">
            {positions.map((position) => (
              <li key={position.id} className="rounded-md border p-3">
                <DefinitionList
                  items={positionItems(position, instrumentName)}
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
