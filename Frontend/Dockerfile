FROM nginx:alpine

ENV NGINX_PORT=80

# Remove default Nginx HTML files
RUN rm -rf /usr/share/nginx/html/*

# Copy frontend files from current directory to nginx html
COPY . /usr/share/nginx/html/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
