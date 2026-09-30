"use client";

import { useState } from "react";
import { createLink } from "@/lib/api";
import type { Link } from "@/types/link";
import LinkHealth from "./LinkHealth";

export default function CreateLinkForm() {
  const [destination, setDestination] = useState("");
  const [alias, setAlias] = useState("");
  const [result, setResult] = useState<Link | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault();

    setError("");
    setResult(null);
    setLoading(true);

    try {
      const link = await createLink(
        destination,
        alias || undefined,
      );

      setResult(link);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Something went wrong",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto mt-10 w-full max-w-2xl">
      <form
        onSubmit={handleSubmit}
        className="space-y-4"
      >
        <input
          type="url"
          required
          value={destination}
          onChange={(event) =>
            setDestination(event.target.value)
          }
          placeholder="https://example.com/your-long-url"
          className="h-12 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 text-sm text-white outline-none placeholder:text-slate-500 focus:border-slate-400"
        />

        <input
          type="text"
          value={alias}
          onChange={(event) =>
            setAlias(event.target.value)
          }
          placeholder="Custom alias (optional)"
          className="h-12 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 text-sm text-white outline-none placeholder:text-slate-500 focus:border-slate-400"
        />

        <button
          type="submit"
          disabled={loading}
          className="h-12 w-full rounded-lg bg-white font-medium text-slate-950 transition hover:bg-slate-200 disabled:opacity-50"
        >
          {loading ? "Creating..." : "Shorten"}
        </button>
      </form>

      {error && (
        <div className="mt-5 rounded-lg border border-red-900 bg-red-950 p-4 text-sm text-red-300">
          {error}
        </div>
      )}

      {result && (
  <>
    <div className="mt-5 rounded-lg border border-slate-800 bg-slate-900 p-5">
      <p className="text-sm text-slate-400">
        Your short link
      </p>

      <a
        href={result.short_url}
        target="_blank"
        className="mt-2 block text-lg text-white underline"
      >
        {result.short_url}
      </a>
    </div>

    <LinkHealth link={result} />
  </>
)}
    </div>
  );
}