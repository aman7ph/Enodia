import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

interface ToolbarProps {
  searchTerm: string;
  onSearchChange: (term: string) => void;
  selectedCount: number;
  onBlock: () => void;
  onUnblock: () => void;
  activeTab: string;
}

export function Toolbar({
  searchTerm,
  onSearchChange,
  selectedCount,
  onBlock,
  onUnblock,
  activeTab,
}: ToolbarProps) {
  const showBlock = activeTab !== "blocked";
  const showUnblock = activeTab !== "unblocked";

  return (
    <div className="flex items-center gap-3 p-4 border-b">
      <Input
        type="text"
        placeholder="Search applications..."
        value={searchTerm}
        onChange={(e) => onSearchChange(e.target.value)}
        className="flex-1"
      />
      {showBlock && (
        <Button
          variant="destructive"
          onClick={onBlock}
          disabled={selectedCount === 0}
        >
          🚫 Block ({selectedCount})
        </Button>
      )}
      {showUnblock && (
        <Button
          variant="outline"
          onClick={onUnblock}
          disabled={selectedCount === 0}
        >
          ✅ Unblock ({selectedCount})
        </Button>
      )}
    </div>
  );
}

