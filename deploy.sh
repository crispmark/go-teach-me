#! /bin/bash
gcloud compute firewall-rules create default-allow-http-8080 \
  --allow tcp:8080 \
  --source-ranges 0.0.0.0/0 \
  --target-tags http-server \
  --description "Allow port 8080 access to http-server"

gcloud compute instances create go-compute \
  --image-family=centos-7 \
  --image-project=centos-cloud \
  --machine-type=f1-micro \
  --metadata-from-file startup-script=startup.sh \
  --tags http-server
