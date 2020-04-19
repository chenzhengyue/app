#app项目的编译脚本
curpwd=`pwd`
cd $HOME/app/src
go build
mv app $HOME/app/bin
cd $curpwd
