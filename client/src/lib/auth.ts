import {
  useMutation,
  useQuery,
  useQueryClient,
  type QueryClient,
  type UseMutationResult,
  type UseQueryResult,
} from "@tanstack/react-query";

import { api, type LoginRequest, type LogoutResponse, type Session } from "./api";

export const sessionQueryKey = ["auth", "session"] as const;

export type SessionQueryKey = typeof sessionQueryKey;

export function sessionQueryOptions() {
  return {
    queryKey: sessionQueryKey,
    queryFn: (): Promise<Session> => api.session(),
  };
}

export function fetchSession(queryClient: QueryClient): Promise<Session> {
  return queryClient.fetchQuery(sessionQueryOptions());
}

export function useSession(): UseQueryResult<Session, Error> {
  return useQuery(sessionQueryOptions());
}

export function useLogin(): UseMutationResult<Session, Error, LoginRequest> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: LoginRequest) => api.login(body),
    onSuccess: (data: Session) => {
      queryClient.setQueryData(sessionQueryKey, data);
    },
  });
}

export function useLogout(): UseMutationResult<LogoutResponse, Error, void> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (): Promise<LogoutResponse> => api.logout(),
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: sessionQueryKey });
      queryClient.removeQueries({ queryKey: ["accounts"] });
      queryClient.removeQueries({ queryKey: ["positions"] });
      queryClient.removeQueries({ queryKey: ["instruments"] });
    },
  });
}
