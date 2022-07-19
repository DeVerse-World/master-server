cp ../conf/prod/ms/.env ../dotenv_file
cp ../conf/prod/ms/config.yml ../config_temp.yml
cd ..
go build main.go
rsync -avz -e ssh * root@206.189.159.101:/root/repos/deploy/main-server-prod