package com.mcs.emkn.ui.auth.viewmodels

sealed class SignUpError {
    object BadNetwork : SignUpError()
    object IllegalLogin : SignUpError()
    object IllegalPassword : SignUpError()
    object LoginIsNotAvailable : SignUpError()
    object EmailIsNotAvailable : SignUpError()
}
