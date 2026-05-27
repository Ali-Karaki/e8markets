import { useMemo, useState } from "react";
import { ChevronDownIcon } from "lucide-react";

import { DefinitionList } from "@/components/shared/definition-list";
import { QueryFeedback } from "@/components/shared/query-feedback";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import type { Instrument } from "@/lib/api";
import { cn } from "@/lib/utils";

type InstrumentsCardProps = {
  instruments: Instrument[];
  selectedInstrument: Instrument | undefined;
  isLoading: boolean;
  error: unknown;
  onSelectInstrument: (tradableInstrumentId: number) => void;
};

function matchesSearch(instrument: Instrument, query: string): boolean {
  const q = query.trim().toLowerCase();
  if (!q) return true;

  return (
    instrument.name.toLowerCase().includes(q) ||
    instrument.type.toLowerCase().includes(q) ||
    (instrument.description?.toLowerCase().includes(q) ?? false)
  );
}

export function InstrumentsCard({
  instruments,
  selectedInstrument,
  isLoading,
  error,
  onSelectInstrument,
}: InstrumentsCardProps) {
  const [search, setSearch] = useState("");
  const [isListOpen, setIsListOpen] = useState(false);

  const filtered = useMemo(
    () => instruments.filter((instrument) => matchesSearch(instrument, search)),
    [instruments, search],
  );

  function handleSelectInstrument(tradableInstrumentId: number) {
    onSelectInstrument(tradableInstrumentId);
    setSearch("");
    setIsListOpen(false);
  }

  const detailItems = selectedInstrument
    ? [
        { label: "Symbol", value: selectedInstrument.name },
        { label: "Type", value: selectedInstrument.type },
        { label: "Tradable ID", value: selectedInstrument.tradableInstrumentId },
        { label: "Trading exchange", value: selectedInstrument.tradingExchange },
        { label: "Market data exchange", value: selectedInstrument.marketDataExchange },
        ...(selectedInstrument.description
          ? [{ label: "Description", value: selectedInstrument.description }]
          : []),
        ...(selectedInstrument.tradeRouteId
          ? [{ label: "Trade route ID", value: selectedInstrument.tradeRouteId }]
          : []),
        ...(selectedInstrument.infoRouteId
          ? [{ label: "Info route ID", value: selectedInstrument.infoRouteId }]
          : []),
      ]
    : [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Instruments</CardTitle>
        <CardDescription>Browse and select a tradable symbol</CardDescription>
      </CardHeader>

      <CardContent className="space-y-4">
        <QueryFeedback
          isLoading={isLoading}
          error={error}
          isEmpty={!isLoading && !error && instruments.length === 0}
          loadingMessage="Loading instruments..."
          errorMessage="Failed to load instruments"
          emptyMessage="No tradable instruments for this account."
        >
          <button
            type="button"
            onClick={() => setIsListOpen((open) => !open)}
            className={cn(
              "flex h-9 w-full items-center justify-between gap-2 rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none hover:bg-muted/50 focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50",
              !selectedInstrument && "text-muted-foreground",
            )}
            aria-expanded={isListOpen}
            aria-haspopup="listbox"
          >
            <span className="truncate">
              {selectedInstrument
                ? `${selectedInstrument.name} · ${selectedInstrument.type}`
                : "Select a symbol"}
            </span>
            <ChevronDownIcon
              className={cn("size-4 shrink-0 text-muted-foreground transition-transform", isListOpen && "rotate-180")}
            />
          </button>

          {isListOpen ? (
            <>
              <Input
                type="search"
                placeholder="Search by symbol, type, or description..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                aria-label="Search instruments"
                autoFocus
              />

              <ul className="max-h-48 space-y-1 overflow-y-auto rounded-md border p-1">
                {filtered.length === 0 ? (
                  <li className="px-3 py-2 text-sm text-muted-foreground">
                    No instruments match your search.
                  </li>
                ) : (
                  filtered.map((instrument) => {
                    const isSelected =
                      selectedInstrument?.tradableInstrumentId === instrument.tradableInstrumentId;

                    return (
                      <li key={instrument.tradableInstrumentId}>
                        <button
                          type="button"
                          onClick={() => handleSelectInstrument(instrument.tradableInstrumentId)}
                          className={cn(
                            "w-full rounded-sm px-3 py-2 text-left text-sm transition-colors hover:bg-muted",
                            isSelected && "bg-muted font-medium",
                          )}
                        >
                          <span>{instrument.name}</span>
                          <span className="ml-2 text-muted-foreground">{instrument.type}</span>
                        </button>
                      </li>
                    );
                  })
                )}
              </ul>
            </>
          ) : null}

          <DefinitionList items={detailItems} />
        </QueryFeedback>
      </CardContent>
    </Card>
  );
}
