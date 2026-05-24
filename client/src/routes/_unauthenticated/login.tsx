import { createFileRoute, isRedirect, redirect } from "@tanstack/react-router";

import { LoginForm } from "@/components/login/login-form";
import { fetchSession } from "@/lib/auth";

export const Route = createFileRoute("/_unauthenticated/login")({
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
  return (
    <div className="flex min-h-svh items-center justify-center p-8">
      <LoginForm />
    </div>
  );
}
