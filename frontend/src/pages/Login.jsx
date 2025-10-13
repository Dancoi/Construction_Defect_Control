import React, { useState, useContext } from "react";
import { AuthContext } from "../auth/AuthContext";
import { useNavigate, Link } from "react-router-dom";

export default function Login() {
  const { login } = useContext(AuthContext);
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const ok = await login(email, password);
      if (ok) navigate('/projects');
    } catch (err) {
      setError(err.response?.data?.error || err.message || "Login failed");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="w-full max-w-md p-8 bg-white rounded shadow">
        <h1 className="text-2xl mb-4">Вход</h1>
        {error && <div className="text-red-600">{error}</div>}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm">Email</label>
            <input value={email} onChange={(e) => setEmail(e.target.value)} className="w-full border p-2 rounded" />
          </div>
          <div>
            <label className="block text-sm">Пароль</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} className="w-full border p-2 rounded" />
          </div>
          <div>
            <button type="submit" className="w-full bg-blue-600 text-white p-2 rounded">Войти</button>
          </div>
          <div className="text-center text-sm text-gray-600">
            Нет аккаунта? <Link to="/register" className="text-blue-600">Зарегистрироваться</Link>
          </div>
        </form>
      </div>
    </div>
  );
}
