import {
  useMutation,
  useQuery,
  useQueryClient,
  type QueryClient,
  type UseMutationResult,
  type UseQueryResult,
} from "@tanstack/react-query";
import { useEffect, useRef } from "react";

import {
  api,
  type Position,
  type PositionHistoryResponse,
  type PositionsResponse,
  type SyncPositionsOptions,
  type SyncPositionsResponse,
} from "./api";

export const POSITIONS_POLL_INTERVAL_MS = 180_000;

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

function applySyncResult(
  queryClient: QueryClient,
  accountId: string,
  accNum: string,
  data: SyncPositionsResponse,
) {
  if (data.positions) {
    queryClient.setQueryData(positionsQueryKey(accountId, accNum), {
      positions: data.positions,
    });
  }
  if (data.recordsStored > 0) {
    queryClient.invalidateQueries({
      queryKey: ["positions", "history", accountId, accNum],
    });
  }
}

type SyncVariables = {
  accountId: string;
  accNum: string;
  options?: SyncPositionsOptions;
};

export function useSyncPositions(
  onSuccessMessage?: (data: SyncPositionsResponse) => void,
): UseMutationResult<SyncPositionsResponse, Error, SyncVariables> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ accountId, accNum, options }: SyncVariables) =>
      api.syncPositions(accountId, accNum, options),
    onSuccess: (data, { accountId, accNum }) => {
      applySyncResult(queryClient, accountId, accNum, data);
      onSuccessMessage?.(data);
    },
  });
}

export function formatSyncMessage(data: SyncPositionsResponse, manual: boolean): string {
  if (data.skipped) {
    return "Auto-sync skipped (no new open positions).";
  }
  const prefix = manual ? "Saved" : "Auto-synced";
  return `${prefix} ${data.recordsStored} snapshot${data.recordsStored === 1 ? "" : "s"} to the database (sync run #${data.syncRunId}).`;
}

export function useAutoPositionSync(
  accountId: string | undefined,
  accNum: string | undefined,
  enabled: boolean,
  onMessage?: (message: string) => void,
) {
  const queryClient = useQueryClient();
  const id = accountId ?? "";
  const num = accNum ?? "";
  const onMessageRef = useRef(onMessage);
  onMessageRef.current = onMessage;

  useEffect(() => {
    if (!enabled || !accountId || !accNum) return;

    let intervalId: ReturnType<typeof setInterval> | undefined;

    const runAutoSync = () => {
      if (document.visibilityState === "hidden") return;

      api
        .syncPositions(id, num, { onlyIfNew: true })
        .then((data) => {
          applySyncResult(queryClient, id, num, data);
          onMessageRef.current?.(formatSyncMessage(data, false));
        })
        .catch(() => {
          // Keep last good cache; do not spam errors on rate limits.
        });
    };

    const startInterval = () => {
      if (intervalId !== undefined) return;
      intervalId = setInterval(runAutoSync, POSITIONS_POLL_INTERVAL_MS);
    };

    const stopInterval = () => {
      if (intervalId === undefined) return;
      clearInterval(intervalId);
      intervalId = undefined;
    };

    const onVisibilityChange = () => {
      if (document.visibilityState === "hidden") {
        stopInterval();
      } else {
        startInterval();
      }
    };

    startInterval();
    document.addEventListener("visibilitychange", onVisibilityChange);

    return () => {
      stopInterval();
      document.removeEventListener("visibilitychange", onVisibilityChange);
    };
  }, [enabled, accountId, accNum, id, num, queryClient]);
}
