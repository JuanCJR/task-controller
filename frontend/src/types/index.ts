export interface User {
  id: string;
  email: string;
  must_change_password: boolean;
  first_name: string;
  last_name: string;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  expires_at: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  assigned_to: string;
  task_state: string;
}

export interface TaskComment {
  id: string;
  task_id: string;
  user_id: string;
  comment: string;
  created_at: string;
}

export interface LoginResponse {
  token: string;
  must_change_password: boolean;
  user_id: string;
  email: string;
  first_name: string;
  last_name: string;
}

export interface AuthUser {
  token: string;
  user_id: string;
  email: string;
  first_name: string;
  last_name: string;
  must_change_password: boolean;
}

export type TaskState =
  | "Pendiente"
  | "Asignado"
  | "Iniciado"
  | "En espera"
  | "Finalizada con exito"
  | "Finalizada con error";
