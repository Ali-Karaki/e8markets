import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useMemo, useState } from "react";

import { AccountStateCard } from "@/components/dashboard/account-state-card";
import { AccountsCard } from "@/components/dashboard/accounts-card";
import { DashboardHeader } from "@/components/dashboard/dashboard-header";
import { InstrumentsCard } from "@/components/dashboard/instruments-card";
import { useAccountState, useAccounts } from "@/lib/accounts";
import { useLogout, useSession } from "@/lib/auth";
import { useInstruments } from "@/lib/instruments";

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

  const selectedInstrument = useMemo(
    () =>
      instrumentList.find((i) => i.tradableInstrumentId === selectedInstrumentId) ??
      instrumentList[0],
    [instrumentList, selectedInstrumentId],
  );

  useEffect(() => {
    if (instrumentList.length > 0 && selectedInstrumentId === null) {
      setSelectedInstrumentId(instrumentList[0].tradableInstrumentId);
    }
  }, [instrumentList, selectedInstrumentId]);

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
    </div>
  );
}
