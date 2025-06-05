# ── Stage 1: Build the Angular app ─────────────────────────────────────────
FROM node:16-alpine AS build

WORKDIR /app

# Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Copy the rest of the frontend code
COPY . .

# Build the production artifacts
RUN npm run build --prod

# ── Stage 2: Serve with nginx ──────────────────────────────────────────────
FROM nginx:alpine

# Copy the compiled Angular output to nginx’s html folder.
# Adjust “quizgenie-ui” below if your dist folder has a different name.
COPY --from=build /app/dist/quizgenie-ui /usr/share/nginx/html

# If you have a custom nginx.conf, copy it here. Otherwise, the default
# will serve index.html and static assets.
# COPY nginx/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
