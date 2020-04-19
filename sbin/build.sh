#app项目的编译脚本
curpwd=`pwd`
cd $HOME/app/src/main
go build -o app
mv app $HOME/app/bin
cd $curpwd
