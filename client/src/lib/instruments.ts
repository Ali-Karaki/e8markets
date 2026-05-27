import { useQuery, type UseQueryResult } from "@tanstack/react-query";

import { api, type InstrumentsResponse } from "./api";

export function instrumentsQueryKey(accountId: string, accNum: string) {
  return ["instruments", accountId, accNum] as const;
}

export type InstrumentsQueryKey = ReturnType<typeof instrumentsQueryKey>;

export function instrumentsQueryOptions(accountId: string, accNum: string) {
  return {
    queryKey: instrumentsQueryKey(accountId, accNum),
    queryFn: (): Promise<InstrumentsResponse> => api.instruments(accountId, accNum),
  };
}

export function useInstruments(
  accountId: string | undefined,
  accNum: string | undefined,
): UseQueryResult<InstrumentsResponse, Error> {
  const id = accountId ?? "";
  const num = accNum ?? "";

  return useQuery({
    ...instrumentsQueryOptions(id, num),
    enabled: Boolean(accountId && accNum),
  });
}
