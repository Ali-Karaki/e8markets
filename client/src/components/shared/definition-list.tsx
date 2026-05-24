import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

export type DefinitionListItem = {
  label: string;
  value: ReactNode;
};

type DefinitionListProps = {
  items: DefinitionListItem[];
  columns?: 2 | 3;
  className?: string;
};

export function DefinitionList({ items, columns = 2, className }: DefinitionListProps) {
  if (items.length === 0) return null;

  return (
    <dl
      className={cn(
        "grid gap-3 text-sm",
        columns === 3 ? "grid-cols-2 md:grid-cols-3" : "grid-cols-2",
        className,
      )}
    >
      {items.map((item) => (
        <div key={item.label}>
          <dt className="text-muted-foreground">{item.label}</dt>
          <dd className="font-medium">{item.value}</dd>
        </div>
      ))}
    </dl>
  );
}
