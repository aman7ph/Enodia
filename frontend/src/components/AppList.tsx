import { ScrollArea } from "@/components/ui/scroll-area";
import { AppCard } from "./AppCard";
import type { InstalledApp } from "@/types";

interface AppListProps {
  apps: InstalledApp[];
  blockedPaths: Set<string>;
  blockedPkgNames: Set<string>;
  selectedAppIds: Set<string>;
  onToggleSelect: (appId: string) => void;
}

export function AppList({
  apps,
  blockedPaths,
  blockedPkgNames,
  selectedAppIds,
  onToggleSelect,
}: AppListProps) {
  const isAppBlocked = (app: InstalledApp): boolean => {
    // For Store apps, check by name in PKG rules
    if (app.appType === "store") {
      return blockedPkgNames.has(app.name.toLowerCase());
    }
    // For Win32 apps, check by executable path
    return (
      app.executables?.some((exe) => blockedPaths.has(exe.toLowerCase())) ||
      false
    );
  };

  if (apps.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center text-muted-foreground">
        <p>No applications found</p>
      </div>
    );
  }

  return (
    <ScrollArea className="h-full w-full">
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 p-4">
        {apps.map((app) => (
          <AppCard
            key={app.id}
            app={app}
            isBlocked={isAppBlocked(app)}
            isSelected={selectedAppIds.has(app.id)}
            onToggleSelect={onToggleSelect}
          />
        ))}
      </div>
    </ScrollArea>
  );
}