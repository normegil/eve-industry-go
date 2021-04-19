def linuxBuildTargets = [
    "386",
    "amd64",
    "arm",
    "arm64"
]

def windowsBuildTargets = [
    "386",
    "amd64"
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
            steps {
                script {
                    def builds = [:]
                    linuxBuildTargets.each { target ->
                        builds["linux-"+target] = {
                            stage("Build Linux - ${target}") {
                                sh "GOOS=linux GOARCH=${target} go build -o eve-industry-linux-${target} ./..."
                            }
                        }
                    }
                    windowsBuildTargets.each { target ->
                        builds["windows-"+target] = {
                            stage("Build Windows - ${target}") {
                                sh "GOOS=windows GOARCH=${target} go build -o eve-industry-windows-${target} ./..."
                            }
                        }
                    }
                    parallel builds
                }
            }
        }
        stage('Build docker image') {
            agent any
            steps {
                script {
                    def img = docker.build("eve-industry:${env.BUILD_ID}")
                    img.push('latest')
                }
            }
        }
    }
}