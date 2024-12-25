import React, { useEffect, useState } from 'react';
import { config } from "../config";

type GetUserInfoResponse = {
  subject: string;
};

export const Callback: React.FC = () => {
  const [error, setError] = useState<string | null>(null);
  const [userInfo, setUserInfo] = useState<GetUserInfoResponse | null>(null);

  useEffect(() => {
    const searchParams = new URLSearchParams(window.location.search);
    const code = searchParams.get('code');

    if (code) {
      const fetchUserInfo = async () => {
        try {
          const response = await fetch(`${config.apiUrl}/api/user-info?code=${code}`);
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          const data: GetUserInfoResponse = await response.json();
          setUserInfo(data);
        } catch (error) {
          setError('Failed to fetch user info'+ error);
        }
      };

      fetchUserInfo();
    } else {
      setError('Missing code parameter');
    }
  }, []);

  return (
    <div>
      <h1>RPのコールバックページ</h1>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {!userInfo && !error && <p>認証処理中...</p>}
      {userInfo && <p>認証成功: {userInfo.subject}</p>}
    </div>
  );

};
