import CreateLinkForm from "@/components/links/CreateLinkForm";
import Dashboard from "@/components/links/Dashboard";

export default function Home() {
  return (
    <main className="min-h-screen bg-slate-950 text-white">
      <div className="mx-auto flex min-h-screen max-w-6xl flex-col px-6 py-8">
        <header className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 overflow-hidden rounded-lg">
              <img
                src="/linkforge-logo.png"
                alt="LinkForge"
                className="h-full w-full object-cover"
              />
            </div>
            
            <span className="text-xl font-semibold tracking-tight">
            LinkForge
            </span>
          </div>

          <div className="text-sm text-slate-400">
            Self-Hosted URL shortener
          </div>
        </header>

        <section className="flex flex-1 items-center justify-center py-20">
          <div className="w-full max-w-3xl text-center">
            <div className="mb-6 inline-flex rounded-full border border-slate-800 bg-slate-900 px-4 py-2 text-sm text-slate-300">
              Fast · Secure · Self-Hosted
            </div>

            <h1 className="text-5xl font-bold tracking-tight sm:text-6xl">
              Turn long URLs into
              <span className="block text-slate-400">
                useful short links.
              </span>
            </h1>

            <p className="mx-auto mt-6 max-w-2xl text-lg leading-8 text-slate-400">
              Create short links with custom aliases, expiration,
              destination validation, QR codes, and link analytics.
            </p>

            <CreateLinkForm />

            <Dashboard />

            <div className="mt-6 flex flex-wrap justify-center gap-x-6 gap-y-2 text-sm text-slate-500">
              <span>✓ HTTPS</span>
              <span>✓ Valid domain</span>
              <span>✓ Safe destination</span>
              <span>✓ QR codes</span>
            </div>
          </div>
        </section>

        <footer className="border-t border-slate-900 py-6 text-center text-sm text-slate-600">
          LinkForge
        </footer>
      </div>
    </main>
  );
}