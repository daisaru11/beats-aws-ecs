FROM nginx:alpine
RUN rm /var/log/nginx/*.log # Write logs to file
