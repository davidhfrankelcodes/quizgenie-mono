# ── Stage 1: Build the Angular app ─────────────────────────────────────────
FROM node:22-alpine AS build

WORKDIR /app

# 1) Install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# 2) Copy source & build
COPY . .
RUN npm run build --configuration=production

# ── Stage 2: Serve with nginx ──────────────────────────────────────────────
FROM nginx:alpine

# 1) Remove default welcome page
RUN rm -rf /usr/share/nginx/html/*

# 2) Add our nginx config
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# 3) Copy *just* the browser build output into nginx’s html root
COPY --from=build /app/dist/quizgenie-frontend/browser/ /usr/share/nginx/html/

# 4) Expose & run
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
