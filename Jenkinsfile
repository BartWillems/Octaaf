pipeline {

    agent any

    environment {
        REPO_SERVER = 'repo.youkebox.be'
        REPO_PATH   = "/var/vhosts/repo/octaaf/"
        NAME        = 'octaaf'
        VERSION     = "${TAG_NAME}"
        DESCRIPTION = 'A Go Telegram bot'
        ARCH        = 'x86_64'
        GO111MODULE = 'on'
    }

    stages {
        stage('Build') {
            agent {
                docker { 
                    image 'golang:1.11.1'
                    args '--user=root'
                }
            }
            steps {
                sh 'go test -mod vendor'
                sh 'go build -mod vendor -ldflags "-s -w" -o octaaf'
                stash includes: 'octaaf', name: 'octaaf'
            }
        }

        stage('Package') {
            when { buildingTag() }
            steps {
                unstash 'octaaf'
                sh "make package"
            }
        }

        stage('Upload') {
            when { buildingTag() }
            steps {
                sh "scp octaaf-*.rpm root@${REPO_SERVER}:${REPO_PATH}/packages/"
                sh "ssh root@${REPO_SERVER} 'createrepo --update ${REPO_PATH}'"
            }
        }

        stage('Deploy') {
            when { 
                allOf {
                    buildingTag()
                    tag "release-*"
                }
            }
            steps {
                sh """
                ssh root@${REPO_SERVER} '\\
                    yum makecache; yum update octaaf -y \\
                    && systemctl daemon-reload \\
                    && systemctl restart octaaf'
                """
            }
        }
    }
}
