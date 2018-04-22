go install ./...
rm ./bin/main
cp ~/go/bin/main ./bin
sudo docker build -t r74anand/logreader .