import { type UseQueryResult, useQuery } from "@tanstack/react-query";

import { api, type Instrument, type InstrumentsResponse } from "./api";

/** Substring match for default symbol selection (e.g. GBPCAD+). */
export const PREFERRED_INSTRUMENT_SYMBOL = "GBPCAD";

export function findPreferredInstrument(
  instruments: Instrument[],
): Instrument | undefined {
  if (instruments.length === 0) return undefined;
  const preferred = instruments.find((i) =>
    i.name.toUpperCase().includes(PREFERRED_INSTRUMENT_SYMBOL),
  );
  return preferred ?? instruments[0];
}

export function instrumentsQueryKey(accountId: string, accNum: string) {
  return ["instruments", accountId, accNum] as const;
}

export type InstrumentsQueryKey = ReturnType<typeof instrumentsQueryKey>;

export function instrumentsQueryOptions(accountId: string, accNum: string) {
  return {
    queryKey: instrumentsQueryKey(accountId, accNum),
    queryFn: (): Promise<InstrumentsResponse> =>
      api.instruments(accountId, accNum),
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
