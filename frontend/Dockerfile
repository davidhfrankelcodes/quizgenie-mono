# ── Stage 1: Build the Angular app ─────────────────────────────────────────
FROM node:22-alpine AS build

WORKDIR /app

# 1) Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# 2) Copy the rest of the frontend code
COPY . .

# 3) Run the production build
RUN npm run build --configuration=production


# ── Stage 2: Serve with nginx ──────────────────────────────────────────────
FROM nginx:alpine

# Remove default nginx index (optional but recommended)
RUN rm -rf /usr/share/nginx/html/*

# Copy Angular’s “dist/quizgenie-frontend” output into nginx’s html root
COPY --from=build /app/dist/quizgenie-frontend /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Run nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
