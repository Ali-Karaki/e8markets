import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import type { Instrument } from "@/lib/api";

import { InstrumentsCard } from "./instruments-card";

const instruments: Instrument[] = [
  {
    id: 1,
    name: "EURUSD",
    type: "FOREX",
    tradableInstrumentId: 1001,
    tradingExchange: "FX",
    marketDataExchange: "FX",
    description: "Euro vs US Dollar",
  },
  {
    id: 2,
    name: "GBPUSD",
    type: "FOREX",
    tradableInstrumentId: 1002,
    tradingExchange: "FX",
    marketDataExchange: "FX",
    description: "British Pound vs US Dollar",
  },
];

describe("InstrumentsCard", () => {
  it("filters instruments by search", () => {
    render(
      <InstrumentsCard
        instruments={instruments}
        selectedInstrument={undefined}
        isLoading={false}
        error={null}
        onSelectInstrument={vi.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: /select a symbol/i }));
    fireEvent.change(screen.getByLabelText("Search instruments"), {
      target: { value: "EUR" },
    });

    expect(screen.getByText("EURUSD")).toBeTruthy();
    expect(screen.queryByText("GBPUSD")).toBeNull();
  });
});
