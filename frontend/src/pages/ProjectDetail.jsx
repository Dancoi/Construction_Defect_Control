import React, { useEffect, useState, useRef, useContext } from "react";
import { useParams, Link } from "react-router-dom";
import api from "../api/axios";
import Header from "../components/Header";
import { AuthContext } from "../auth/AuthContext";

export default function ProjectDetail() {
  const { id } = useParams();
  const [project, setProject] = useState(null);
  const [defects, setDefects] = useState([]);
  const [form, setForm] = useState({ title: "", description: "", assignee_id: "", due_date: "", priority: "medium" });
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [users, setUsers] = useState([]);
  const [assigneeQuery, setAssigneeQuery] = useState("");
  const [assigneeResults, setAssigneeResults] = useState([]);
  const [assigneeOpen, setAssigneeOpen] = useState(false);
  const [previews, setPreviews] = useState([]);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const previewUrlsRef = useRef([]);
  const { user } = useContext(AuthContext);
  const canCreateDefect = user && ["engineer", "manager", "admin"].includes(user.role);

  // debounced search for assignees
  useEffect(() => {
    if (!assigneeQuery) {
      setAssigneeResults([]);
      return;
    }
    const t = setTimeout(async () => {
      try {
        const r = await api.get(`/users?search=${encodeURIComponent(assigneeQuery)}`);
        setAssigneeResults(r.data.data || r.data || []);
      } catch (e) {
        console.error("assignee search failed", e);
        setAssigneeResults([]);
      }
    }, 300);
    return () => clearTimeout(t);
  }, [assigneeQuery]);

  useEffect(() => {
    const load = async () => {
      try {
        const resP = await api.get(`/projects/${id}`);
        setProject(resP.data.data || resP.data);
      } catch (e) {
        console.error(e);
      }
      try {
        const res = await api.get(`/projects/${id}/defects`);
        setDefects(res.data.data || res.data || []);
      } catch (e) {
        console.error(e);
      }
      try {
        const ru = await api.get(`/users`);
        setUsers(ru.data.data || ru.data || []);
      } catch (e) {
        // ignore if endpoint not present
      }
    };
    load();
  }, [id]);

  const handleFilesChange = (e) => {
    const list = Array.from(e.target.files || []);
    previewUrlsRef.current.forEach((u) => URL.revokeObjectURL(u));
    previewUrlsRef.current = [];
    const p = list.map((f) => {
      const url = URL.createObjectURL(f);
      previewUrlsRef.current.push(url);
      return { file: f, url };
    });
    setPreviews(p);
    setFiles(list);
  };

  const createDefect = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const payload = { ...form };
      if (!payload.assignee_id) delete payload.assignee_id;
      else payload.assignee_id = Number(payload.assignee_id);

      const res = await api.post(`/projects/${id}/defects`, payload);
      const created = res.data.data || res.data;
      setDefects((d) => [created, ...d]);

      if (files && files.length > 0) {
        try {
          setUploading(true);
          setUploadProgress(0);
          const fd = new FormData();
          for (let i = 0; i < files.length; i++) fd.append("files", files[i]);
          const token = localStorage.getItem("token");
          const headers = token ? { Authorization: `Bearer ${token}` } : undefined;
          await api.post(`/projects/${id}/attachments?defect_id=${created.id}`, fd, {
            headers,
            onUploadProgress: (e) => {
              if (e.lengthComputable) setUploadProgress(Math.round((e.loaded / e.total) * 100));
            },
          });
        } catch (upErr) {
          console.error("attachments upload failed", upErr);
        } finally {
          setUploading(false);
          setUploadProgress(0);
          previewUrlsRef.current.forEach((u) => URL.revokeObjectURL(u));
          previewUrlsRef.current = [];
          setPreviews([]);
        }
      }

      setForm({ title: "", description: "", assignee_id: "", due_date: "", priority: "medium" });
      setFiles([]);
    } catch (err) {
      console.error(err);
      setError(err.response?.data?.message || "Ошибка при создании дефекта");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="max-w-3xl mx-auto p-6">
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-2xl font-bold">Проект: {project?.name || id}</h1>
          <Link to="/projects" className="text-sm text-blue-600">← Все проекты</Link>
        </div>

        <section className="mb-6">
          <h2 className="font-semibold mb-2">Создать дефект</h2>
          {canCreateDefect ? (
            <form onSubmit={createDefect} className="space-y-3 bg-white p-4 rounded">
              <div>
                <label className="block text-sm">Заголовок</label>
                <input
                  value={form.title}
                  onChange={(e) => setForm({ ...form, title: e.target.value })}
                  className="w-full border rounded p-2"
                  required
                />
              </div>

              <div className="relative">
                <label className="block text-sm">Исполнитель (assignee)</label>
                <input
                  value={assigneeQuery}
                  onChange={(e) => {
                    const v = e.target.value;
                    setAssigneeQuery(v);
                    setAssigneeOpen(true);
                  }}
                  onFocus={() => setAssigneeOpen(true)}
                  placeholder={form.assignee_id ? users.find(u => String(u.id) === String(form.assignee_id))?.name || "" : "Поиск исполнителя..."}
                  className="w-full border rounded p-2"
                />
                {assigneeOpen && assigneeResults.length > 0 && (
                  <ul className="absolute z-20 left-0 right-0 bg-white border rounded mt-1 max-h-48 overflow-auto">
                    {assigneeResults.map((u) => (
                      <li key={u.id} className="p-2 hover:bg-gray-100 cursor-pointer" onMouseDown={() => {
                        setForm({ ...form, assignee_id: u.id });
                        setAssigneeQuery(u.name || u.email);
                        setAssigneeOpen(false);
                      }}>{u.name || u.email}</li>
                    ))}
                  </ul>
                )}
                <input type="hidden" value={form.assignee_id} />
              </div>

              <div>
                <label className="block text-sm">Срок (due date)</label>
                <input
                  type="date"
                  value={form.due_date}
                  onChange={(e) => setForm({ ...form, due_date: e.target.value })}
                  className="w-full border rounded p-2"
                />
              </div>

              <div>
                <label className="block text-sm">Приоритет</label>
                <select
                  value={form.priority}
                  onChange={(e) => setForm({ ...form, priority: e.target.value })}
                  className="w-full border rounded p-2"
                >
                  <option value="low">Низкий</option>
                  <option value="medium">Средний</option>
                  <option value="high">Высокий</option>
                </select>
              </div>

              <div>
                <label className="block text-sm">Вложения</label>
                <input type="file" multiple onChange={handleFilesChange} className="w-full" />
                {previews.length > 0 && (
                  <div className="mt-2 grid grid-cols-3 gap-2">
                    {previews.map((p, idx) => (
                      <div key={idx} className="border rounded overflow-hidden">
                        {p.file.type.startsWith("image/") ? (
                          <img src={p.url} alt={p.file.name} className="w-full h-24 object-cover" />
                        ) : (
                          <div className="p-2 text-sm">{p.file.name}</div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div>
                <label className="block text-sm">Описание</label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  className="w-full border rounded p-2"
                  rows={4}
                />
              </div>

              {error && <div className="text-red-600">{error}</div>}
              <div>
                <button className="bg-blue-600 text-white px-4 py-2 rounded" disabled={loading || uploading}>
                  {loading ? "Сохранение..." : uploading ? `Загрузка... ${uploadProgress}%` : "Создать"}
                </button>
              </div>
            </form>
          ) : (
            <div className="bg-white p-4 rounded text-sm text-gray-600">У вас нет прав для создания дефекта в этом проекте.</div>
          )}
        </section>

        <section>
          <h2 className="font-semibold mb-2">Дефекты</h2>
          <div className="space-y-3">
            {defects.length === 0 && <div className="text-gray-600">Нет дефектов</div>}
            {defects.map((d) => (
              <Link key={d.id} to={`/projects/${id}/defects/${d.id}`} className="block bg-white p-3 rounded border hover:shadow">
                <h3 className="font-bold">{d.title}</h3>
                <div className="text-sm text-gray-700">{d.description}</div>
              </Link>
            ))}
          </div>
        </section>
      </main>
    </div>
  );
}
