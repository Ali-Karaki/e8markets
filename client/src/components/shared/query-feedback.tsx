import { apiErrorMessage } from "@/lib/errors";

type QueryFeedbackProps = {
  isLoading?: boolean;
  error?: unknown;
  isEmpty?: boolean;
  loadingMessage?: string;
  errorMessage?: string;
  emptyMessage?: string;
  children?: React.ReactNode;
};

export function QueryFeedback({
  isLoading,
  error,
  isEmpty,
  loadingMessage = "Loading...",
  errorMessage = "Something went wrong",
  emptyMessage,
  children,
}: QueryFeedbackProps) {
  if (isLoading) {
    return <p className="text-sm text-muted-foreground">{loadingMessage}</p>;
  }

  if (error) {
    return (
      <p className="text-sm text-destructive">{apiErrorMessage(error, errorMessage)}</p>
    );
  }

  if (isEmpty && emptyMessage) {
    return <p className="text-sm text-muted-foreground">{emptyMessage}</p>;
  }

  return children ?? null;
}
