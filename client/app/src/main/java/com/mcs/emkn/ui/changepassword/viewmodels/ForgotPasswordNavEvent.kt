package com.mcs.emkn.ui.changepassword.viewmodels

sealed class ForgotPasswordNavEvent {
    object ContinueForgotPassword : ForgotPasswordNavEvent()
}