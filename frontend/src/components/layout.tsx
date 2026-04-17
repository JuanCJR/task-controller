import { Outlet, Link, useNavigate } from "react-router";
import { useAuth } from "@/context/auth";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";

export function Layout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await api.logout();
    } catch {
      // ignore
    }
    logout();
    navigate("/login");
  };

  // Determine role from user permissions - we'll store it after first task fetch
  // For now we show all nav items and the backend handles permissions
  return (
    <div className="min-h-screen">
      <header className="border-b px-4 py-3 flex items-center justify-between">
        <nav className="flex items-center gap-4">
          <Link to="/tasks" className="font-bold text-lg">
            Task Controller
          </Link>
          <Link to="/tasks">
            <Button variant="ghost" size="sm">Tareas</Button>
          </Link>
          <Link to="/users">
            <Button variant="ghost" size="sm">Usuarios</Button>
          </Link>
        </nav>
        <div className="flex items-center gap-3">
          <span className="text-sm text-muted-foreground">
            {user?.first_name} {user?.last_name}
          </span>
          <Button variant="outline" size="sm" onClick={handleLogout}>
            Cerrar sesion
          </Button>
        </div>
      </header>
      <main className="p-4">
        <Outlet />
      </main>
    </div>
  );
}
