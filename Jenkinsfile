pipeline {

    agent any

    environment {
        REPO_SERVER = 'repo.youkebox.be'
        REPO_PATH   = "/var/vhosts/repo/${env.GIT_BRANCH}"
        NAME        = 'octaaf'
        VERSION     = '0.2.1'
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
            steps {
                sh "make package --environment-overrides BUILD_NO=${env.BUILD_NUMBER}"
            }
        }

        stage('Upload') {
            when {
                allOf {
                    expression { BRANCH_NAME ==~ /(master|development)/ }
                    expression { env.CHANGE_ID == null  }
                }
            }
            steps {
                sh "scp octaaf-*.rpm root@${REPO_SERVER}:${REPO_PATH}/packages/"
                ssh root@${REPO_SERVER} '\\
                    cd ${REPO_PATH}/packages/ \\
                    && rm -rf \$(ls ${REPO_PATH}/packages/ -1t | grep ${NAME}-${VERSION} | tail -n +4) \\
                    && createrepo --update ${REPO_PATH}'
                """
            }
        }

        stage('Deploy') {
            agent any
            when {
                allOf {
                    expression { BRANCH_NAME == "master" }
                    expression { env.CHANGE_ID == null  }
                    expression { tag != null }
                }
            }
            steps {
                sh """
                ssh root@${REPO_SERVER} '\\
                    yum -y install https://repo.youkebox.be/master/packages/${NAME}-${VERSION}-${env.BUILD_NUMBER}.${ARCH}.rpm \\
                    && systemctl restart octaaf'
                """
            }
        }
    }
}
