import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { api } from "@/lib/api";
import type { Task, TaskComment, User, TaskState } from "@/types";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from "@/components/ui/select";

const TASK_STATES: TaskState[] = [
  "Pendiente", "Asignado", "Iniciado", "En espera", "Finalizada con exito", "Finalizada con error",
];

function isExpired(expiresAt: string): boolean {
  return new Date(expiresAt) < new Date();
}

export default function TaskDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [task, setTask] = useState<Task | null>(null);
  const [comments, setComments] = useState<TaskComment[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // State change
  const [newState, setNewState] = useState("");
  const [changingState, setChangingState] = useState(false);

  // Comment
  const [comment, setComment] = useState("");
  const [addingComment, setAddingComment] = useState(false);

  // Edit
  const [editing, setEditing] = useState(false);
  const [editTitle, setEditTitle] = useState("");
  const [editDescription, setEditDescription] = useState("");
  const [editExpiresAt, setEditExpiresAt] = useState("");
  const [editAssignedTo, setEditAssignedTo] = useState("");
  const [saving, setSaving] = useState(false);

  const fetchData = async () => {
    if (!id) return;
    try {
      const [t, c] = await Promise.all([api.getTask(id), api.getTaskComments(id)]);
      setTask(t);
      setComments(c || []);
      setNewState(t.task_state);
      setEditTitle(t.title);
      setEditDescription(t.description || "");
      setEditExpiresAt(t.expires_at ? new Date(t.expires_at).toISOString().slice(0, 16) : "");
      setEditAssignedTo(t.assigned_to);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    api.getUsers().then((u) => setUsers(u || [])).catch(() => {});
  }, [id]);

  const handleStateChange = async () => {
    if (!id || !newState) return;
    setChangingState(true);
    try {
      await api.updateTaskState(id, newState);
      fetchData();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Error al cambiar estado");
    } finally {
      setChangingState(false);
    }
  };

  const handleAddComment = async () => {
    if (!id || !comment.trim()) return;
    setAddingComment(true);
    try {
      await api.addTaskComment(id, comment);
      setComment("");
      fetchData();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Error al agregar comentario");
    } finally {
      setAddingComment(false);
    }
  };

  const handleSaveEdit = async () => {
    if (!id) return;
    setSaving(true);
    try {
      await api.updateTask(id, {
        title: editTitle,
        description: editDescription,
        expires_at: new Date(editExpiresAt).toISOString(),
        assigned_to: editAssignedTo,
      });
      setEditing(false);
      fetchData();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Error al actualizar");
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <p>Cargando...</p>;
  if (error) return <p className="text-red-500">{error}</p>;
  if (!task) return <p>Tarea no encontrada</p>;

  const expired = isExpired(task.expires_at);

  return (
    <div className="space-y-6 max-w-3xl">
      <Button variant="outline" onClick={() => navigate("/tasks")}>
        Volver
      </Button>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>{task.title}</CardTitle>
            <div className="flex gap-2">
              <Badge variant={expired ? "destructive" : "outline"}>
                {expired ? "Vencida" : "Vigente"}
              </Badge>
              <Badge>{task.task_state}</Badge>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {!editing ? (
            <>
              <div>
                <strong>Descripcion:</strong>
                <p>{task.description || "Sin descripcion"}</p>
              </div>
              <div>
                <strong>Vencimiento:</strong> {new Date(task.expires_at).toLocaleString()}
              </div>
              <div>
                <strong>Asignado a:</strong>{" "}
                {users.find((u) => u.id === task.assigned_to)?.email || task.assigned_to}
              </div>
              <div>
                <strong>Creado:</strong> {new Date(task.created_at).toLocaleString()}
              </div>
              {task.task_state === "Asignado" && users.length > 0 && (
                <Button variant="outline" onClick={() => setEditing(true)}>
                  Editar
                </Button>
              )}
            </>
          ) : (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label>Titulo</Label>
                <Input value={editTitle} onChange={(e) => setEditTitle(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Descripcion</Label>
                <Textarea value={editDescription} onChange={(e) => setEditDescription(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Vencimiento</Label>
                <Input type="datetime-local" value={editExpiresAt} onChange={(e) => setEditExpiresAt(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>Asignar a</Label>
                <Select value={editAssignedTo} onValueChange={setEditAssignedTo}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {users.map((u) => (
                      <SelectItem key={u.id} value={u.id}>
                        {u.first_name} {u.last_name} ({u.email})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="flex gap-2">
                <Button onClick={handleSaveEdit} disabled={saving}>
                  {saving ? "Guardando..." : "Guardar"}
                </Button>
                <Button variant="outline" onClick={() => setEditing(false)}>Cancelar</Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* State change */}
      <Card>
        <CardHeader>
          <CardTitle>Cambiar Estado</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {expired ? (
            <p className="text-muted-foreground">La tarea esta vencida, no se puede cambiar el estado.</p>
          ) : (
            <div className="flex gap-2 items-end">
              <div className="flex-1 space-y-2">
                <Label>Nuevo estado</Label>
                <Select value={newState} onValueChange={setNewState}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {TASK_STATES.map((s) => (
                      <SelectItem key={s} value={s}>{s}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <Button onClick={handleStateChange} disabled={changingState || newState === task.task_state}>
                {changingState ? "Actualizando..." : "Actualizar"}
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Comments */}
      <Card>
        <CardHeader>
          <CardTitle>Comentarios ({comments.length})</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {comments.length === 0 && (
            <p className="text-muted-foreground">Sin comentarios.</p>
          )}
          {comments.map((c) => (
            <div key={c.id} className="border rounded p-3">
              <div className="text-sm text-muted-foreground">
                {users.find((u) => u.id === c.user_id)?.email || c.user_id} -{" "}
                {new Date(c.created_at).toLocaleString()}
              </div>
              <p className="mt-1">{c.comment}</p>
            </div>
          ))}

          <div className="space-y-2 pt-4 border-t">
            <Label>Agregar comentario</Label>
            <Textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Escribe un comentario..."
            />
            <Button onClick={handleAddComment} disabled={addingComment || !comment.trim()}>
              {addingComment ? "Enviando..." : "Enviar"}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
