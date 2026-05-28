import { type UseQueryResult, useQuery } from "@tanstack/react-query";

import { type AccountState, type AccountsResponse, api } from "./api";

export const accountsQueryKey = ["accounts"] as const;

export type AccountsQueryKey = typeof accountsQueryKey;

export function accountsQueryOptions() {
  return {
    queryKey: accountsQueryKey,
    queryFn: (): Promise<AccountsResponse> => api.accounts(),
  };
}

export function useAccounts(): UseQueryResult<AccountsResponse, Error> {
  return useQuery(accountsQueryOptions());
}

export function accountStateQueryKey(accountId: string, accNum: string) {
  return ["accounts", "state", accountId, accNum] as const;
}

export type AccountStateQueryKey = ReturnType<typeof accountStateQueryKey>;

export function accountStateQueryOptions(accountId: string, accNum: string) {
  return {
    queryKey: accountStateQueryKey(accountId, accNum),
    queryFn: (): Promise<AccountState> => api.accountState(accountId, accNum),
  };
}

export function useAccountState(
  accountId: string | undefined,
  accNum: string | undefined,
): UseQueryResult<AccountState, Error> {
  const id = accountId ?? "";
  const num = accNum ?? "";

  return useQuery({
    ...accountStateQueryOptions(id, num),
    enabled: Boolean(accountId && accNum),
  });
}
