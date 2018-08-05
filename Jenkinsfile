pipeline {

    agent any

    environment {
        REPO_SERVER = 'repo.youkebox.be'
        REPO_PATH   = "/var/vhosts/repo/octaaf/"
        NAME        = 'octaaf'
        VERSION     = "${TAG_NAME}"
        DESCRIPTION = 'A Go Telegram bot'
        ARCH        = 'x86_64'
    }

    stages {
        stage('Build') {
            steps {
                sh 'make build'
            }
        }

        stage('Package') {
            when { buildingTag() }
            steps {
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
