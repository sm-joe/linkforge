# syntax=docker/dockerfile:1.7

# ---------- Dependencies ----------
FROM node:26-bookworm-slim AS deps

WORKDIR /app

# Apply latest Debian security updates
RUN apt-get update \
    && apt-get upgrade -y \
    && rm -rf /var/lib/apt/lists/*

COPY apps/web/package.json apps/web/package-lock.json ./

RUN npm install -g npm@latest \
    && npm cache clean --force

RUN --mount=type=cache,id=linkforge-web-npm,target=/root/.npm \
    npm ci \
    --no-audit \
    --no-fund

# ---------- Build ----------
FROM node:26-bookworm-slim AS builder

WORKDIR /app

# Apply latest Debian security updates
RUN apt-get update \
    && apt-get upgrade -y \
    && rm -rf /var/lib/apt/lists/*

COPY --from=deps /app/node_modules ./node_modules
COPY apps/web ./

ENV NEXT_TELEMETRY_DISABLED=1

RUN --mount=type=cache,target=/root/.npm npm run build

# ---------- Runtime ----------
FROM node:26-bookworm-slim AS runner

WORKDIR /app

ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000 \
    HOSTNAME=0.0.0.0

# Apply latest Debian security updates and remove npm from runtime
RUN apt-get update \
    && apt-get upgrade -y \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /usr/local/lib/node_modules/npm \
    && rm -f /usr/local/bin/npm /usr/local/bin/npx \
    && groupadd --system linkforge \
    && useradd --system --gid linkforge --create-home linkforge

COPY --from=builder --chown=linkforge:linkforge /app/.next/standalone ./
COPY --from=builder --chown=linkforge:linkforge /app/.next/static ./.next/static
COPY --from=builder --chown=linkforge:linkforge /app/public ./public

USER linkforge

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD node -e "fetch('http://127.0.0.1:3000').then(r => process.exit(r.ok ? 0 : 1)).catch(() => process.exit(1))"

CMD ["node", "server.js"]