"use client";

import type { Link } from "@/types/link";
import {
  disableLink,
  enableLink,
  deleteLink,
} from "@/lib/api";

interface Props {
  link: Link;
  onChange: () => void;
}

export default function LinkCard({
  link,
  onChange,
}: Props) {

  async function disable() {
    await disableLink(link.short_code);
    onChange();
  }

  async function enable() {
    await enableLink(link.short_code);
    onChange();
  }

  async function remove() {
    await deleteLink(link.short_code);
    onChange();
  }


  return (
    <div className="rounded-xl border border-slate-800 bg-slate-900 p-5">

      <div className="flex justify-between">

        <div>
          <p className="font-semibold text-white">
            /{link.short_code}
          </p>

          <p className="mt-1 max-w-xl truncate text-sm text-slate-400">
            {link.destination}
          </p>
        </div>


        <span
          className={
            link.status === "active"
              ? "text-green-400"
              : "text-red-400"
          }
        >
          ● {link.status}
        </span>

      </div>


      <div className="mt-5 flex gap-3">

        <a
          href={link.short_url}
          target="_blank"
          className="rounded-lg bg-white px-4 py-2 text-sm text-slate-950"
        >
          Open
        </a>


        {link.status === "active" ? (
          <button
            onClick={disable}
            className="rounded-lg border border-slate-700 px-4 py-2 text-sm"
          >
            Disable
          </button>
        ) : (
          <button
            onClick={enable}
            className="rounded-lg border border-slate-700 px-4 py-2 text-sm"
          >
            Enable
          </button>
        )}


        <button
          onClick={remove}
          className="rounded-lg border border-red-900 px-4 py-2 text-sm text-red-400"
        >
          Delete
        </button>

      </div>

    </div>
  );
}