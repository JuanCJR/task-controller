const API_BASE = "/api/v1/task-controller";

function getToken(): string | null {
  const stored = localStorage.getItem("auth");
  if (!stored) return null;
  try {
    return JSON.parse(stored).token;
  } catch {
    return null;
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  if (res.status === 401) {
    localStorage.removeItem("auth");
    window.location.href = "/login";
    throw new Error("Unauthorized");
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || body.message || `Error ${res.status}`);
  }

  if (res.status === 204) return {} as T;
  return res.json();
}

// Auth
export const api = {
  login: (email: string, password: string) =>
    request<import("@/types").LoginResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),

  logout: () =>
    request<{ message: string }>("/auth/logout", { method: "POST" }),

  changePassword: (old_password: string, new_password: string) =>
    request<{ message: string }>("/auth/change-password", {
      method: "PUT",
      body: JSON.stringify({ old_password, new_password }),
    }),

  // Users
  getUsers: () => request<import("@/types").User[]>("/users"),

  createUser: (data: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    role_name: string;
  }) =>
    request<import("@/types").User>("/users", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  updateUser: (id: string, data: { email?: string; first_name?: string; last_name?: string }) =>
    request<{ message: string }>(`/users/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  deleteUser: (id: string) =>
    request<{ message: string }>(`/users/${id}`, { method: "DELETE" }),

  // Tasks
  getTasks: () => request<import("@/types").Task[]>("/tasks"),

  getTask: (id: string) => request<import("@/types").Task>(`/tasks/${id}`),

  createTask: (data: {
    title: string;
    description?: string;
    expires_at: string;
    assigned_to: string;
  }) =>
    request<import("@/types").Task>("/tasks", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  updateTask: (
    id: string,
    data: { title?: string; description?: string; expires_at?: string; assigned_to?: string }
  ) =>
    request<{ message: string }>(`/tasks/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  updateTaskState: (id: string, task_state: string) =>
    request<{ message: string }>(`/tasks/${id}/state`, {
      method: "PATCH",
      body: JSON.stringify({ task_state }),
    }),

  deleteTask: (id: string) =>
    request<{ message: string }>(`/tasks/${id}`, { method: "DELETE" }),

  // Comments
  getTaskComments: (taskId: string) =>
    request<import("@/types").TaskComment[]>(`/tasks/${taskId}/comments`),

  addTaskComment: (taskId: string, comment: string) =>
    request<import("@/types").TaskComment>(`/tasks/${taskId}/comments`, {
      method: "POST",
      body: JSON.stringify({ comment }),
    }),
};
