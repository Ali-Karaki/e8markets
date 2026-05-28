import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useMemo, useState } from "react";

import { AccountStateCard } from "@/components/dashboard/account-state-card";
import { AccountsCard } from "@/components/dashboard/accounts-card";
import { DashboardHeader } from "@/components/dashboard/dashboard-header";
import { InstrumentsCard } from "@/components/dashboard/instruments-card";
import { PositionHistoryCard } from "@/components/dashboard/position-history-card";
import { PositionsCard } from "@/components/dashboard/positions-card";
import { useAccountState, useAccounts } from "@/lib/accounts";
import { useLogout, useSession } from "@/lib/auth";
import { findPreferredInstrument, useInstruments } from "@/lib/instruments";
import {
  filterPositionsByInstrument,
  usePositionHistory,
  usePositions,
  useSyncPositions,
} from "@/lib/positions";

export const Route = createFileRoute("/_authenticated/")({
  component: DashboardPage,
});

function DashboardPage() {
  const navigate = useNavigate();
  const session = useSession();
  const accounts = useAccounts();
  const logout = useLogout();

  const accountList = accounts.data?.accounts ?? [];
  const [selectedId, setSelectedId] = useState<string>("");
  const [selectedInstrumentId, setSelectedInstrumentId] = useState<number | null>(null);
  const [lastSyncMessage, setLastSyncMessage] = useState<string | undefined>();

  const selectedAccount = useMemo(
    () => accountList.find((a) => a.id === selectedId) ?? accountList[0],
    [accountList, selectedId],
  );

  useEffect(() => {
    if (accountList.length > 0 && !selectedId) {
      setSelectedId(accountList[0].id);
    }
  }, [accountList, selectedId]);

  useEffect(() => {
    setSelectedInstrumentId(null);
  }, [selectedAccount?.id]);

  const accountState = useAccountState(selectedAccount?.id, selectedAccount?.accNum);
  const instruments = useInstruments(selectedAccount?.id, selectedAccount?.accNum);
  const instrumentList = instruments.data?.instruments ?? [];

  const selectedInstrument = useMemo(() => {
    if (selectedInstrumentId !== null) {
      const found = instrumentList.find((i) => i.tradableInstrumentId === selectedInstrumentId);
      if (found) return found;
    }
    return findPreferredInstrument(instrumentList);
  }, [instrumentList, selectedInstrumentId]);

  const instrumentFilter = selectedInstrument?.tradableInstrumentId;
  const positions = usePositions(selectedAccount?.id, selectedAccount?.accNum);
  const allPositions = positions.data?.positions ?? [];
  const filteredPositions = useMemo(
    () => filterPositionsByInstrument(allPositions, instrumentFilter),
    [allPositions, instrumentFilter],
  );
  const positionsEmptyMessage = useMemo(() => {
    if (allPositions.length > 0 && filteredPositions.length === 0 && selectedInstrument?.name) {
      return `No open positions for ${selectedInstrument.name}.`;
    }
    return "No open positions for this account.";
  }, [allPositions.length, filteredPositions.length, selectedInstrument?.name]);
  const positionHistory = usePositionHistory(
    selectedAccount?.id,
    selectedAccount?.accNum,
    instrumentFilter,
  );
  const syncPositions = useSyncPositions();

  useEffect(() => {
    if (instrumentList.length > 0 && selectedInstrumentId === null) {
      const preferred = findPreferredInstrument(instrumentList);
      if (preferred) {
        setSelectedInstrumentId(preferred.tradableInstrumentId);
      }
    }
  }, [instrumentList, selectedInstrumentId]);

  useEffect(() => {
    setLastSyncMessage(undefined);
  }, [selectedAccount?.id, instrumentFilter]);

  async function handleSyncPositions() {
    if (!selectedAccount) return;

    try {
      const result = await syncPositions.mutateAsync({
        accountId: selectedAccount.id,
        accNum: selectedAccount.accNum,
      });
      setLastSyncMessage(
        `Saved ${result.recordsStored} snapshot${result.recordsStored === 1 ? "" : "s"} to the database (sync run #${result.syncRunId}).`,
      );
    } catch {
      setLastSyncMessage(undefined);
    }
  }

  async function handleLogout() {
    await logout.mutateAsync();
    navigate({ to: "/login" });
  }

  return (
    <div className="mx-auto flex min-h-svh w-full max-w-4xl flex-col gap-6 p-8">
      <DashboardHeader
        email={session.data?.email}
        onLogout={handleLogout}
        isLoggingOut={logout.isPending}
      />
      <AccountsCard
        accounts={accountList}
        selectedAccount={selectedAccount}
        onSelectAccount={setSelectedId}
        isLoading={accounts.isLoading}
        error={accounts.error}
      />
      <InstrumentsCard
        instruments={instrumentList}
        selectedInstrument={selectedInstrument}
        onSelectInstrument={setSelectedInstrumentId}
        isLoading={instruments.isLoading}
        error={instruments.error}
      />
      <AccountStateCard
        state={accountState.data?.state ?? {}}
        isLoading={accountState.isLoading}
        error={accountState.error}
      />
      <PositionsCard
        positions={filteredPositions}
        instrumentName={selectedInstrument?.name}
        isLoading={positions.isLoading}
        error={positions.error}
        isSyncing={syncPositions.isPending}
        syncError={syncPositions.error}
        lastSyncMessage={lastSyncMessage}
        emptyMessage={positionsEmptyMessage}
        onSync={handleSyncPositions}
      />
      <PositionHistoryCard
        snapshots={positionHistory.data?.snapshots ?? []}
        instrumentName={selectedInstrument?.name}
        isLoading={positionHistory.isLoading}
        error={positionHistory.error}
      />
    </div>
  );
}
