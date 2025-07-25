'use client';
import { signIn } from 'next-auth/react';
import { useState } from 'react';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    const res = await signIn('credentials', {
      email,
      password,
      callbackUrl: '/dashboard',
      redirect: true,
    });
    // NextAuth will handle redirect, but you can show error if needed
    if (res?.error) setError('Invalid credentials');
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50">
      <div className="bg-white rounded-lg shadow-lg p-8 max-w-sm w-full flex flex-col items-center">
        <h1 className="text-2xl font-bold mb-4">Sign in to Scain</h1>
        <form onSubmit={handleSubmit} className="w-full flex flex-col gap-4">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="submit"
            className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-blue-600 text-white rounded font-semibold hover:bg-blue-700 transition"
          >
            Sign in
          </button>
          {error && <div className="text-red-600 text-sm text-center">{error}</div>}
        </form>
        <div className="mt-4 text-xs text-gray-500 text-center">
          <div>Demo: Use any email and password</div>
        </div>
      </div>
    </div>
  );
} 