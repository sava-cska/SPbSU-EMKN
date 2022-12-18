package com.mcs.emkn.core

import android.os.Bundle

interface Router {
    fun back()

    fun goToEmailConfirmationScreen()

    fun goToForgotPasswordScreen()

    fun goToChangePasswordConfirmationScreen()

    fun goToCommitChangePasswordScreen()

    fun goToRegistrationNavGraph()

    fun goToUserNavGraph()

    fun goToProfile()

    fun goToCoursePage(bundle: Bundle)

    fun goToCourseHomeworks(bundle: Bundle)
}