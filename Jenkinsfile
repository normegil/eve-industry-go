pipeline {
    agent { docker { image 'golang:1.16.3' } }
    stages {
        stage('lint') {
            steps {
                sh './setup-dev-env.sh'
                withEnv(["XDG_CONFIG_HOME=/tmp"]) {
                    sh 'bin/golangci-lint --version'
                }
            }
        }
    }
}