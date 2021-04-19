def linuxBuildTargets = [
    "386",
    "amd64",
    "arm",
    "arm64"
]

pipeline {
    agent { docker { image 'golang:1.16.3' } }
    environment {
        XDG_CACHE_HOME = '/tmp/.cache'
    }
    stages {
        stage('Validate code') {
            parallel {
                stage('Lint') {
                    steps {
                        sh './setup-dev-env.sh'
                        sh 'bin/golangci-lint run'
                    }
                }
                stage('Test') {
                    steps {
                        sh 'go test ./...'
                    }
                }
            }
        }
        stage('Build code') {
            parallel {
                stage('Linux') {
                    steps {
                        script {
                            linuxBuildTargets.each { target ->
                                withEnv["GOOS=linux", "GOARCH=" + target] {
                                    sh "go build -o eve-industry-linux-${target} ./..."
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}