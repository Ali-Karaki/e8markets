import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationResult,
  type UseQueryResult,
} from "@tanstack/react-query";

import {
  api,
  type Position,
  type PositionHistoryResponse,
  type PositionsResponse,
  type SyncPositionsResponse,
} from "./api";

export function positionsQueryKey(accountId: string, accNum: string) {
  return ["positions", accountId, accNum] as const;
}

export function filterPositionsByInstrument(
  positions: Position[],
  tradableInstrumentId?: number,
): Position[] {
  if (tradableInstrumentId === undefined) {
    return positions;
  }
  return positions.filter((p) => p.tradableInstrumentId === tradableInstrumentId);
}

export function positionHistoryQueryKey(
  accountId: string,
  accNum: string,
  tradableInstrumentId?: number,
) {
  return ["positions", "history", accountId, accNum, tradableInstrumentId ?? "all"] as const;
}

export function usePositions(
  accountId: string | undefined,
  accNum: string | undefined,
): UseQueryResult<PositionsResponse, Error> {
  const id = accountId ?? "";
  const num = accNum ?? "";

  return useQuery({
    queryKey: positionsQueryKey(id, num),
    queryFn: (): Promise<PositionsResponse> => api.positions(id, num),
    enabled: Boolean(accountId && accNum),
  });
}

export function usePositionHistory(
  accountId: string | undefined,
  accNum: string | undefined,
  tradableInstrumentId?: number,
): UseQueryResult<PositionHistoryResponse, Error> {
  const id = accountId ?? "";
  const num = accNum ?? "";

  return useQuery({
    queryKey: positionHistoryQueryKey(id, num, tradableInstrumentId),
    queryFn: (): Promise<PositionHistoryResponse> =>
      api.positionHistory(id, num, tradableInstrumentId),
    enabled: Boolean(accountId && accNum),
  });
}

type SyncVariables = {
  accountId: string;
  accNum: string;
};

export function useSyncPositions(): UseMutationResult<
  SyncPositionsResponse,
  Error,
  SyncVariables
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ accountId, accNum }: SyncVariables) =>
      api.syncPositions(accountId, accNum),
    onSuccess: (_data, { accountId, accNum }) => {
      queryClient.invalidateQueries({ queryKey: ["positions", accountId, accNum] });
      queryClient.invalidateQueries({
        queryKey: ["positions", "history", accountId, accNum],
      });
    },
  });
}
