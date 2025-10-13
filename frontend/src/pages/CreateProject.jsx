import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../api/axios'
import Header from '../components/Header'

export default function CreateProject(){
  const [name, setName] = useState('')
  const [address, setAddress] = useState('')
  const [loading, setLoading] = useState(false)
  const [errors, setErrors] = useState({})
  const nav = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    const errs = {}
    if (!name.trim()) errs.name = 'Название обязательно'
    if (address && address.length < 5) errs.address = 'Адрес слишком короткий'
    setErrors(errs)
    if (Object.keys(errs).length > 0) return
    setLoading(true)
    try {
      const res = await api.post('/projects', { name: name.trim(), address: address.trim() })
      const created = res.data.data || res.data
      if (created && created.id) {
        nav(`/projects/${created.id}`)
      } else {
        // fallback: go to projects list
        nav('/projects')
      }
    } catch (err) {
      console.error('create project failed', err)
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main className="max-w-2xl mx-auto p-6">
        <h1 className="text-2xl font-semibold mb-4">Создать проект</h1>
        <form onSubmit={handleSubmit} className="bg-white p-6 rounded shadow">
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700">Название</label>
            <input placeholder="Например: ЖК Речной" value={name} onChange={(e)=>setName(e.target.value)} className={`mt-1 block w-full rounded border px-3 py-2 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-200 ${errors.name ? 'border-red-400' : 'border-gray-300'}`} />
            {errors.name && <div className="text-xs text-red-600 mt-1">{errors.name}</div>}
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700">Адрес</label>
            <input placeholder="Город, улица, дом" value={address} onChange={(e)=>setAddress(e.target.value)} className={`mt-1 block w-full rounded border px-3 py-2 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-200 ${errors.address ? 'border-red-400' : 'border-gray-300'}`} />
            {errors.address && <div className="text-xs text-red-600 mt-1">{errors.address}</div>}
          </div>
          <div className="flex items-center space-x-2">
            <button type="submit" disabled={loading} className="px-3 py-1.5 bg-blue-600 text-white rounded disabled:opacity-50">{loading ? 'Создание...' : 'Создать'}</button>
          </div>
        </form>
      </main>
    </div>
  )
}
