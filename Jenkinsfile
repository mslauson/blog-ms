pipeline {
    stages {
        stage('Compile') {
            steps {
                sh 'go build'
            }
        }
        stage('Test') {
            environment {
                CODECOV_TOKEN = credentials('codecov_token')
            }
            steps {
                sh 'go test ./...
                /* sh "curl -s https://codecov.io/bash | bash -s -" */
            }
        }
        stage('Code Analysis') {
            steps {
                sh 'curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b $GOPATH/bin v1.12.5'
                sh 'golangci-lint run'
            }
        }
        /* stage('publish') { */
        /*     when { */
        /*         buildingTag() */
        /*     } */
        /*     environment { */
        /*         GITHUB_TOKEN = credentials('github_token') */
        /*     } */
        /*     steps { */
        /*         sh 'curl -sL https://git.io/goreleaser | bash' */
        /*     } */
        /* } */
    }
}
