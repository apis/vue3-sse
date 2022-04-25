FROM node:lts-alpine as builder

WORKDIR /app

# copy both 'package.json' and 'package-lock.json' (if available)
COPY package*.json ./

# install project dependencies
RUN npm install

# copy project files and folders to the current working directory (i.e. 'app' folder)
COPY babel.config.js ./
COPY public ./public
COPY src ./src

# build app for production with minification
RUN npm run build

FROM caddy:2

WORKDIR /app

COPY --from=builder /app/dist dist

WORKDIR /

COPY ./Caddyfile /etc/caddy/Caddyfile

CMD [ "caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile" ]