import React, { useContext } from "react";
import { Link, useNavigate } from "react-router-dom";
import { AuthContext } from "../auth/AuthContext";

export default function Header() {
  const { user, logout } = useContext(AuthContext);
  const nav = useNavigate();

  const onLogout = () => {
    logout();
    nav('/login');
  };

  return (
    <header className="bg-white shadow">
      <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
        <Link to="/projects" className="font-bold">Система контроля</Link>
        <nav>
          <Link to="/projects" className="mr-4">Проекты</Link>
          {user ? (
            <>
              {(user.role === 'manager' || user.role === 'admin') && <Link to="/projects/create" className="mr-4">Создать проект</Link>}
              {user.role === 'admin' && <Link to="/admin/users" className="mr-4 text-sm text-gray-700">Администрирование</Link>}
              <Link to="/me" className="mr-4" aria-label={`Профиль пользователя ${user.name || ''}`}>
                <span className="font-medium">{user.name || 'Профиль'}</span>
                <span className="ml-2 inline-block text-xs text-gray-500 px-2 py-0.5 bg-gray-100 rounded">{(user.role || '—').toUpperCase()}</span>
              </Link>
              <button onClick={onLogout} className="text-sm text-red-600">Выйти</button>
            </>
          ) : (
            <Link to="/login">Вход</Link>
          )}
        </nav>
      </div>
    </header>
  );
}
