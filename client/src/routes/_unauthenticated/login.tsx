import { createFileRoute, isRedirect, redirect } from "@tanstack/react-router";

import { LoginForm } from "@/components/login/login-form";
import { fetchSession } from "@/lib/auth";

type LoginSearch = {
  reason?: "expired";
};

export const Route = createFileRoute("/_unauthenticated/login")({
  validateSearch: (search: Record<string, unknown>): LoginSearch => ({
    reason: search.reason === "expired" ? "expired" : undefined,
  }),
  beforeLoad: async ({ context }) => {
    try {
      const session = await fetchSession(context.queryClient);
      if (session.authenticated) {
        throw redirect({ to: "/" });
      }
    } catch (err) {
      if (isRedirect(err)) {
        throw err;
      }
    }
  },
  component: LoginPage,
});

function LoginPage() {
  const { reason } = Route.useSearch();

  return (
    <div className="flex min-h-svh items-center justify-center p-8">
      <LoginForm sessionExpired={reason === "expired"} />
    </div>
  );
}
