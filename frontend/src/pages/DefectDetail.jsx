import React, { useEffect, useState, useRef } from "react";
import { useParams } from "react-router-dom";
import api from "../api/axios";
import Header from "../components/Header";

export default function DefectDetail() {
  const { id, defectId } = useParams();
  const [defect, setDefect] = useState(null);
  const [attachments, setAttachments] = useState([]);
  const [previews, setPreviews] = useState({}); // map attachment id -> objectURL
  const [comments, setComments] = useState([]);
  const [commentsLoading, setCommentsLoading] = useState(true);
  const [newComment, setNewComment] = useState("");
  const [posting, setPosting] = useState(false);
  const isMounted = useRef(true);

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

  // ensure isMounted correctly reflects component lifecycle (mount/unmount)
  useEffect(() => {
    isMounted.current = true;
    return () => { isMounted.current = false };
  }, []);

  // load comments
  useEffect(() => {
    let cancelled = false;
    const loadComments = async () => {
      setCommentsLoading(true);
      try {
        const rc = await api.get(`/comments?defect_id=${defectId}`);
        // store serializable debug snapshot; if axios didn't populate data, try to read raw responseText
        let rawText = null;
        try {
          if ((rc?.data === null || rc?.data === undefined) && rc?.request && rc.request.responseText) {
            rawText = rc.request.responseText;
          }
        } catch (e) { rawText = null }

        // attempt to build parsed body: prefer rc.data, else try JSON.parse(rawText)
        let parsedBody = null;
        if (rc?.data !== null && rc?.data !== undefined) parsedBody = rc.data;
        else if (rawText) {
          try { parsedBody = JSON.parse(rawText); } catch (e) { parsedBody = null }
        }

  const snap = { status: rc?.status ?? null, statusText: rc?.statusText ?? null, headers: rc?.headers ?? null, body: parsedBody, bodyText: rawText };

        // parse comments: api returns { data: [...] , status: 'ok' }
        const body = parsedBody ?? null;
        let list = [];
        if (Array.isArray(body)) list = body;
        else if (body && Array.isArray(body.data)) list = body.data;
        else list = [];

        // If axios didn't return useful data, fallback to fetch to ensure browser sees the same response as DevTools
        if ((list.length === 0) && typeof window !== 'undefined') {
          try {
            const token = localStorage.getItem('token');
            const url = `/api/v1/comments?defect_id=${encodeURIComponent(defectId)}`;
            const fresp = await fetch(url, { headers: token ? { Authorization: `Bearer ${token}` } : {} });
            const ftext = await fresp.text();
            let fbody = null;
            try { fbody = ftext ? JSON.parse(ftext) : null } catch (e) { fbody = null }
            const fList = Array.isArray(fbody) ? fbody : (fbody && Array.isArray(fbody.data) ? fbody.data : []);
            if (fList.length > 0) {
              list = fList;
            }
          } catch (e) {
            console.warn('fetch fallback failed', e);
          }
        }

        if (!cancelled && isMounted.current) setComments(list);
      } catch (err) {
        console.error('load comments failed', err);
        if (!cancelled && isMounted.current) setComments([]);
      } finally {
        if (!cancelled && isMounted.current) setCommentsLoading(false);
      }
    };
    loadComments();
    return () => { cancelled = true };
  }, [id, defectId]);

  // create previews for image attachments using authenticated request
  useEffect(() => {
    let cancelled = false;
    const createPreviews = async () => {
      const map = {};
      const createdUrls = [];
      for (const a of attachments) {
        try {
          const isImage = (a.content_type && a.content_type.startsWith("image/")) || (a.filename && /\.(jpe?g|png|gif|webp)$/i.test(a.filename));
          if (!isImage) continue;
          // fetch via axios to include Authorization header
          const resp = await api.get(`/attachments/${a.id}`, { responseType: 'blob' });
          const url = URL.createObjectURL(resp.data);
          map[a.id] = url;
          createdUrls.push(url);
        } catch (err) {
          // ignore preview for this file
          console.error('preview failed', a.id, err?.message || err);
        }
        if (cancelled) break;
      }
      if (!cancelled && isMounted.current) setPreviews((p) => ({ ...p, ...map }));
      // cleanup helper: revoke created urls when this effect is torn down
      return () => {
        createdUrls.forEach((u) => URL.revokeObjectURL(u));
      };
    };
    const cleanupPromise = createPreviews();
    return () => {
      cancelled = true;
      // if the async createPreviews returned a cleanup function, call it
      if (cleanupPromise && typeof cleanupPromise.then === 'function') {
        // wait for it and then call the returned cleanup
        cleanupPromise.then((fn) => { if (typeof fn === 'function') fn(); }).catch(() => {});
      }
      // also clear previews state
      setPreviews({});
    };
  }, [attachments]);

  const handleDownload = async (e, a) => {
    e.preventDefault();
    try {
      const resp = await api.get(`/attachments/${a.id}`, { responseType: 'blob' });
      const url = URL.createObjectURL(resp.data);
      const link = document.createElement('a');
      link.href = url;
      link.download = a.filename || `file-${a.id}`;
      document.body.appendChild(link);
      link.click();
      link.remove();
      setTimeout(() => URL.revokeObjectURL(url), 5000);
    } catch (err) {
      console.error('download failed', err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="max-w-3xl mx-auto p-6">
        <h1 className="text-2xl font-bold mb-4">Дефект: {defect?.title || defectId}</h1>
        <section className="mb-4 bg-white p-4 rounded">
          <div className="text-sm text-gray-700 mb-2">{defect?.description}</div>
          <div className="text-xs text-gray-500">Приоритет: {defect?.priority || "-"}</div>
          <div className="text-xs text-gray-500">Исполнитель: {defect?.assignee?.name || defect?.assignee_id || "-"}</div>
          <div className="text-xs text-gray-500">Срок: {defect?.due_date ? (typeof defect.due_date === 'string' ? defect.due_date : new Date(defect.due_date).toLocaleDateString()) : "-"}</div>
        </section>

        <section className="bg-white p-4 rounded">
          <h2 className="font-semibold mb-2">Вложения</h2>
          {attachments.length === 0 && <div className="text-gray-600">Нет вложений</div>}
          <ul className="space-y-2">
            {attachments.map((a) => (
              <li key={a.id} className="flex items-center space-x-3">
                {previews[a.id] ? (
                  <a href={`/api/v1/attachments/${a.id}`} onClick={(e) => handleDownload(e, a)}>
                    <img src={previews[a.id]} alt={a.filename} className="w-24 h-16 object-cover rounded" />
                  </a>
                ) : null}
                <div>
                  <a className="text-blue-600 cursor-pointer" onClick={(e) => handleDownload(e, a)} href={`/api/v1/attachments/${a.id}`}>{a.filename || a.original_name || `file-${a.id}`}</a>
                  <div className="text-xs text-gray-500">{a.content_type || ''} {a.size ? `· ${Math.round(a.size/1024)} KB` : ''}</div>
                </div>
              </li>
            ))}
          </ul>
        </section>

        <section className="mt-4 bg-white p-4 rounded">
          <h2 className="font-semibold mb-2">Комментарии</h2>
          {commentsLoading ? (
            <div className="text-gray-600 mb-3">Загрузка комментариев...</div>
          ) : (
            comments.length === 0 && <div className="text-gray-600 mb-3">Нет комментариев</div>
          )}
          <ul className="space-y-3 mb-4">
            {comments.map((c) => (
              <li key={c.id} className="border p-3 rounded">
                <div className="text-xs text-gray-500 mb-1">{c.author && c.author.name ? c.author.name : (c.author_id ? `User #${c.author_id}` : 'Unknown')} · {c.created_at ? new Date(c.created_at).toLocaleString() : ''}</div>
                <div className="text-sm text-gray-800 whitespace-pre-wrap">{c.body}</div>
              </li>
            ))}
          </ul>
            {/* Debug: raw comments payload (remove in production)
            <div className="mb-3">
              <div className="text-xs text-gray-500 mb-1">Debug: raw response / parsed comments (remove in prod)</div>
              <pre className="text-xs text-gray-600 bg-gray-50 p-2 rounded overflow-x-auto">Raw: {JSON.stringify(commentsRaw, null, 2)}
  Parsed: {JSON.stringify(comments, null, 2)}</pre>
            </div> */}

          <div>
            <label className="block text-sm font-medium text-gray-700">Добавить комментарий</label>
            <textarea value={newComment} onChange={(e) => setNewComment(e.target.value)} rows={4} className="mt-1 block w-full rounded border-gray-300 shadow-sm" />
            <div className="mt-2 flex items-center space-x-2">
              <button disabled={posting || newComment.trim() === ''} onClick={async () => {
                if (newComment.trim() === '') return;
                setPosting(true);
                try {
                  const res = await api.post(`/projects/${id}/defects/${defectId}/comments`, { body: newComment.trim() });
                  const created = res.data.data || res.data;
                  // append to list
                  setComments((s) => [...s, created]);
                  setNewComment("");
                } catch (err) {
                  console.error('post comment failed', err);
                  // optionally show an error toast
                } finally {
                  setPosting(false);
                }
              }} className="inline-flex items-center px-3 py-1.5 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50">{posting ? 'Отправка...' : 'Отправить'}</button>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
