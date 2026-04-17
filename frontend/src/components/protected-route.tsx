import { Navigate, Outlet } from "react-router";
import { useAuth } from "@/context/auth";

export function ProtectedRoute() {
  const { isAuthenticated, user } = useAuth();

  if (!isAuthenticated) return <Navigate to="/login" replace />;

  if (user?.must_change_password) return <Navigate to="/change-password" replace />;

  return <Outlet />;
}

export function GuestRoute() {
  const { isAuthenticated, user } = useAuth();

  if (isAuthenticated && !user?.must_change_password) return <Navigate to="/tasks" replace />;

  return <Outlet />;
}

export function ChangePasswordRoute() {
  const { isAuthenticated, user } = useAuth();

  if (!isAuthenticated) return <Navigate to="/login" replace />;
  if (!user?.must_change_password) return <Navigate to="/tasks" replace />;

  return <Outlet />;
}
