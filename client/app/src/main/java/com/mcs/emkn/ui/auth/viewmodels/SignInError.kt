package com.mcs.emkn.ui.auth.viewmodels

sealed class SignInError {
    object IncorrectLoginOrPassword : SignInError()
    object BadNetwork : SignInError()
}
