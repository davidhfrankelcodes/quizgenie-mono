# Stage 1: build
FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build --configuration=production

# Stage 2: serve + inject env.js at container start
FROM nginx:alpine
ENV ALLOW_SIGNUP=false

# clean out default
RUN rm -rf /usr/share/nginx/html/*

# copy nginx config
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# copy built app
COPY --from=build /app/dist/quizgenie-frontend/browser/ /usr/share/nginx/html/
# add entrypoint script
COPY frontend-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.sh"]
