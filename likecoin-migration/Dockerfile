FROM --platform=$BUILDPLATFORM node:23.7-slim as builder

WORKDIR /usr/src/app
COPY . ./
RUN npm install
RUN npm run build


FROM nginx:alpine

COPY --from=builder /usr/src/app/dist /usr/share/nginx/html
