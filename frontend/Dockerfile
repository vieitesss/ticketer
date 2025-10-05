FROM cgr.dev/chainguard/wolfi-base AS builder

WORKDIR /app

RUN apk add --no-cache nodejs yarn

COPY package.json yarn.lock ./

RUN yarn install --frozen-lockfile

COPY . .

RUN yarn build

FROM cgr.dev/chainguard/wolfi-base

WORKDIR /app

RUN apk add --no-cache nodejs

ENV NODE_ENV=production

COPY --chown=nonroot:nonroot --from=builder /app/.next/standalone ./
COPY --chown=nonroot:nonroot --from=builder /app/.next/static ./.next/static
COPY --chown=nonroot:nonroot --from=builder /app/public ./public

EXPOSE 3000

USER nonroot

CMD ["node", "server.js"]
