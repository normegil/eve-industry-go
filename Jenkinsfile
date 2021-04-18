pipeline {
    agent { docker { image 'golang:1.16.3' } }
    stages {
        stage('lint') {
            steps {
                sh './setup-dev-env.sh'
                withEnv(["XDG_CACHE_HOME=/tmp/.cache"]) {
                    sh 'bin/golangci-lint run'
                }
            }
        }
        stage('test') {
            steps {
                sh 'go test ./...'
            }
        }
    }
}