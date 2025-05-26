# Use official Nginx image to serve static files
FROM nginx:alpine

# Set working directory
WORKDIR /usr/share/nginx/html

# Copy entire Frontend directory contents to Nginx html directory
COPY . .

# Replace default Nginx index.html with a symlink to LandingPage/index.html
RUN rm -f index.html && ln -s LandingPage/index.html index.html

# Expose port 80 to Docker host
EXPOSE 80

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
