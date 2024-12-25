import React, { useState } from "react";
import { config } from "../config";

type GetOidcUrlResponse = {
  redirect_url: string;
}

export const Login: React.FC = () => {
  const [error, setError] = useState<string | null>(null);

  const handleLogin = async () => {
    try {
      const response = await fetch(`${config.apiUrl}/api/oidc-url`);
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data: GetOidcUrlResponse = await response.json();
      window.location.href = data.redirect_url;
    } catch (error) {
      setError('failed to get oidc-url: ' + error);
    }
  };

  return (
    <>
      <h1>RPのログインページ</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <button onClick={handleLogin}>
        IdP Appでログインする
      </button>
    </>
  );
};
