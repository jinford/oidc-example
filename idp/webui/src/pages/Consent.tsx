import React, { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { config } from "../config";

type GetConsentResponse = {
  requested_scope: string[];
  grant_access_token_audience: string[];
};

type PostConsentResponse = {
  redirect_to: string;
};

export const Consent: React.FC = () => {
  const [requestedScope, setRequestedScope] = useState<string[]>([]);
  const [grantScope, setGrantScope] = useState<string[]>([]);
  const [grantAccessTokenAudience, setGrantAccessTokenAudience] = useState<string[]>([]);
  const [consentChallenge, setConsentChallenge] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [searchParams] = useSearchParams();

  useEffect(() => {
    setError(null);

    const challenge = searchParams.get("consent_challenge");
    if (challenge) {
      setConsentChallenge(challenge);
    } else {
      setError("Missing consent_challenge parameter");
      return
    }

    const fetchConsentData = async () => {
      try {
        const response = await fetch(`${config.apiUrl}/api/consent?consent_challenge=${consentChallenge}`);
        if (!response.ok) {
          throw new Error(`Error fetching consent data: ${response.statusText}`);
        }
        const data: GetConsentResponse = await response.json();
        setRequestedScope(data.requested_scope);
        setGrantAccessTokenAudience(data.grant_access_token_audience);
      } catch (error) {
        console.error(error);
        setError("Failed to fetch consent data");
      }
    };

    fetchConsentData();
  }, [searchParams, consentChallenge]);

  const handleScopeChange = (scope: string) => {
    setGrantScope((prev) =>
      prev.includes(scope)
        ? prev.filter((s) => s !== scope)
        : [...prev, scope]
    );
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError(null);

    if (!consentChallenge) {
      setError("No consent_challenge available.");
      return;
    }

    try {
      const response = await fetch(`${config.apiUrl}/api/consent`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
            consent_challenge: consentChallenge,
            grant_scope: grantScope,
            grant_access_token_audience: grantAccessTokenAudience,
         }),
      });

      if (!response.ok) {
        throw new Error("Failed to submit consent");
      }

      const data: PostConsentResponse = await response.json();
      window.location.href = data.redirect_to; // redirect to hydra
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    }
  };

  return (
    <div>
      <h1>同意ページ</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="grant_scope">スコープ:</label>
          {requestedScope.map((scope) => (
            <div key={scope}>
              <input
                type="checkbox"
                className="grant_scope"
                id={scope}
                value={scope}
                name="grant_scope"
                checked={grantScope.includes(scope)}
                onChange={() => handleScopeChange(scope)}
              />
              <label htmlFor={scope}>{scope}</label>
              <br />
            </div>
          ))}
        </div>
        <button type="submit" disabled={!consentChallenge}>
          許可
        </button>
      </form>
    </div>
  );
};