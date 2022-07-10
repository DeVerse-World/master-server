sudo certbot --nginx -d api.staging.deverse.world
cp ../conf/staging/ms/nginx_conf /etc/nginx/sites-enabled/api.staging.deverse.world
cp ../conf/staging/ms/service /lib/systemd/system/deversems_staging.service