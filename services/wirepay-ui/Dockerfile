# Multi-stage build for wirepay-ui.
# TODO: pin to a specific node tag once we settle on one.

FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci || npm install
COPY . .
RUN npm run build

FROM nginx:1.27-alpine
COPY --from=build /app/.output/public /usr/share/nginx/html
COPY --from=build /app/.output/server /usr/share/nginx/html/_nuxt/server
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
