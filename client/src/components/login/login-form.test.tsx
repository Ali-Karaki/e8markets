import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ApiError } from "@/lib/api";

const mutateAsync = vi.fn();
const navigate = vi.fn();

vi.mock("@/lib/auth", () => ({
  useLogin: () => ({
    mutateAsync,
    isPending: false,
  }),
}));

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => navigate,
}));

import { LoginForm } from "./login-form";

describe("LoginForm", () => {
  beforeEach(() => {
    mutateAsync.mockReset();
    navigate.mockReset();
  });

  it("submits credentials", async () => {
    mutateAsync.mockResolvedValue({ authenticated: true });

    render(<LoginForm />);

    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "user@example.com" },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "secret" },
    });
    fireEvent.submit(screen.getByRole("button", { name: "Sign in" }).closest("form")!);

    await waitFor(() => {
      expect(mutateAsync).toHaveBeenCalledWith({
        email: "user@example.com",
        password: "secret",
        server: "E8",
      });
    });
    expect(navigate).toHaveBeenCalledWith({ to: "/" });
  });

  it("shows login error", async () => {
    mutateAsync.mockRejectedValue(
      new ApiError(401, "Invalid credentials", "invalid_credentials"),
    );

    render(<LoginForm />);

    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "user@example.com" },
    });
    fireEvent.change(screen.getByLabelText("Password"), {
      target: { value: "wrong" },
    });
    fireEvent.submit(screen.getByRole("button", { name: "Sign in" }).closest("form")!);

    expect(await screen.findByText("Invalid credentials")).toBeTruthy();
    expect(navigate).not.toHaveBeenCalled();
  });

  it("shows session expired banner", () => {
    render(<LoginForm sessionExpired />);

    expect(screen.getByText("Your session expired. Please sign in again.")).toBeTruthy();
  });
});
