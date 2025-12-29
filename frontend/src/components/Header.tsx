import { Button } from "@/components/ui/button";
import { useTheme } from "./ThemeProvider";

interface HeaderProps {
  onRefresh?: () => void;
  isRefreshing?: boolean;
}

export function Header({ onRefresh, isRefreshing }: HeaderProps) {
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme(theme === "dark" ? "light" : "dark");
  };

  return (
    <header className="flex items-center justify-between p-4 border-b">
      <div>
        <h1 className="text-2xl font-bold bg-gradient-to-r from-purple-500 to-blue-500 bg-clip-text text-transparent">
          Enodia
        </h1>
        <p className="text-sm text-muted-foreground">Network Controller</p>
      </div>

      <div className="flex items-center gap-2">
        {onRefresh && (
          <Button 
            variant="outline" 
            size="sm" 
            onClick={onRefresh}
            disabled={isRefreshing}
          >
            {isRefreshing ? "🔄 Refreshing..." : "🔄 Refresh"}
          </Button>
        )}
        <Button variant="ghost" size="sm" onClick={toggleTheme}>
          {theme === "dark" ? "☀️ Light" : "🌙 Dark"}
        </Button>
      </div>
    </header>
  );
}
