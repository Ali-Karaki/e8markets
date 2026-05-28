const API_URL = import.meta.env.VITE_API_URL;

const apiPaths = {
  login: "/api/auth/login",
  session: "/api/auth/session",
  logout: "/api/auth/logout",
  accounts: "/api/accounts",
  accountState: "/api/accounts/state",
  instruments: "/api/instruments",
  positions: "/api/positions",
  positionsSync: "/api/positions/sync",
  positionsHistory: "/api/positions/history",
} as const;

type ApiPath = (typeof apiPaths)[keyof typeof apiPaths];
type ApiRequestPath = ApiPath | `${ApiPath}?${string}`;

export type LoginRequest = {
  email: string;
  password: string;
  server?: string;
};

export type Session = {
  authenticated: boolean;
  email: string;
  server: string;
  expiresAt: string;
};

export type LogoutResponse = {
  ok: boolean;
};

export type Account = {
  id: string;
  name: string;
  currency: string;
  status: string;
  accNum: string;
  accountBalance?: number;
};

export type AccountsResponse = {
  accounts: Account[];
};

export type AccountState = {
  accountId: string;
  accNum: string;
  state: Record<string, number>;
};

export type Instrument = {
  id: number;
  tradableInstrumentId: number;
  name: string;
  description?: string;
  type: string;
  tradingExchange: string;
  marketDataExchange: string;
  tradeRouteId?: number;
  infoRouteId?: number;
};

export type InstrumentsResponse = {
  instruments: Instrument[];
};

export type Position = {
  id: string;
  tradableInstrumentId?: number;
  side?: string;
  qty?: number;
  avgPrice?: number;
  routeId?: number;
  stopLossId?: number;
  takeProfitId?: number;
  fields?: Record<string, string>;
};

export type PositionsResponse = {
  positions: Position[];
};

export type SyncPositionsOptions = {
  force?: boolean;
  onlyIfNew?: boolean;
};

export type SyncPositionsResponse = {
  syncRunId: number;
  status: string;
  recordsStored: number;
  skipped?: boolean;
  positions: Position[];
};

export type PositionSnapshot = {
  id: number;
  syncRunId: number;
  syncedAt: string;
  positionId: string;
  position: Position;
};

export type PositionHistoryResponse = {
  snapshots: PositionSnapshot[];
};

type ApiErrorBody = {
  error: string;
  code?: string;
};

export class ApiError extends Error {
  status: number;
  code?: string;

  constructor(status: number, message: string, code?: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

export type Api = {
  login(body: LoginRequest): Promise<Session>;
  session(): Promise<Session>;
  logout(): Promise<LogoutResponse>;
  accounts(): Promise<AccountsResponse>;
  accountState(accountId: string, accNum: string): Promise<AccountState>;
  instruments(accountId: string, accNum: string): Promise<InstrumentsResponse>;
  positions(
    accountId: string,
    accNum: string,
    tradableInstrumentId?: number,
  ): Promise<PositionsResponse>;
  syncPositions(
    accountId: string,
    accNum: string,
    options?: SyncPositionsOptions,
  ): Promise<SyncPositionsResponse>;
  positionHistory(
    accountId: string,
    accNum: string,
    tradableInstrumentId?: number,
  ): Promise<PositionHistoryResponse>;
};

function isApiErrorBody(data: unknown): data is ApiErrorBody {
  return (
    typeof data === "object" &&
    data !== null &&
    "error" in data &&
    typeof data.error === "string"
  );
}

async function request<T>(
  path: ApiRequestPath,
  init?: RequestInit,
): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });

  const data: unknown = await res.json().catch(() => null);

  if (!res.ok) {
    const message = isApiErrorBody(data) ? data.error : "Request failed";
    const code = isApiErrorBody(data) ? data.code : undefined;

    if (res.status === 401 && path !== apiPaths.login) {
      if (window.location.pathname !== "/login") {
        window.location.assign("/login?reason=expired");
      }
    }

    throw new ApiError(res.status, message, code);
  }

  return data as T;
}

export const api: Api = {
  login(body) {
    return request<Session>(apiPaths.login, {
      method: "POST",
      body: JSON.stringify(body),
    });
  },

  session() {
    return request<Session>(apiPaths.session);
  },

  logout() {
    return request<LogoutResponse>(apiPaths.logout, { method: "POST" });
  },

  accounts() {
    return request<AccountsResponse>(apiPaths.accounts);
  },

  accountState(accountId, accNum) {
    const params = new URLSearchParams({ accountId, accNum });
    return request<AccountState>(`${apiPaths.accountState}?${params}`);
  },

  instruments(accountId, accNum) {
    const params = new URLSearchParams({ accountId, accNum });
    return request<InstrumentsResponse>(`${apiPaths.instruments}?${params}`);
  },

  positions(accountId, accNum, tradableInstrumentId) {
    const params = new URLSearchParams({ accountId, accNum });
    if (tradableInstrumentId !== undefined) {
      params.set("tradableInstrumentId", String(tradableInstrumentId));
    }
    return request<PositionsResponse>(`${apiPaths.positions}?${params}`);
  },

  syncPositions(accountId, accNum, options) {
    const params = new URLSearchParams({ accountId, accNum });
    if (options?.force) params.set("force", "true");
    if (options?.onlyIfNew) params.set("onlyIfNew", "true");
    return request<SyncPositionsResponse>(
      `${apiPaths.positionsSync}?${params}`,
      {
        method: "POST",
      },
    );
  },

  positionHistory(accountId, accNum, tradableInstrumentId) {
    const params = new URLSearchParams({ accountId, accNum });
    if (tradableInstrumentId !== undefined) {
      params.set("tradableInstrumentId", String(tradableInstrumentId));
    }
    return request<PositionHistoryResponse>(
      `${apiPaths.positionsHistory}?${params}`,
    );
  },
};
