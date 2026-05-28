import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useForm } from "react-hook-form";

import { FormField } from "@/components/shared/form-field";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Login_FormSchema, type Login_FormType } from "@/form-schemas/login";
import { useLogin } from "@/lib/auth";
import { apiErrorMessage } from "@/lib/errors";

type LoginFormProps = {
  sessionExpired?: boolean;
};

export function LoginForm({ sessionExpired = false }: LoginFormProps) {
  const navigate = useNavigate();
  const login = useLogin();
  const [error, setError] = useState<string | null>(null);

  const form = useForm<Login_FormType>({
    resolver: zodResolver(Login_FormSchema),
    defaultValues: {
      email: "",
      password: "",
      server: "E8",
    },
  });

  async function onSubmit(values: Login_FormType) {
    setError(null);

    try {
      await login.mutateAsync(values);
      navigate({ to: "/" });
    } catch (err) {
      setError(apiErrorMessage(err, "Login failed"));
    }
  }

  const isSubmitting = form.formState.isSubmitting || login.isPending;

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Sign in</CardTitle>
        <CardDescription>Log in with your TradeLocker credentials</CardDescription>
      </CardHeader>

      <CardContent>
        <form className="space-y-4" noValidate onSubmit={form.handleSubmit(onSubmit)}>
          <FormField label="Email" htmlFor="email" error={form.formState.errors.email?.message}>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              aria-invalid={Boolean(form.formState.errors.email)}
              {...form.register("email")}
            />
          </FormField>

          <FormField label="Password" htmlFor="password" error={form.formState.errors.password?.message}>
            <Input
              id="password"
              type="password"
              autoComplete="current-password"
              aria-invalid={Boolean(form.formState.errors.password)}
              {...form.register("password")}
            />
          </FormField>

          <FormField label="Server" htmlFor="server" error={form.formState.errors.server?.message}>
            <Input
              id="server"
              aria-invalid={Boolean(form.formState.errors.server)}
              {...form.register("server")}
            />
          </FormField>

          {sessionExpired ? (
            <p className="text-sm text-muted-foreground">
              Your session expired. Please sign in again.
            </p>
          ) : null}

          {error ? <p className="text-sm text-destructive">{error}</p> : null}

          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? "Signing in..." : "Sign in"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
