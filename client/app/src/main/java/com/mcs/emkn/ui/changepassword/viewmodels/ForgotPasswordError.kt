package com.mcs.emkn.ui.changepassword.viewmodels

sealed class ForgotPasswordError {
    object BadNetwork : ForgotPasswordError()
    object InvalidEmail : ForgotPasswordError()
}
