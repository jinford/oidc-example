import React, { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { config } from "../config";

type PostLoginResponse = {
  redirect_to: string;
}

export const Login: React.FC = () => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [loginChallenge, setLoginChallenge] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [searchParams] = useSearchParams();

  useEffect(() => {
    const challenge = searchParams.get("login_challenge");
    if (challenge) {
      setLoginChallenge(challenge);
    } else {
      setError("Missing login_challenge parameter");
    }
  }, [searchParams]);

  // ログイン処理
  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError(null);

    try {
      const response = await fetch(`${config.apiUrl}/api/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: username,
          password: password,
          login_challenge: loginChallenge,
        })
      });

      if (!response.ok) {
        throw new Error("Failed to login");
      }

      const data: PostLoginResponse = await response.json();
      console.log(data);
      if (!data.redirect_to) {
        throw new Error("Redirect URL is undefined");
      }
      window.location.href = data.redirect_to; // redirect to hydra
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    }
  };

  return (
    <div>
      <h1>IdPのログインページ</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="username">ユーザー名：</label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">パスワード：</label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit" disabled={!loginChallenge}>
          ログイン
        </button>
      </form>
    </div>
  );
};
