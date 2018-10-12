node {
    stage('Clone Repo') {
        checkout scm
    }

    env.GIT_TAG = gitTagName()

    withEnv(["GO111MOD=on"]) {
        stage('Build') {
            docker.image('golang:1.11.1').inside("--user=root") {
                sh 'go vet -mod vendor'
                sh 'go build -mod vendor -ldflags "-s -w" -o octaaf'
            }
        }
    }

    if( env.GIT_TAG.startsWith("release-") ) {
        withEnv([
            "VERSION=${env.GIT_TAG}"]) {
             stage("Package") {
                sh "make package"
            }
        }

        withEnv([
            "REPO_SERVER=repo.youkebox.be",
            "REPO_PATH=/var/vhosts/repo/octaaf/"]) {
            stage("Upload") {
                sh "scp octaaf-*.rpm root@${REPO_SERVER}:${REPO_PATH}/packages/"
                sh "ssh root@${REPO_SERVER} 'createrepo --update ${REPO_PATH}'"
            }
        }

        withEnv(["REPO_SERVER=repo.youkebox.be"]) {
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
