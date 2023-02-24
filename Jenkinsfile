pipeline {
  agent any
  stages {
    stage('pull code') {
      steps {
        git(url: 'https://github.com/malphitee/cos-backup.git', branch: 'master', credentialsId: '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')
      }
    }

  }
}