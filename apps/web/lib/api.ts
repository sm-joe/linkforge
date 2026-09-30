import type { Link } from "@/types/link";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080";

export async function createLink(
  destination: string,
  alias?: string,
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
      }),
    },
  );

  if (!response.ok) {
    const error = await response.json();

    throw new Error(
      error.error ?? "Unable to create link",
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
    throw new Error(
      "Unable to fetch links",
    );
  }

  return response.json();
}


export async function disableLink(
  code: string,
) {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}/disable`,
    {
      method: "POST",
    },
  );

  if (!response.ok) {
    throw new Error(
      "Unable to disable link",
    );
  }
}


export async function enableLink(
  code: string,
) {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}/enable`,
    {
      method: "POST",
    },
  );

  if (!response.ok) {
    throw new Error(
      "Unable to enable link",
    );
  }
}


export async function deleteLink(
  code: string,
) {
  const response = await fetch(
    `${API_URL}/api/v1/links/${code}`,
    {
      method: "DELETE",
    },
  );

  if (!response.ok) {
    throw new Error(
      "Unable to delete link",
    );
  }
}