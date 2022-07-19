cd ..
cp dotenv_file .env
cp config_temp.yml config.yml
go build main.go
./migrate.sh up
sudo service deversems stop
sudo service deversems start
#journalctl -u deversems.service
