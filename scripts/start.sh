#!bin/sh

# .envファイルを読み込んで変数として参照できるようにする
. ./.env

export GOOGLE_APPLICATION_CREDENTIALS=$GOOGLE_APPLICATION_CREDENTIALS

echo $GOOGLE_APPLICATION_CREDENTIALS

go run main.go
