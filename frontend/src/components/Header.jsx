import React from "react";
import { Link } from "react-router-dom";

export default function Header() {
  return (
    <header className="bg-white shadow">
      <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
        <Link to="/projects" className="font-bold">Система контроля</Link>
        <nav>
          <Link to="/projects" className="mr-4">Проекты</Link>
          <Link to="/login">Вход</Link>
        </nav>
      </div>
    </header>
  );
}
