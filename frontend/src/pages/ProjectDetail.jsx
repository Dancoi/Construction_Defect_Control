import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../api/axios";
import Header from "../components/Header";

export default function ProjectDetail() {
  const { id } = useParams();
  const [project, setProject] = useState(null);
  const [defects, setDefects] = useState([]);
  const [form, setForm] = useState({ title: "", description: "", assignee: "", due_date: "", priority: "medium" });
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [users, setUsers] = useState([]);

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
      // try to load users for assignee autocomplete (optional)
      try {
        const ru = await api.get(`/users`);
        setUsers(ru.data.data || ru.data || []);
      } catch (e) {
        // ignore if endpoint not present
      }
    };
    load();
  }, [id]);

  const createDefect = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const res = await api.post(`/projects/${id}/defects`, form);
      const created = res.data.data || res.data;
      setDefects((d) => [created, ...d]);
      // if files selected, upload them to attachments endpoint referencing defect id
      if (files && files.length > 0) {
        try {
          const fd = new FormData();
          for (let i = 0; i < files.length; i++) fd.append("files", files[i]);
          // some backends expect defect_id as query param or form field; use query param here
          await api.post(`/projects/${id}/attachments?defect_id=${created.id}`, fd, {
            headers: { "Content-Type": "multipart/form-data" },
          });
        } catch (upErr) {
          console.error("attachments upload failed", upErr);
        }
      }
      setForm({ title: "", description: "", assignee: "", due_date: "", priority: "medium" });
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
        <h1 className="text-2xl font-bold mb-4">Проект: {project?.name || id}</h1>

        <section className="mb-6">
          <h2 className="font-semibold mb-2">Создать дефект</h2>
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
            <div>
              <label className="block text-sm">Исполнитель (assignee)</label>
              <input
                list="assignees"
                value={form.assignee}
                onChange={(e) => setForm({ ...form, assignee: e.target.value })}
                className="w-full border rounded p-2"
                placeholder="user id or email"
              />
              <datalist id="assignees">
                {users.map((u) => (
                  <option key={u.id} value={u.name || u.email || u.id} />
                ))}
              </datalist>
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
              <input
                type="file"
                multiple
                onChange={(e) => setFiles(Array.from(e.target.files || []))}
                className="w-full"
              />
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
              <button className="bg-blue-600 text-white px-4 py-2 rounded" disabled={loading}>
                {loading ? "Сохранение..." : "Создать"}
              </button>
            </div>
          </form>
        </section>

        <section>
          <h2 className="font-semibold mb-2">Дефекты</h2>
          <div className="space-y-3">
            {defects.length === 0 && <div className="text-gray-600">Нет дефектов</div>}
            {defects.map((d) => (
              <a key={d.id} href={`/projects/${id}/defects/${d.id}`} className="block bg-white p-3 rounded border hover:shadow">
                <h3 className="font-bold">{d.title}</h3>
                <div className="text-sm text-gray-700">{d.description}</div>
              </a>
            ))}
          </div>
        </section>
      </main>
    </div>
  );
}
