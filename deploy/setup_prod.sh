sudo certbot --nginx -d api.deverse.world
cp ../conf/prod/ms/nginx_conf /etc/nginx/sites-enabled/api.deverse.world
cp ../conf/prod/ms/service /lib/systemd/system/deversems.service
