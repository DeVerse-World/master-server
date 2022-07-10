cd ..
cp dotenv_file .env
go build main.go
./migrate.sh up
sudo service deversems stop
sudo service deversems start
#journalctl -u deversems.service
