cp ../conf/staging/ms/.env ../dotenv_file
cp ../conf/staging/ms/config.yml ../config.yml
cd ..
go build main.go
rsync -avz -e ssh * root@206.189.159.101:/root/repos/deploy/main-server-staging
