pipeline {
      agent any
	environment {
	  GT_CREDS= credentials('Gitea')
	  IAM_KEY="${env.IAM_KEY}"
	  IAM_HOST="${env.IAM_HOST}"
	  IAM_PROJECT="${env.IAM_PROJECT}"
	  OAUTH_CLIENT_ID="${env.OAUTH_CLIENT_ID}"
	  OAUTH_CLIENT_SECRET="${env.OAUTH_CLIENT_SECRET}"
	  OAUTH_ADMIN_BASE="${env.OAUTH_ADMIN_BASE}"
	  OAUTH_ISSUER_BASE="${env.OAUTH_ISSUER_BASE}"
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
                withEnv(['IAM_KEY=${env.IAM_KEY}', 'IAM_HOST=${env.IAM_HOST}', 'IAM_PROJECT=${env.IAM_PROJECT}', 'OAUTH_CLIENT_ID=${env.OAUTH_CLIENT_ID}', 'OAUTH_CLIENT_SECRET=${env.OAUTH_CLIENT_SECRET}', 'OAUTH_ADMIN_BASE=${env.OAUTH_ADMIN_BASE}', 'OAUTH_ISSUER_BASE=${env.OAUTH_ISSUER_BASE}']) {
                    sh 'printenv'
                    sh 'go test ./...'
                }
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
