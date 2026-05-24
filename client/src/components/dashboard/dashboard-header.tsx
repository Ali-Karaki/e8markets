import { Button } from "@/components/ui/button";

type DashboardHeaderProps = {
  email?: string;
  onLogout: () => void;
  isLoggingOut: boolean;
};

export function DashboardHeader({ email, onLogout, isLoggingOut }: DashboardHeaderProps) {
  return (
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">{email}</p>
      </div>
      <Button variant="outline" onClick={onLogout} disabled={isLoggingOut}>
        {isLoggingOut ? "Signing out..." : "Sign out"}
      </Button>
    </div>
  );
}
