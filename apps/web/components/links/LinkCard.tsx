"use client";

import { useEffect, useState } from "react";

import type { Link } from "@/types/link";
import type {
  LinkAnalytics,
  LinkHealth,
} from "@/lib/api";

import {
  disableLink,
  enableLink,
  deleteLink,
  getLinkAnalytics,
  getLinkHealth,
  getLinkHealthHistory,
} from "@/lib/api";

interface Props {
  link: Link;
  onChange: () => void;
}

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080";

export default function LinkCard({
  link,
  onChange,
}: Props) {
  const [analytics, setAnalytics] =
    useState<LinkAnalytics | null>(null);

  const [analyticsError, setAnalyticsError] =
    useState("");

  const [showAnalytics, setShowAnalytics] =
    useState(false);

  const [showQRCode, setShowQRCode] =
    useState(false);

  const [health, setHealth] =
    useState<LinkHealth | null>(null);

  const [healthError, setHealthError] =
    useState("");

  const [healthHistory, setHealthHistory] =
    useState<LinkHealth[]>([]);

  const [healthHistoryLoading, setHealthHistoryLoading] =
    useState(false);

  const [healthHistoryError, setHealthHistoryError] =
    useState("");

  const [showHealth, setShowHealth] =
    useState(false);

  const [healthLoading, setHealthLoading] =
    useState(false);

  useEffect(() => {
    let cancelled = false;

    async function loadAnalytics() {
      try {
        setAnalyticsError("");

        const data =
          await getLinkAnalytics(
            link.short_code,
          );

        if (!cancelled) {
          setAnalytics(data);
        }
      } catch {
        if (!cancelled) {
          setAnalyticsError(
            "Analytics unavailable",
          );
        }
      }
    }

    loadAnalytics();

    return () => {
      cancelled = true;
    };
  }, [link.short_code]);

  const loadHealth = async () => {
    setHealthLoading(true);
    setHealthError("");

    setHealthHistoryLoading(true);
    setHealthHistoryError("");

    try {
      const [currentHealth, history] =
        await Promise.all([
          getLinkHealth(link.short_code),
          getLinkHealthHistory(
            link.short_code,
          ),
        ]);

      setHealth(currentHealth);
      setHealthHistory(history);
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : "Health check unavailable";

      setHealthError(message);
      setHealthHistoryError(message);
    } finally {
      setHealthLoading(false);
      setHealthHistoryLoading(false);
    }
  };

  async function toggleHealth() {
    const nextVisible = !showHealth;

    setShowHealth(nextVisible);

    if (
      nextVisible &&
      health === null
    ) {
      await loadHealth();
    }
  }

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

  function formatDate(
    value?: string | null,
  ) {
    if (!value) {
      return "—";
    }

    return new Date(
      value,
    ).toLocaleString();
  }

  function healthStatusClass(
    status?: LinkHealth["status"],
  ) {
    switch (status) {
      case "healthy":
        return "text-green-400";

      case "degraded":
        return "text-yellow-400";

      case "unhealthy":
        return "text-red-400";

      default:
        return "text-slate-500";
    }
  }

  const qrURL =
    `${API_URL}/api/v1/links/` +
    `${link.short_code}/qr`;

  return (
    <div className="rounded-xl border border-slate-800 bg-slate-900 p-5">
      <div className="flex justify-between gap-4">
        <div className="min-w-0">
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
              ? "shrink-0 text-green-400"
              : "shrink-0 text-red-400"
          }
        >
          ● {link.status}
        </span>
      </div>

      <div className="mt-5 flex items-center gap-6 border-y border-slate-800 py-4">
        <div>
          <p className="text-xs uppercase tracking-wide text-slate-500">
            Clicks
          </p>

          <p className="mt-1 text-xl font-semibold text-white">
            {analytics?.total_clicks ?? 0}
          </p>
        </div>

        {analyticsError && (
          <p className="text-xs text-slate-500">
            {analyticsError}
          </p>
        )}
      </div>

      {showQRCode && (
        <div className="mt-5 border-b border-slate-800 pb-5">
          <div className="flex flex-col items-start gap-3">
            <p className="text-sm font-medium text-white">
              QR Code
            </p>

            <div className="rounded-lg bg-white p-3">
              <img
                src={qrURL}
                alt={`QR code for ${link.short_url}`}
                width={256}
                height={256}
                className="h-64 w-64"
              />
            </div>
          </div>
        </div>
      )}

      {showHealth && (
        <div className="mt-5 border-b border-slate-800 pb-5">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium text-white">
              Link Health
            </p>

            {health && (
              <span
                className={`text-sm font-semibold capitalize ${healthStatusClass(
                  health.status,
                )}`}
              >
                ● {health.status}
              </span>
            )}
          </div>

          {healthLoading && (
            <p className="mt-4 text-sm text-slate-500">
              Checking destination…
            </p>
          )}

          {healthError && (
            <p className="mt-4 text-sm text-red-400">
              {healthError}
            </p>
          )}

          {health && !healthLoading && (
            <div className="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div className="rounded-lg border border-slate-800 bg-slate-950 p-4">
                <p className="text-xs text-slate-500">
                  HTTP status
                </p>

                <p className="mt-1 text-lg font-semibold text-white">
                  {health.http_status ?? "—"}
                </p>
              </div>

              <div className="rounded-lg border border-slate-800 bg-slate-950 p-4">
                <p className="text-xs text-slate-500">
                  Response time
                </p>

                <p className="mt-1 text-lg font-semibold text-white">
                  {health.response_time_ms} ms
                </p>
              </div>

              <div className="rounded-lg border border-slate-800 bg-slate-950 p-4">
                <p className="text-xs text-slate-500">
                  Redirects
                </p>

                <p className="mt-1 text-lg font-semibold text-white">
                  {health.redirects}
                </p>
              </div>

              <div className="rounded-lg border border-slate-800 bg-slate-950 p-4">
                <p className="text-xs text-slate-500">
                  HTTPS
                </p>

                <p className="mt-1 text-lg font-semibold text-white">
                  {health.https ? "Yes" : "No"}
                </p>
              </div>
            </div>
          )}

          {health?.final_url && (
            <div className="mt-4">
              <p className="text-xs text-slate-500">
                Final URL
              </p>

              <p className="mt-1 truncate text-sm text-slate-300">
                {health.final_url}
              </p>
            </div>
          )}

          {health?.error && (
            <div className="mt-4 rounded-lg border border-red-900 bg-red-950/30 p-4">
              <p className="text-xs text-red-400">
                Health check error
              </p>

              <p className="mt-1 text-sm text-red-300">
                {health.error}
              </p>
            </div>
          )}

          {health && (
            <p className="mt-4 text-xs text-slate-600">
              Checked {formatDate(health.checked_at)}
            </p>
          )}

          {health && (
            <button
              onClick={loadHealth}
              disabled={healthLoading}
              className="mt-4 rounded-lg border border-slate-700 px-4 py-2 text-sm text-white disabled:opacity-50"
            >
              {healthLoading
                ? "Checking…"
                : "Refresh Health"}
            </button>
          )}

          <div className="mt-6 border-t border-slate-800 pt-5">
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-white">
                Health History
              </p>

              {!healthHistoryLoading &&
                healthHistory.length > 0 && (
                  <span className="text-xs text-slate-500">
                    {healthHistory.length} check
                    {healthHistory.length === 1
                      ? ""
                      : "s"}
                  </span>
                )}
            </div>

            {healthHistoryLoading && (
              <p className="mt-4 text-sm text-slate-500">
                Loading history…
              </p>
            )}

            {healthHistoryError && (
              <p className="mt-4 text-sm text-red-400">
                {healthHistoryError}
              </p>
            )}

            {!healthHistoryLoading &&
              !healthHistoryError &&
              healthHistory.length === 0 && (
                <p className="mt-4 text-sm text-slate-500">
                  No health checks recorded yet.
                </p>
              )}

            {!healthHistoryLoading &&
              !healthHistoryError &&
              healthHistory.length > 0 && (
                <div className="mt-4 space-y-3">
                  {healthHistory.map(
                    (check, index) => (
                      <div
                        key={`${check.checked_at}-${index}`}
                        className="rounded-lg border border-slate-800 bg-slate-950 p-4"
                      >
                        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                          <div className="flex items-center gap-3">
                            <span
                              className={`text-sm font-semibold capitalize ${healthStatusClass(
                                check.status,
                              )}`}
                            >
                              ● {check.status}
                            </span>

                            <span className="text-xs text-slate-500">
                              HTTP{" "}
                              {check.http_status ??
                                "—"}
                            </span>
                          </div>

                          <span className="text-xs text-slate-600">
                            {formatDate(
                              check.checked_at,
                            )}
                          </span>
                        </div>

                        <div className="mt-3 grid gap-3 sm:grid-cols-3">
                          <div>
                            <p className="text-xs text-slate-600">
                              Response
                            </p>

                            <p className="mt-1 text-sm text-slate-300">
                              {
                                check.response_time_ms
                              }{" "}
                              ms
                            </p>
                          </div>

                          <div>
                            <p className="text-xs text-slate-600">
                              Redirects
                            </p>

                            <p className="mt-1 text-sm text-slate-300">
                              {check.redirects}
                            </p>
                          </div>

                          <div>
                            <p className="text-xs text-slate-600">
                              HTTPS
                            </p>

                            <p className="mt-1 text-sm text-slate-300">
                              {check.https
                                ? "Yes"
                                : "No"}
                            </p>
                          </div>
                        </div>

                        {check.final_url && (
                          <div className="mt-3">
                            <p className="text-xs text-slate-600">
                              Final URL
                            </p>

                            <p className="mt-1 truncate text-xs text-slate-400">
                              {check.final_url}
                            </p>
                          </div>
                        )}

                        {check.error && (
                          <div className="mt-3 rounded-lg border border-red-900 bg-red-950/30 p-3">
                            <p className="text-xs text-red-300">
                              {check.error}
                            </p>
                          </div>
                        )}
                      </div>
                    ),
                  )}
                </div>
              )}
          </div>
        </div>
      )}

      <div className="mt-5 flex flex-wrap gap-3">
        <a
          href={`${API_URL}/${link.short_code}`}
          target="_blank"
          rel="noreferrer"
          className="rounded-lg bg-white px-4 py-2 text-sm text-slate-950"
        >
          Open
        </a>

        <button
          onClick={() =>
            setShowAnalytics(
              (visible) => !visible,
            )
          }
          className="rounded-lg border border-slate-700 px-4 py-2 text-sm text-white"
        >
          {showAnalytics
            ? "Hide Analytics"
            : "Analytics"}
        </button>

        <button
          onClick={toggleHealth}
          className="rounded-lg border border-slate-700 px-4 py-2 text-sm text-white"
        >
          {showHealth
            ? "Hide Health"
            : "Link Health"}
        </button>

        <button
          onClick={() =>
            setShowQRCode(
              (visible) => !visible,
            )
          }
          className="rounded-lg border border-slate-700 px-4 py-2 text-sm text-white"
        >
          {showQRCode
            ? "Hide QR"
            : "QR Code"}
        </button>

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

      {showAnalytics && (
        <div className="mt-5 space-y-6 border-t border-slate-800 pt-5">
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-400">
              Analytics
            </h3>

            <div className="mt-4 grid gap-4 sm:grid-cols-3">
              <div className="rounded-lg border border-slate-800 bg-slate-950 p-4">
                <p className="text-xs text-slate-500">
                  Total clicks
                </p>

                <p className="mt-1 text-xl font-semibold text-white">
                  {analytics?.total_clicks ?? 0}
                </p>
              </div>

              <div className="rounded-lg border border-slate-800 bg-slate-950 p-4">
                <p className="text-xs text-slate-500">
                  First click
                </p>

                <p className="mt-1 text-sm text-slate-200">
                  {formatDate(
                    analytics?.first_click,
                  )}
                </p>
              </div>

              <div className="rounded-lg border border-slate-800 bg-slate-950 p-4">
                <p className="text-xs text-slate-500">
                  Last click
                </p>

                <p className="mt-1 text-sm text-slate-200">
                  {formatDate(
                    analytics?.last_click,
                  )}
                </p>
              </div>
            </div>
          </div>

          <AnalyticsGroup
            title="Referrers"
            items={analytics?.referrers ?? []}
          />

          <AnalyticsGroup
            title="Browsers"
            items={analytics?.browsers ?? []}
          />

          <AnalyticsGroup
            title="Devices"
            items={analytics?.devices ?? []}
          />

          <AnalyticsGroup
            title="Client IPs"
            items={analytics?.client_ips ?? []}
          />
        </div>
      )}
    </div>
  );
}

interface AnalyticsGroupProps {
  title: string;
  items: {
    name?: string;
    referrer?: string;
    ip?: string;
    clicks: number;
  }[];
}

function AnalyticsGroup({
  title,
  items,
}: AnalyticsGroupProps) {
  return (
    <div>
      <h4 className="text-sm font-semibold text-white">
        {title}
      </h4>

      {items.length === 0 ? (
        <p className="mt-3 text-sm text-slate-500">
          No data yet.
        </p>
      ) : (
        <div className="mt-3 space-y-2">
          {items.map((item) => (
            <div
              key={`${title}-${item.name ?? item.referrer ?? item.ip}`}
              className="flex justify-between rounded-lg border border-slate-800 bg-slate-950 px-4 py-3 text-sm"
            >
              <span className="truncate text-slate-300">
                {item.name ?? item.referrer ?? item.ip}
              </span>

              <span className="ml-4 shrink-0 text-slate-500">
                {item.clicks}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}