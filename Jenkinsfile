pipeline {
    agent any

    environment {
        REPO_SERVER    = 'repo.youkebox.be'
        REPO_PATH      = "/var/vhosts/repo/${env.GIT_BRANCH}"
        NAME           = 'octaaf'
        OCTAAF_VERSION = '0.1.1'
        DESCRIPTION    = 'A Go Telegram bot'
        ARCH           = 'x86_64'
    }

    stages {
       stage('Build') {
            steps {
                sh 'make'
            }
            post {
                always {
                    sh 'sudo rm -rf vendor'
                }
            }
        }

        stage('Package') {
            steps {
                sh "make package --environment-overrides BUILD_NO=${env.BUILD_NUMBER}"
            }
        }

        stage('Upload') {
            steps {
                sh "scp octaaf-*.rpm root@${REPO_SERVER}:${REPO_PATH}/packages/"
                sh "ssh root@${REPO_SERVER} 'cd ${REPO_PATH}/packages/ && rm -rf \$(ls ${REPO_PATH}/packages/ -1t | tail -n +4)'"
                sh "ssh root@${REPO_SERVER} 'createrepo --update ${REPO_PATH}'"
            }
        }

        stage('Deploy') {
            when {
                expression { return ${env.GIT_BRANCH} == 'master' }
            }
            steps {
                 sh "ssh root@${REPO_SERVER} 'yum makecache; yum update octaaf'"
                 sh "ssh root@${REPO_SERVER} 'systemctl restart octaaf'"
            }
        }
    }
}