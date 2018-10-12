node {
    stage('Clone Repo') {
        checkout scm
    }

    withEnv(["GO111MOD=on"]) {
        stage('Build') {
            docker.image('golang:1.11.1').inside("--user=root") {
                sh 'go vet -mod vendor'
                sh 'go build -mod vendor -ldflags "-s -w" -o octaaf'
            }
        }
    }

    withEnv(["VERSION=${env.BUILD_TAG}"]) {
        stage("Package") {
            if( env.TAG_NAME.startsWith("release-") ) {
                sh "make package"
            }
        }
    }

    withEnv([
        "REPO_SERVER=repo.youkebox.be",
        "REPO_PATH=/var/vhosts/repo/octaaf/"]) {
        stage("Upload") {
            if( env.TAG_NAME.startsWith("release-") ) {
                sh "scp octaaf-*.rpm root@${REPO_SERVER}:${REPO_PATH}/packages/"
                sh "ssh root@${REPO_SERVER} 'createrepo --update ${REPO_PATH}'"
            }
        }
    }

    withEnv(["REPO_SERVER=repo.youkebox.be"]) {
        if( env.TAG_NAME.startsWith("release-") ) {
            stage('Deploy') {
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
