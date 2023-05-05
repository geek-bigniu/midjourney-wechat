CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/mj-wechat-bot-linux -trimpath main.go
CGO_ENABLED=0 go build -o release/mj-wechat-bot-mac -trimpath main.go
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o release/mj-wechat-bot.exe -trimpath main.go
tar -czvf release/mj-wechat-bot-all.tar.gz ./release/*