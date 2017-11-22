FROM scratch
ADD gcloud-directory-service /gcloud-directory-service
ENTRYPOINT ["/gcloud-directory-service"]
CMD ["--help"]