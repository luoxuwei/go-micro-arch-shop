pipeline {
    agent any

    stages {
        stage('pull code') {
            steps {
                git branch: 'main', credentialsId: 'github', url: 'https://github.com/luoxuwei/go-micro-arch-shop.git'
            }
        }
        
        stage('build project') {
            steps {
                sh '''
                echo "配置go环境"
                export GOROOT=/usr/local/go
                export PATH=$PATH:$GOROOT/bin
                go env -w GO111MODULE=on
                go env -w GOPROXY=https://goproxy.io

                echo "准备目录"
                chmod -R 777 ${srv_or_api}/${project_name}/ 
                mkdir -vp target/${project_name}/
                cp ${srv_or_api}/${project_name}/config-pro.yaml target/${project_name}/config-pro.yaml
                cp start.sh target/start.sh

                echo "go build"
                cd ${srv_or_api}
                go build -o ../target/${project_name}_main ${project_name}/main.go
                echo "构建完成"'''
            }
        }
        
        stage('deploy project') {
            steps {
                sshPublisher(publishers: [sshPublisherDesc(configName: '192.168.139.9', transfers: [sshTransfer(cleanRemote: false, excludes: '', execCommand: 'chmod +x /docker/go/${project_name}/start.sh && cd /docker/go/${project_name}/ && ./start.sh ${project_name}_main', execTimeout: 120000, flatten: false, makeEmptyDirs: false, noDefaultExcludes: false, patternSeparator: '[, ]+', remoteDirectory: '/docker/go/${project_name}/', remoteDirectorySDF: false, removePrefix: 'target', sourceFiles: 'target/**')], usePromotionTimestamp: false, useWorkspaceInPromotion: false, verbose: false)])
            }
        }
    }
}
