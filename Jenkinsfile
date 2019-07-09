#!/usr/bin/env groovy

pipeline {
  environment {
    repoName = "pgexport"
    registry = "744716359313.dkr.ecr.us-west-2.amazonaws.com/toolkits/$repoName"
    ecrUrl = "https://744716359313.dkr.ecr.us-west-2.amazonaws.com"
  }
  agent any
  stages {
    stage('Building image') {
      steps {
        script {
          shortCommit = sh(returnStdout: true, script: "git log -n 1 --pretty=format:'%h'").trim()
          revNo = sh(returnStdout: true, script: "git rev-list --count HEAD").trim()
          currentBranch = env.BRANCH_NAME.replaceAll(/[^A-Za-z0-9]/, "-").toLowerCase()
          echo "Current branch $currentBranch"
          if (currentBranch == "master") {
            buildno = "prod-$BUILD_NUMBER"
          } else {
            buildno = "v$revNo-$shortCommit-$currentBranch-$BUILD_NUMBER"
          }
          env.buildno = buildno
          echo "Build No: $buildno"

          // push build no tags to Gitlab
          sh("git tag -a build/$buildno -m 'Tagged by Jenkins build $currentBranch#$BUILD_NUMBER'")
          try {
            withCredentials([usernamePassword(credentialsId: 'gitlab-credential', passwordVariable: 'GIT_PASSWORD', usernameVariable: 'GIT_USERNAME')]) {
              sh("git config credential.username ${env.GIT_USERNAME}")
              sh("git config credential.helper '!echo password=\$GIT_PASSWORD; echo'")
              sh("GIT_ASKPASS=true git push origin build/$buildno")
            }
          } finally {
            sh("git config --unset credential.username")
            sh("git config --unset credential.helper")
          }

          // build image
          customImage = docker.build(registry + ":$buildno", "--build-arg BUILDNO=$buildno --build-arg REV=$shortCommit .")
          env.imageID = customImage.id
        }
      }
    }
    stage('Deploy image') {
      steps {
        script {
          docker.withRegistry(ecrUrl, 'ecr:us-west-2:aws-ecr-credentials') {
            customImage.push()
            if (env.BRANCH_NAME != "master") { customImage.push('latest') }
          }
        }
      }
    }
  }
  post {
    always {
      script {
        def imageID = "${env.imageID}"
        if (imageID != "") {
          sh "docker rmi ${imageID}"
        }
      }
    }
    success {
      script {
        def info = collectChangeLogs()
        def encodedInfo = java.net.URLEncoder.encode(info, "UTF-8")
        httpRequest(url: "http://buildtracking.bluex.ai/component", httpMode: "POST", responseHandle: "NONE", contentType: "APPLICATION_FORM", requestBody: "component_name=${env.JOB_NAME}&build_number=${env.buildno}&build_info=$encodedInfo")
        slackSend(color:'#00FF00', message:"Build completed on [${env.JOB_NAME}], Build:${env.buildno}, Duration:${currentBuild.duration/1000}s (<${env.BUILD_URL}|Open>). Info:\n$info")
      }
    }
    failure {
      script {
        def info = collectChangeLogs()
        slackSend(color:'#FF0000', message:"⚠️⚠️⚠️ Build failed on [${env.JOB_NAME}], Build:${env.buildno}, Duration:${currentBuild.duration/1000}s (<${env.BUILD_URL}|Open>). Info:\n$info")
      }
    }
  }
}

@NonCPS
def collectChangeLogs() {
  if (currentBuild.changeSets.size() == 0) {
    return "[$GIT_COMMIT]"
  }
  def info = ""
  def changeLogSets = currentBuild.changeSets
  for (int i = 0; i < changeLogSets.size(); i++) {
    def entries = changeLogSets[i].items
    for (int j = 0; j < entries.length; j++) {
      def entry = entries[j]
      echo "${entry.commitId} by ${entry.author} on ${new Date(entry.timestamp)}: ${entry.msg}\n"
      info = info + "[${entry.commitId.substring(0,6)}](${entry.author}) ${entry.msg}\n"
    }
  }
  return info
}
