import React, { useEffect, useState } from "react";
import api from "../api/axios";
import Header from "../components/Header";
import { Link } from "react-router-dom";

export default function Projects() {
  const [list, setList] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get("/projects").then((res) => {
      setList(res.data.data || res.data || []);
    }).catch(() => setList([])).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="p-4">Загрузка...</div>;

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="max-w-5xl mx-auto p-4">
        <h1 className="text-2xl mb-4">Проекты</h1>
        <div className="grid gap-4">
          {list.map((p) => (
            <Link key={p.id} to={`/projects/${p.id}`} className="block border p-3 rounded bg-white hover:shadow">
              <h2 className="font-bold">{p.name}</h2>
              <div className="text-sm text-gray-600">{p.address}</div>
            </Link>
          ))}
        </div>
      </main>
    </div>
  );
}
