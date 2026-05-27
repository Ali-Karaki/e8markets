const API_URL = import.meta.env.VITE_API_URL;

const apiPaths = {
  login: "/api/auth/login",
  session: "/api/auth/session",
  logout: "/api/auth/logout",
  accounts: "/api/accounts",
  accountState: "/api/accounts/state",
  instruments: "/api/instruments",
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

type ApiErrorBody = {
  error: string;
};

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export type Api = {
  login(body: LoginRequest): Promise<Session>;
  session(): Promise<Session>;
  logout(): Promise<LogoutResponse>;
  accounts(): Promise<AccountsResponse>;
  accountState(accountId: string, accNum: string): Promise<AccountState>;
  instruments(accountId: string, accNum: string): Promise<InstrumentsResponse>;
};

function isApiErrorBody(data: unknown): data is ApiErrorBody {
  return (
    typeof data === "object" &&
    data !== null &&
    "error" in data &&
    typeof data.error === "string"
  );
}

async function request<T>(path: ApiRequestPath, init?: RequestInit): Promise<T> {
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
    if (res.status === 401 && path !== apiPaths.login) {
      if (window.location.pathname !== "/login") {
        window.location.assign("/login");
      }
    }

    const message = isApiErrorBody(data) ? data.error : "Request failed";
    throw new ApiError(res.status, message);
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
};
