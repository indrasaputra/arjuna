FROM nginx:1.27.0

COPY ../blueprint/index.html /usr/share/nginx/html/index.html
COPY ../openapiv2/arjuna.swagger.json /usr/share/nginx/html/api/arjuna.swagger.json
