import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({ component: Home });

function Home() {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center p-8">
      <h1 className="text-4xl font-bold">e8markets client</h1>
      <p className="mt-4 text-lg text-muted-foreground">
        Vite · React · TanStack Router · TanStack Query · shadcn/ui
      </p>
    </div>
  );
}
