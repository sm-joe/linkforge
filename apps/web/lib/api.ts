import type { Link } from "@/types/link";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080";

export interface LinkAnalytics {
  short_code: string;
  total_clicks: number;
  first_click?: string | null;
  last_click?: string | null;
  referrers: {
    referrer: string;
    clicks: number;
  }[];
  browsers: {
    name: string;
    clicks: number;
  }[];
  devices: {
    name: string;
    clicks: number;
  }[];
}

export type LinkHealthStatus =
  | "healthy"
  | "degraded"
  | "unhealthy";

export interface LinkHealth {
  short_code: string;
  destination: string;
  status: LinkHealthStatus;
  http_status?: number;
  response_time_ms: number;
  final_url?: string;
  redirects: number;
  https: boolean;
  checked_at: string;
  error?: string;
}

interface LinkHealthHistoryResponse {
  short_code: string;
  checks: LinkHealth[];
}

export async function createLink(
  destination: string,
  alias?: string,
  expiresAt?: string,
): Promise<Link> {
  const response = await fetch(
    `${API_URL}/api/v1/links`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        destination,
        alias,
        expires_at: expiresAt,
      }),
    },
  );

  if (!response.ok) {
    const error = await response.json().catch(() => null);

    throw new Error(
      error?.error ?? "Unable to create link",
    );
  }

  return response.json();
}

export async function listLinks(): Promise<Link[]> {
  const response = await fetch(
    `${API_URL}/api/v1/links`,
    {
      cache: "no-store",
    },
  );

  if (!response.ok) {
    const error = await response.json().catch(() => null);

    throw new Error(
      error?.error ?? "Unable to fetch links",
    );
  }

  return response.json();
}

export async function getLinkAnalytics(
  code: string,
): Promise<LinkAnalytics> {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}/analytics`,
    {
      cache: "no-store",
    },
  );

  if (!response.ok) {
    const error = await response.json().catch(() => null);

    throw new Error(
      error?.error ?? "Unable to fetch analytics",
    );
  }

  return response.json();
}

export async function getLinkHealth(
  code: string,
): Promise<LinkHealth> {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}/health`,
    {
      cache: "no-store",
    },
  );

  if (!response.ok) {
    const error = await response.json().catch(() => null);

    throw new Error(
      error?.error ?? "Unable to fetch link health",
    );
  }

  return response.json();
}

export async function getLinkHealthHistory(
  code: string,
  limit = 50,
): Promise<LinkHealth[]> {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}/health/history?limit=${limit}`,
    {
      cache: "no-store",
    },
  );

  if (!response.ok) {
    const error = await response.json().catch(() => null);

    throw new Error(
      error?.error ?? "Unable to fetch health history",
    );
  }

  const data =
    (await response.json()) as LinkHealthHistoryResponse;

  return data.checks ?? [];
}

export async function disableLink(
  code: string,
): Promise<void> {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}/disable`,
    {
      method: "POST",
    },
  );

  if (!response.ok) {
    const error = await response.json().catch(() => null);

    throw new Error(
      error?.error ?? "Unable to disable link",
    );
  }
}

export async function enableLink(
  code: string,
): Promise<void> {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}/enable`,
    {
      method: "POST",
    },
  );

  if (!response.ok) {
    const error = await response.json().catch(() => null);

    throw new Error(
      error?.error ?? "Unable to enable link",
    );
  }
}

export async function deleteLink(
  code: string,
): Promise<void> {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}`,
    {
      method: "DELETE",
    },
  );

  if (!response.ok) {
    const error = await response.json().catch(() => null);

    throw new Error(
      error?.error ?? "Unable to delete link",
    );
  }
}