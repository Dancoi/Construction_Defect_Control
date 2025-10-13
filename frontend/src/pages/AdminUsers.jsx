import React, { useEffect, useState, useContext } from 'react'
import { AuthContext } from '../auth/AuthContext'
import Header from '../components/Header'
import api from '../api/axios'

export default function AdminUsers(){
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [confirm, setConfirm] = useState({ show: false, userId: null, role: null, name: '' })
  const { user: me } = useContext(AuthContext)

  useEffect(()=>{
    const load = async ()=>{
      setLoading(true)
      try{
        const r = await api.get('/users')
        const body = r.data?.data || r.data
        const list = (body || []).slice().sort((a,b) => Number(a.id) - Number(b.id))
        setUsers(list)
      }catch(e){
        console.error(e)
        setError('Не удалось загрузить пользователей')
      }finally{ setLoading(false) }
    }
    load()
  },[])

  const doChangeRole = async (userId, role) => {
    // prevent changing own role
    if (me && Number(me.id) === Number(userId)) {
      alert('Нельзя изменить свою роль');
      setConfirm({ show: false, userId: null, role: null, name: '' })
      return
    }
    try{
      await api.patch(`/users/${userId}`, { role })
      setUsers((s)=> {
        const updated = s.map(u => u.id === userId ? { ...u, role } : u)
        return updated.slice().sort((a,b) => Number(a.id) - Number(b.id))
      })
      setConfirm({ show: false, userId: null, role: null, name: '' })
    }catch(e){ console.error('change role failed', e); alert('Не удалось изменить роль') }
  }

  const changeRole = (userId, role, name) => {
    setConfirm({ show: true, userId, role, name })
  }

  const cancelConfirm = () => setConfirm({ show: false, userId: null, role: null, name: '' })

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main className="max-w-4xl mx-auto p-6">
        <h1 className="text-2xl font-semibold mb-4">Управление пользователями</h1>
        {loading ? <div>Загрузка...</div> : null}
        {error && <div className="text-red-600">{error}</div>}
        <table className="w-full bg-white rounded shadow overflow-hidden">
          <thead className="bg-gray-100">
            <tr>
              <th className="p-3 text-left">ID</th>
              <th className="p-3 text-left">Имя</th>
              <th className="p-3 text-left">Email</th>
              <th className="p-3 text-left">Роль</th>
              <th className="p-3 text-left">Действия</th>
            </tr>
          </thead>
          <tbody>
            {users.map(u => (
              <tr key={u.id} className="border-t">
                <td className="p-3">{u.id}</td>
                <td className="p-3">{u.name}</td>
                <td className="p-3">{u.email}</td>
                <td className="p-3">{u.role || '—'}</td>
                <td className="p-3 space-x-2">
                  <button onClick={()=>changeRole(u.id, 'engineer', u.name)} disabled={me && Number(me.id) === Number(u.id)} className="px-2 py-1 bg-gray-200 rounded disabled:opacity-50">Engineer</button>
                  <button onClick={()=>changeRole(u.id, 'manager', u.name)} disabled={me && Number(me.id) === Number(u.id)} className="px-2 py-1 bg-yellow-200 rounded disabled:opacity-50">Manager</button>
                  <button onClick={()=>changeRole(u.id, 'admin', u.name)} disabled={me && Number(me.id) === Number(u.id)} className="px-2 py-1 bg-red-200 rounded disabled:opacity-50">Admin</button>
                  {me && Number(me.id) === Number(u.id) && <span className="ml-2 text-xs text-gray-500">(это вы)</span>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {confirm.show && (
          <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
            <div className="bg-white rounded shadow-lg max-w-md w-full p-6">
              <h3 className="text-lg font-semibold mb-2">Подтвердите изменение роли</h3>
              <p className="text-sm text-gray-700 mb-4">Вы собираетесь назначить пользователю <strong>{confirm.name}</strong> роль <strong>{confirm.role}</strong>. Продолжить?</p>
              <div className="flex justify-end space-x-2">
                <button onClick={cancelConfirm} className="px-3 py-1 rounded border">Отмена</button>
                <button onClick={()=>doChangeRole(confirm.userId, confirm.role)} className="px-3 py-1 rounded bg-blue-600 text-white">Подтвердить</button>
              </div>
            </div>
          </div>
        )}
      </main>
    </div>
  )
}
