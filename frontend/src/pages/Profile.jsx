import React, { useState, useEffect, useContext } from 'react'
import api from '../api/axios'
import Header from '../components/Header'
import { AuthContext } from '../auth/AuthContext'

export default function Profile(){
  const { user, setUser } = useContext(AuthContext)
  const [form, setForm] = useState({ name: '', email: '' })
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [errors, setErrors] = useState({})

  useEffect(()=>{
    const load = async ()=>{
      setLoading(true)
      try{
        const r = await api.get('/users/me')
        const body = r.data.data || r.data
        setForm({ name: body.name || '', email: body.email || '' })
      }catch(e){ console.error(e) }
      setLoading(false)
    }
    load()
  },[])

  const validate = () => {
    const errs = {}
    if (!form.name || form.name.trim().length < 2) errs.name = 'Введите имя (мин. 2 символа)'
    if (!form.email || !/^\S+@\S+\.\S+$/.test(form.email)) errs.email = 'Введите корректный email'
    setErrors(errs)
    return Object.keys(errs).length === 0
  }

  const handleSave = async () => {
    if (!validate()) return
    setSaving(true)
    try{
      const r = await api.patch('/users/me', { name: form.name.trim(), email: form.email.trim() })
      const body = r.data.data || r.data
      // optimistic update in auth context if present
      try{ setUser && setUser(body) }catch(e){}
      // small inline success message rather than alert
      setErrors({ _success: 'Сохранено' })
    }catch(e){
      console.error(e)
      setErrors({ _error: 'Сохранение не удалось' })
    }finally{ setSaving(false) }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main className="max-w-2xl mx-auto p-6">
        <h1 className="text-2xl font-semibold mb-4">Личный кабинет</h1>
        {loading ? <div>Загрузка...</div> : (
          <div className="bg-white p-6 rounded shadow">
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">Имя</label>
              <input placeholder="Ваше полное имя" value={form.name} onChange={(e)=>setForm(s=>({...s, name: e.target.value}))} className={`mt-1 block w-full rounded border px-3 py-2 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-200 ${errors.name ? 'border-red-400' : 'border-gray-300'}`} />
              {errors.name && <div className="text-xs text-red-600 mt-1">{errors.name}</div>}
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">Email</label>
              <input placeholder="email@company.com" value={form.email} onChange={(e)=>setForm(s=>({...s, email: e.target.value}))} className={`mt-1 block w-full rounded border px-3 py-2 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-200 ${errors.email ? 'border-red-400' : 'border-gray-300'}`} />
              {errors.email && <div className="text-xs text-red-600 mt-1">{errors.email}</div>}
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">Роль</label>
              <div className="mt-1 text-sm text-gray-700">{user?.role || '—'}</div>
            </div>
            <div className="flex items-center space-x-2">
              <button onClick={handleSave} disabled={saving} className="px-3 py-1.5 bg-blue-600 text-white rounded disabled:opacity-50">{saving ? 'Сохранение...' : 'Сохранить'}</button>
              {errors._success && <div className="text-sm text-green-600">{errors._success}</div>}
              {errors._error && <div className="text-sm text-red-600">{errors._error}</div>}
            </div>
          </div>
        )}
      </main>
    </div>
  )
}
