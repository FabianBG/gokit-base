pipeline {
    /*
        Seleccion agente correcto para deploy
    */
    agent any

    /*
        Si el Job demora más de 30 minutos es cancelado.
    */
    options {
        timeout(time: 30 , unit: 'MINUTES')
    }

    /*
        Variables de entorno
    */
    environment {
        ENVIRONMENT           = "develop"

        IMAGE_OWNER          = "up_company"

        REPOSITORY_BASE_NAME = """${sh(
                returnStdout: true,
                script: 'git config --get remote.origin.url | cut -d "/" -f 5 | cut -d "." -f 1 | tr -d \'[[:space:]]\'').trim()}"""
        PROYECT_NAME         = """${sh(
                returnStdout: true,
                script: 'git config --get remote.origin.url | cut -d "/" -f 5 | cut -d "." -f 1 | tr -d \'[[:space:]]\'')}"""

        REGISTRY_PATH = "up-company.dev:5000"

        REGISTRY_PROJECT = """${REGISTRY_PATH}/up-registry"""

        KUBEFILE = "/var/lib/jenkins/.kube/config"

        KUBE_BASE_NAME = REPOSITORY_BASE_NAME.replace("_", "-")

        GIT_AUTH = credentials("bitbucker") 

        SONAR_SCANER_HOME = tool 'sonar-scaner-default'

    }


    stages {
        stage('Init'){
            stages{
                stage('Init: Confirm deploy to producction') {
                    when {
                        anyOf { branch 'master' }
                    }
                    steps {
                        timeout(time: 5, unit: 'MINUTES') {
                            input "Really want to deploy to production?"
                        }

                    }
                }
                stage('Init: Define Version and Variables') {
                    steps {
                        wrap([$class: 'BuildUser']) {
                            slackSend (color: '#2196f3', message: "STARTED by ${env.BUILD_USER_ID}: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
                            bitbucketStatusNotify buildState: "INPROGRESS"
                        }
                        
                        
                        timeout(time: 5, unit: 'MINUTES') {
                            script {
                                VERSION     = new Date().format( 'yyyy.MM.dd.H.m.s' )
                                if (env.BRANCH_NAME == 'master') {
                                    ENVIRONMENT = "production"
                                    VERSION = "${VERSION}-${ENVIRONMENT}"
                                } else {
                                    ENVIRONMENT = input id: 'Environment-deployment',
                                                    message: 'What environment do you want to deploy?',
                                                    ok: 'OK',
                                                    parameters: [[$class: 'ChoiceParameterDefinition', choices: ['develop', 'staging'], description: 'select one of the environments to deploy', name: 'Environment']]
                                }
                                IMAGE       = "${IMAGE_OWNER}/${PROYECT_NAME}:${VERSION}"
                            }
                        }
                    }
                }
            }
        }

        /*
            Unit Test
        */
        stage('Unit Test'){
            stages{

                stage ('Unit Test: Check') {
                    
                    steps {
                        sh "go test ./... -cover"
                        script {
                            TEST_STATUS = sh(
                                        returnStdout: true,
                                        script: 'echo $?').trim()
                            
                            echo TEST_STATUS
                        
                            if (TEST_STATUS != '0') {
                                currentBuild.result = 'ABORTED'
                                error('Error on the tests …')
                                bitbucketStatusNotify buildState: "FAILED"
                            }
                        }
                        sh "go test ./... -coverprofile cover.out; go tool cover -func cover.out"
                    }
                }

            }
        }

        /*
            Code Quality
        */
        stage('Code Quality Analysis'){
            stages{

                stage ('SonarQube: Scan') {
                    
                    steps {
                        script {
                            if (env.BRANCH_NAME == 'master') {
                             sh "echo master build avoid scan."
                            } else {
                            withSonarQubeEnv('up-sonarqube') { 
                                    // If you have configured more than one global server connection, you can specify its name
                                    sh "${SONAR_SCANER_HOME}/bin/sonar-scanner"
                                }
                            }  
                        }    
                    }
                }

            }
        }

        /*
            Build
        */
        stage('Build'){
            stages {

                stage ('Build: Publish Tag Git') {
                    when {
                        anyOf { branch 'master' }
                    }
                    steps {
                            sh """
                                git config --local credential.helper "!f() { echo username=\\$GIT_AUTH_USR; echo password=\\$GIT_AUTH_PSW; }; f"
                                git config --global user.email "fbastidas.up@gmail.com" -m 
                                git config --global user.name "Jenkins"
                                git tag -a "${VERSION}"  -m "${VERSION}"
                                git push --follow-tags
                            """
                    }
                }

                stage ('Build: Publish Docker Image') {
                    steps {
                         withDockerRegistry([credentialsId: "docker-up-registry", url: "http://${REGISTRY_PATH}/v2"]){
                                sh """docker build -t ${REGISTRY_PROJECT}/${REPOSITORY_BASE_NAME}:${VERSION} ."""
                                sh """docker push ${REGISTRY_PROJECT}/${REPOSITORY_BASE_NAME}:${VERSION}"""
                            }
                    }
                }

                stage ('Build: Prune docker') {
                    steps {
                        sh """docker image prune -af"""
                    }
                }
            }

        }

        /*
            Deployment
        */
        stage('Deployment'){
            stages{
                stage ('Deployment: K8S cluster') {
                    steps {
                        sh """ kubectl --kubeconfig=${KUBEFILE} set image deployment/${KUBE_BASE_NAME} *=${REGISTRY_PROJECT}/${REPOSITORY_BASE_NAME}:${VERSION} -n ${ENVIRONMENT}"""
                        sh """ kubectl --kubeconfig=${KUBEFILE} rollout status deployment/${KUBE_BASE_NAME}  -n ${ENVIRONMENT} """
                    }
                }
             }
        }

    }

    post {
        success {
            bitbucketStatusNotify buildState: "SUCCESSFUL"

            slackSend (color: '#00FF00', message: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")

        }

        aborted {
            bitbucketStatusNotify buildState: "FAILED"

            slackSend (color: '#FF0000', message: "ABORTED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")

        }   

        failure {
            bitbucketStatusNotify buildState: "FAILED"

            slackSend (color: '#FF0000', message: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
        }
    }
}