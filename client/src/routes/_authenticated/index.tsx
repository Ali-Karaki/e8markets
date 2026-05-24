import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useMemo, useState } from "react";

import { AccountStateCard } from "@/components/dashboard/account-state-card";
import { AccountsCard } from "@/components/dashboard/accounts-card";
import { DashboardHeader } from "@/components/dashboard/dashboard-header";
import { useAccountState, useAccounts } from "@/lib/accounts";
import { useLogout, useSession } from "@/lib/auth";

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

  const selectedAccount = useMemo(
    () => accountList.find((a) => a.id === selectedId) ?? accountList[0],
    [accountList, selectedId],
  );

  useEffect(() => {
    if (accountList.length > 0 && !selectedId) {
      setSelectedId(accountList[0].id);
    }
  }, [accountList, selectedId]);

  const accountState = useAccountState(selectedAccount?.id, selectedAccount?.accNum);

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
      <AccountStateCard
        state={accountState.data?.state ?? {}}
        isLoading={accountState.isLoading}
        error={accountState.error}
      />
    </div>
  );
}
