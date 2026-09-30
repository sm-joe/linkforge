"use client";

import { useEffect, useState } from "react";
import { listLinks } from "@/lib/api";
import type { Link } from "@/types/link";
import LinkCard from "./LinkCard";

export default function Dashboard() {
  const [links, setLinks] = useState<Link[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    async function fetchLinks() {
      try {
        const data = await listLinks();

        if (!cancelled) {
          setLinks(data);
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
    setLoading(true);

    const data = await listLinks();

    setLinks(data);
    setLoading(false);
  }

  if (loading) {
    return (
      <p className="mt-10 text-slate-400">
        Loading links...
      </p>
    );
  }

  return (
    <section className="mt-12">
      <h2 className="mb-5 text-2xl font-semibold">
        Your Links
      </h2>

      <div className="space-y-4">
        {links.length === 0 && (
          <p className="text-slate-400">
            No links created yet.
          </p>
        )}

        {links.map((link) => (
          <LinkCard
            key={link.id}
            link={link}
            onChange={reload}
          />
        ))}
      </div>
    </section>
  );
}