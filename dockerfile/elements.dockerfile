FROM nginx:1.21

COPY ../web/index.html /usr/share/nginx/html/index.html
COPY ../openapiv2/api/v1/user.swagger.json /usr/share/nginx/html/api/user.swagger.json
