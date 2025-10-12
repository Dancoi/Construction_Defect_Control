import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../api/axios";
import Header from "../components/Header";

export default function DefectDetail() {
  const { id, defectId } = useParams();
  const [defect, setDefect] = useState(null);
  const [attachments, setAttachments] = useState([]);

  useEffect(() => {
    const load = async () => {
      try {
        const res = await api.get(`/projects/${id}/defects/${defectId}`);
        setDefect(res.data.data || res.data);
      } catch (e) {
        console.error(e);
      }
      try {
        const ra = await api.get(`/projects/${id}/defects/${defectId}/attachments`);
        setAttachments(ra.data.data || ra.data || []);
      } catch (e) {
        // attachments endpoint may be different; try generic
        try {
          const ra2 = await api.get(`/attachments?defect_id=${defectId}`);
          setAttachments(ra2.data.data || ra2.data || []);
        } catch (e2) {
          // ignore
        }
      }
    };
    load();
  }, [id, defectId]);

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="max-w-3xl mx-auto p-6">
        <h1 className="text-2xl font-bold mb-4">Дефект: {defect?.title || defectId}</h1>
        <section className="mb-4 bg-white p-4 rounded">
          <div className="text-sm text-gray-700 mb-2">{defect?.description}</div>
          <div className="text-xs text-gray-500">Приоритет: {defect?.priority || "-"}</div>
          <div className="text-xs text-gray-500">Исполнитель: {defect?.assignee || "-"}</div>
          <div className="text-xs text-gray-500">Срок: {defect?.due_date || "-"}</div>
        </section>

        <section className="bg-white p-4 rounded">
          <h2 className="font-semibold mb-2">Вложения</h2>
          {attachments.length === 0 && <div className="text-gray-600">Нет вложений</div>}
          <ul className="space-y-2">
            {attachments.map((a) => (
              <li key={a.id}>
                <a className="text-blue-600" href={`/api/v1/attachments/${a.id}`} target="_blank" rel="noreferrer">{a.filename || a.original_name || `file-${a.id}`}</a>
              </li>
            ))}
          </ul>
        </section>
      </main>
    </div>
  );
}
