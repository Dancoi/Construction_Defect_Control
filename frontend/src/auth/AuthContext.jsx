import React, { createContext, useState, useEffect } from "react";
import api from "../api/axios";

export const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      api.get("/auth/me").then((res) => {
        setUser(res.data);
      }).catch(() => setUser(null)).finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (email, password) => {
    const res = await api.post("/auth/login", { email, password });
    // backend returns { status: 'ok', data: { token: ..., user: {...} } }
    const token = res?.data?.data?.token || res?.data?.token || res?.data?.access_token || null;
    if (token) {
      localStorage.setItem("token", token);
      // if backend returned user object, use it to avoid extra roundtrip
      const userObj = res?.data?.data?.user || null;
      if (userObj) {
        setUser(userObj);
      } else {
        // fetch user
        const me = await api.get("/auth/me");
        setUser(me.data?.data || me.data);
      }
      return true;
    }
    return false;
  };

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
