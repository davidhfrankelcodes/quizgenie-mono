# Stage 1: build
FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build --configuration=production

# Stage 2: serve + inject env.js at container start
FROM nginx:alpine
# default ALLOW_SIGNUP (overridden by docker-compose)
ENV ALLOW_SIGNUP=false

# clean out default welcome page
RUN rm -rf /usr/share/nginx/html/*

# copy our custom nginx config
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# copy **browser/** build output into nginx’s html root
COPY --from=build /app/dist/quizgenie-frontend/browser/ /usr/share/nginx/html/

# add runtime‐env injector
COPY frontend-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.sh"]
