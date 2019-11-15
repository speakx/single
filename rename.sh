#!/bin/bash

# 通过此脚本直接将 go-simple 重命名为自己的服务工程
repository=$1
oldrepository=${PWD##*/}
if [ "$repository" = "" ] ;then
    echo "没有输入新的項目名"
    exit
fi
echo "項目名重命名 $oldrepository -> $repository"

# 1, 修改go代码中的所有当前mod的引用
function rename_go() {
	files=`ls`
    for file in ${files[@]};do
		if [ -d "$(pwd)/$file" ] ;then
            cd ./$file
            rename_go
            cd ../
        else
            if [[ $file == *.go ]] ;then
                if [[ `cat $(pwd)/$file | grep "\"$oldrepository/"` != "" ]]; then
                    rm -rf $file.tmp
                    sed "s/\"$oldrepository\//\"$repository\//g" $file >> $file.tmp
                    rm -rf $file
                    mv $file.tmp $file
                fi

                if [[ `cat $(pwd)/$file | grep "pb$oldrepository"` != "" ]]; then
                    rm -rf $file.tmp
                    sed "s/pb$oldrepository/pb$repository/g" $file >> $file.tmp
                    rm -rf $file
                    mv $file.tmp $file
                fi
            fi
        fi
	done
}

# step1 重命名import
cd ./src
    rename_go
cd ../

# step2 重命名文件夹
cd ../
mv $oldrepository $repository

# step3 go验证一下
cd $repository
sh go.sh