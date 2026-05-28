import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { AccountsCard } from "./accounts-card";

describe("AccountsCard", () => {
  it("shows empty state when no accounts", () => {
    render(
      <AccountsCard
        accounts={[]}
        selectedAccount={undefined}
        isLoading={false}
        error={null}
        onSelectAccount={vi.fn()}
      />,
    );

    expect(screen.getByText("No trading accounts were returned for this login.")).toBeTruthy();
    expect(screen.getByText("No accounts available")).toBeTruthy();
  });
});
