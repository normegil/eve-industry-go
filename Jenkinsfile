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

def builtImage

pipeline {
    agent none
    environment {
        XDG_CACHE_HOME = '/tmp/.cache'
    }
    stages {
        stage('Validate code') {
            parallel {
                stage('Lint') {
                    agent any
                    tools {
                        go '1.16.3'
                    }
                    steps {
                        sh './setup-dev-env.sh'
                        sh 'bin/golangci-lint run'
                    }
                }
                stage('Test') {
                    agent any
                    tools {
                        go '1.16.3'
                    }
                    steps {
                        sh 'go test ./...'
                    }
                }
            }
        }
        stage('Build code') {
            agent any
            tools {
                go '1.16.3'
            }
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
            agent {
                label 'docker-build'
            }
            steps {
                script {
                    builtImage = docker.build("normegil/eve-industry:${env.BUILD_ID}")
                }
            }
        }
        stage('Acceptance test') {
            agent {
                label 'docker-build'
            }
            tools {
                go '1.16.3'
            }
            steps {
                script {
                    builtImage.withRun('-p 18080:18080') {
                        sh 'go test --tags=acceptance ./...'
                    }
                }
            }
        }
        stage('Publish artefacts') {
            parallel {
                stage('Publish docker') {
                    agent {
                        label 'docker-build'
                    }
                    steps {
                        script {
                            docker.withRegistry('https://index.docker.io/v1/', 'DockerHub') {
                                builtImage.push('latest')
                            }
                        }
                    }
                }
                stage('Publish binaries') {
                    agent any
                    tools {
                        go '1.16.3'
                    }
                    steps {
                        withCredentials([usernamePassword(credentialsId: 'GithubToken', usernameVariable: 'GITHUB_USER', passwordVariable: 'GITHUB_TOKEN')]) {
                            script {
                                sh 'wget https://github.com/github-release/github-release/releases/download/v0.10.0/linux-amd64-github-release.bz2 -O bin/github-release.bz2'
                                sh 'bzip2 -d bin/github-release.bz2 || true'
                                sh 'chmod +x bin/github-release'

                                sh 'bin/github-release delete --user ${GITHUB_USER} --repo eve-industry-go --tag latest || true'
                                sh 'bin/github-release release --user ${GITHUB_USER} --repo eve-industry-go --tag latest --name latest'
                                linuxBuildTargets.each { target ->
                                    sh "bin/github-release upload --user ${GITHUB_USER} --repo eve-industry-go --tag latest --name eve-industry-linux-${target} --file eve-industry-linux-${target}"
                                }
                                windowsBuildTargets.each { target ->
                                    sh "bin/github-release upload --user ${GITHUB_USER} --repo eve-industry-go --tag latest --name eve-industry-windows-386 --file eve-industry-windows-${target}"
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            node('docker-build') {
                sh "docker rmi ${builtImage.id}"
                cleanWs()
            }
        }
    }
}