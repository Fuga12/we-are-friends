# The frontend is compiled by the release workflow before the image is built.
# Keeping this image runtime-only makes production builds independent of npm
# registry/client failures and leaves a minimal static Nginx image.
FROM nginx:1.27-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY dist /usr/share/nginx/html

EXPOSE 80
