import type { Link } from "@/types/link";

interface Props {
  link: Link;
}

export default function LinkHealth({
  link,
}: Props) {
  const isHTTPS =
    link.destination.startsWith("https://");

  return (
    <div className="mt-6 rounded-xl border border-slate-800 bg-slate-900 p-6">
      <h2 className="text-lg font-semibold text-white">
        Link Health
      </h2>

      <div className="mt-5 space-y-3 text-sm">

        <div className="flex justify-between">
          <span className="text-slate-400">
            Status
          </span>

          <span className="text-green-400">
            ● {link.status}
          </span>
        </div>


        <div className="flex justify-between gap-4">
          <span className="text-slate-400">
            Destination
          </span>

          <span className="max-w-xs truncate text-slate-200">
            {link.destination}
          </span>
        </div>


        <div className="flex justify-between">
          <span className="text-slate-400">
            HTTPS
          </span>

          <span>
            {isHTTPS ? "✓ Enabled" : "✗ Disabled"}
          </span>
        </div>


        <div className="flex justify-between">
          <span className="text-slate-400">
            Short code
          </span>

          <span className="text-slate-200">
            {link.short_code}
          </span>
        </div>


        <div className="flex justify-between">
          <span className="text-slate-400">
            Created
          </span>

          <span className="text-slate-200">
            {new Date(
              link.created_at,
            ).toLocaleString()}
          </span>
        </div>

      </div>
    </div>
  );
}