FROM nginx:1.27.0

COPY ../blueprint/nginx.conf /etc/nginx/conf.d/default.conf
COPY ../blueprint/index.html /usr/share/nginx/html/index.html
COPY ../openapiv2/arjuna.swagger.yaml /usr/share/nginx/html/api/arjuna.swagger.yaml
