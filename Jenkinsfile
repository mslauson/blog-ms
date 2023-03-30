pipeline {
      agent any
	environment {
	  GT_CREDS= credentials('Gitea')
	}
    tools {
        go 'main'
    }
    stages {
		stage('Prepare'){
		  steps {
		sh 	'git config --global url."https://GT_CREDS_USR:$GT_CREDS_PSW@gitea.slauson.io".insteadOf "https://gitea.slauson.io"'
			}
		  }
        stage('Compile') {
            steps {
                sh 'go build'
            }
        }
        stage('Test') {
            steps {
                sh 'go test ./...'
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
