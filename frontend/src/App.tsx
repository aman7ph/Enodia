import { useState, useEffect } from "react";
import { Header } from "@/components/Header";
import { Toolbar } from "@/components/Toolbar";
import { AppList } from "@/components/AppList";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { InstalledApp, BlockedApp } from "@/types";

import {
  GetInstalledApps,
  GetBlockedApps,
  BlockInstalledApp,
  UnblockInstalledApp,
} from "../wailsjs/go/main/App";

function App() {
  const [apps, setApps] = useState<InstalledApp[]>([]);
  const [blockedPaths, setBlockedPaths] = useState<Set<string>>(new Set());
  const [blockedPkgNames, setBlockedPkgNames] = useState<Set<string>>(new Set());
  const [selectedAppIds, setSelectedAppIds] = useState<Set<string>>(new Set());
  const [searchTerm, setSearchTerm] = useState("");
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("all");

  useEffect(() => {
    // Retry a few times on startup in case backend isn't ready
    const tryLoad = async (retries = 3) => {
      await loadData();
      if (apps.length === 0 && retries > 0) {
        setTimeout(() => tryLoad(retries - 1), 500);
      }
    };
    tryLoad();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [installedApps, blockedApps] = await Promise.all([
        GetInstalledApps(),
        GetBlockedApps(),
      ]);

      setApps(installedApps || []);

      // Build sets for both path-based and PKG-based rules
      const blockedPathSet = new Set<string>();
      const blockedPkgSet = new Set<string>();
      
      (blockedApps || []).forEach((app: BlockedApp) => {
        const name = app.appPath || app.displayName;
        if (name.includes("PKG-")) {
          // Extract the display name from PKG rules
          blockedPkgSet.add(name.replace("PKG-", "").toLowerCase());
        } else {
          blockedPathSet.add(name.toLowerCase());
        }
      });
      
      setBlockedPaths(blockedPathSet);
      setBlockedPkgNames(blockedPkgSet);
    } catch (err) {
      console.error("Failed to load apps:", err);
    }
    setLoading(false);
  };

  // Check if an app is blocked (handles both Win32 and Store apps)
  const isAppBlocked = (app: InstalledApp): boolean => {
    // For Store apps, check if name matches blocked PKG rules
    if (app.appType === "store") {
      return blockedPkgNames.has(app.name.toLowerCase());
    }
    // For Win32 apps, check if any executable is blocked
    return app.executables?.some(exe => 
      blockedPaths.has(exe.toLowerCase())
    ) || false;
  };

  // Filter apps by search term
  const filteredApps = apps.filter(
    (app) =>
      app.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      app.publisher?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  // Separate blocked and unblocked apps
  const blockedApps = filteredApps.filter(isAppBlocked);
  const unblockedApps = filteredApps.filter(app => !isAppBlocked(app));

  const handleToggleSelect = (appId: string) => {
    setSelectedAppIds((prev) => {
      const next = new Set(prev);
      if (next.has(appId)) {
        next.delete(appId);
      } else {
        next.add(appId);
      }
      return next;
    });
  };

  const handleBlock = async () => {
    // Get selected apps from current view
    const appsToBlock = apps.filter(app => selectedAppIds.has(app.id));
    
    for (const app of appsToBlock) {
      await BlockInstalledApp(app);
    }

    await loadData();
    setSelectedAppIds(new Set());
  };

  const handleUnblock = async () => {
    // Get selected apps
    const appsToUnblock = apps.filter(app => selectedAppIds.has(app.id));
    
    for (const app of appsToUnblock) {
      await UnblockInstalledApp(app);
    }

    await loadData();
    setSelectedAppIds(new Set());
  };

  if (loading) {
    return (
      <div className="h-screen flex items-center justify-center bg-background">
        <p className="text-muted-foreground">Loading applications...</p>
      </div>
    );
  }

  return (
    <div className="h-screen flex flex-col overflow-hidden bg-background text-foreground">
      <Header onRefresh={loadData} isRefreshing={loading} />

      <Toolbar
        searchTerm={searchTerm}
        onSearchChange={setSearchTerm}
        selectedCount={selectedAppIds.size}
        onBlock={handleBlock}
        onUnblock={handleUnblock}
        activeTab={activeTab}
      />

      <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 min-h-0 flex flex-col">
        <div className="px-4 border-b">
          <TabsList>
            <TabsTrigger value="all">
              All Apps ({filteredApps.length})
            </TabsTrigger>
            <TabsTrigger value="blocked">
              🚫 Blocked ({blockedApps.length})
            </TabsTrigger>
            <TabsTrigger value="unblocked">
              ✅ Allowed ({unblockedApps.length})
            </TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="all" className="flex-1 min-h-0 overflow-hidden m-0">
          <AppList
            apps={filteredApps}
            blockedPaths={blockedPaths}
            blockedPkgNames={blockedPkgNames}
            selectedAppIds={selectedAppIds}
            onToggleSelect={handleToggleSelect}
          />
        </TabsContent>

        <TabsContent value="blocked" className="flex-1 min-h-0 overflow-hidden m-0">
          <AppList
            apps={blockedApps}
            blockedPaths={blockedPaths}
            blockedPkgNames={blockedPkgNames}
            selectedAppIds={selectedAppIds}
            onToggleSelect={handleToggleSelect}
          />
        </TabsContent>

        <TabsContent value="unblocked" className="flex-1 min-h-0 overflow-hidden m-0">
          <AppList
            apps={unblockedApps}
            blockedPaths={blockedPaths}
            blockedPkgNames={blockedPkgNames}
            selectedAppIds={selectedAppIds}
            onToggleSelect={handleToggleSelect}
          />
        </TabsContent>
      </Tabs>

      <footer className="p-4 border-t text-center text-sm text-muted-foreground">
        {apps.length} apps discovered • {blockedApps.length} blocked
      </footer>
    </div>
  );
}

export default App;
