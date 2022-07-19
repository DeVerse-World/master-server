cd ..
cp dotenv_file .env
cp config_temp.yml config.yml
go build main.go
./migrate.sh up
sudo service deversems_staging stop
sudo service deversems_staging start
#sudo service deversems_staging status
#journalctl -u deversems_staging.service
