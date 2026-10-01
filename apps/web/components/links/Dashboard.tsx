"use client";

import { useEffect, useState } from "react";
import { listLinks } from "@/lib/api";
import type { Link } from "@/types/link";
import LinkCard from "./LinkCard";

export default function Dashboard() {
  const [links, setLinks] = useState<Link[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function fetchLinks() {
      try {
        setError("");

        const data = await listLinks();

        if (!cancelled) {
          setLinks(data);
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load links",
          );
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    fetchLinks();

    return () => {
      cancelled = true;
    };
  }, []);

  async function reload() {
    try {
      setError("");
      setLoading(true);

      const data = await listLinks();

      setLinks(data);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Unable to load links",
      );
    } finally {
      setLoading(false);
    }
  }

  if (loading) {
    return (
      <section className="mt-12">
        <h2 className="mb-5 text-2xl font-semibold">
          Your Links
        </h2>

        <p className="text-slate-400">
          Loading links...
        </p>
      </section>
    );
  }

  return (
    <section className="mt-12">
      <h2 className="mb-5 text-2xl font-semibold">
        Your Links
      </h2>

      {error && (
        <div className="mb-5 rounded-lg border border-red-900 bg-red-950 p-4 text-sm text-red-300">
          {error}
        </div>
      )}

      <div className="space-y-4">
        {links.length === 0 ? (
          <p className="text-slate-400">
            No links created yet.
          </p>
        ) : (
          links.map((link) => (
            <LinkCard
              key={link.id}
              link={link}
              onChange={reload}
            />
          ))
        )}
      </div>
    </section>
  );
}