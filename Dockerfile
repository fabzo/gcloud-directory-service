FROM scratch
ADD gcloud-directory-service /gcloud-directory-service
ADD ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/gcloud-directory-service"]
CMD ["--help"]