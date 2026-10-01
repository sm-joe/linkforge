# syntax=docker/dockerfile:1.7

# ---------- Dependencies ----------
FROM node:22-alpine AS deps

WORKDIR /app

COPY apps/web/package.json apps/web/package-lock.json ./

RUN npm ci

# ---------- Build ----------
FROM node:22-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY apps/web ./

ENV NEXT_TELEMETRY_DISABLED=1

RUN npm run build

# ---------- Runtime ----------
FROM node:22-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1

RUN addgroup -S linkforge \
    && adduser -S -G linkforge linkforge

COPY --from=builder --chown=linkforge:linkforge /app/.next/standalone ./
COPY --from=builder --chown=linkforge:linkforge /app/.next/static ./.next/static
COPY --from=builder --chown=linkforge:linkforge /app/public ./public

USER linkforge

EXPOSE 3000

ENV PORT=3000 \
    HOSTNAME=0.0.0.0

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD node -e "fetch('http://127.0.0.1:3000').then(r => process.exit(r.ok ? 0 : 1)).catch(() => process.exit(1))"

CMD ["node", "server.js"]