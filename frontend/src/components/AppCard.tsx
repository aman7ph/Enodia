import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import type { InstalledApp } from "@/types";

interface AppCardProps {
  app: InstalledApp;
  isBlocked: boolean;
  isSelected: boolean;
  onToggleSelect: (appId: string) => void;
}

export function AppCard({
  app,
  isBlocked,
  isSelected,
  onToggleSelect,
}: AppCardProps) {
  return (
    <Card
      className={`
        relative cursor-pointer transition-all duration-200
        hover:bg-accent/50 hover:scale-[1.02]
        ${isSelected ? "ring-2 ring-primary bg-accent" : ""}
        ${isBlocked ? "border-destructive/50 bg-destructive/5" : ""}
      `}
      onClick={() => onToggleSelect(app.id)}
    >
      {/* Checkbox in top-left corner */}
      <div className="absolute top-2 left-2">
        <Checkbox
          checked={isSelected}
          onClick={(e) => e.stopPropagation()}
          onCheckedChange={() => onToggleSelect(app.id)}
        />
      </div>

      {/* Blocked badge in top-right corner */}
      {isBlocked && (
        <div className="absolute top-2 right-2">
          <Badge variant="destructive" className="text-xs">
            Blocked
          </Badge>
        </div>
      )}

      <CardContent className="flex flex-col items-center p-4 pt-8">
        {/* Large centered icon */}
        <div className="w-16 h-16 flex items-center justify-center rounded-xl bg-muted mb-3">
          {app.iconBase64 ? (
            <img
              src={`data:image/png;base64,${app.iconBase64}`}
              alt=""
              className="w-12 h-12 object-contain"
            />
          ) : (
            <span className="text-3xl">📦</span>
          )}
        </div>

        {/* App name - centered, 2 lines max */}
        <h3 className="font-medium text-center text-sm line-clamp-2 mb-1">
          {app.name}
        </h3>

        {/* Publisher - subtle */}
        <p className="text-xs text-muted-foreground text-center truncate w-full">
          {app.publisher || "Unknown"}
        </p>

        {/* Executable count at bottom */}
        <Badge variant="outline" className="text-xs mt-2">
          {app.executables?.length || 0} exe
        </Badge>
      </CardContent>
    </Card>
  );
}
