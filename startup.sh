#!/bin/bash

DOMAIN_LIST="$(/app/nginx-setup-tool -config=/app/config.json -op=domains)"
CERT_NAME="$(/app/nginx-setup-tool -config=/app/config.json -op=certname)"
EMAIL="$(/app/nginx-setup-tool -config=/app/config.json -op=email)"
WEBROOT='/var/www/acme'

nginx_pid=
sleep_pid=

function shut_down() {
    [[ $nginx_pid ]] && kill -s SIGQUIT "$nginx_pid"
    [[ $sleep_pid ]] && kill "$sleep_pid"
    exit 0
}
trap shut_down SIGQUIT

mkdir -p $WEBROOT
mkdir -p /etc/letsencrypt/

nginx -g 'daemon off;' & nginx_pid=$!

found_certs=0
certbot certificates --domains $DOMAIN_LIST | grep "Certificate Name: $CERT_NAME" > /dev/null 2>&1 || found_certs=$?

if [ "${found_certs}" -ne 0 ]; then
  echo "### Couldn't find existing cert for all specified subdomains; creating one:"
  certbot certonly --cert-name $CERT_NAME --webroot -w $WEBROOT -m $EMAIL -n --agree-tos -d $DOMAIN_LIST

  if [[ ! -f /etc/letsencrypt/ssl-dhparams.pem ]]; then
      openssl dhparam -out /etc/letsencrypt/ssl-dhparams.pem 2048
  fi
fi

/app/nginx-setup-tool -config=/app/config.json -op=write

nginx -s reload

while :
do
    echo "### assessing certbot renewal..."
    certbot renew --webroot -w $WEBROOT
    sleep 12h & sleep_pid=$!
    wait $sleep_pid
done

