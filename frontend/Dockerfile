# ── Stage 1: Build the Angular app ─────────────────────────────────────────
FROM node:20-alpine AS build

WORKDIR /app

# 1) Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# 2) Copy the rest of the frontend code
COPY . .

# 3) Run the production build
#    (if your Angular project’s name is not “quizgenie-ui”, adjust the --configuration and output path accordingly)
RUN npm run build --configuration=production

# ── Stage 2: Serve with nginx ──────────────────────────────────────────────
FROM nginx:alpine

# Copy whatever folder(s) are under /app/dist into nginx’s html root
COPY --from=build /app/dist/* /usr/share/nginx/html
# Expose port 80
EXPOSE 80

# Run nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
