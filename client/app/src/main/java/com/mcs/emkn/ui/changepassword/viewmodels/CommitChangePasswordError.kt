package com.mcs.emkn.ui.changepassword.viewmodels

sealed class CommitChangePasswordError {
    object InvalidPassword : CommitChangePasswordError()
    object ChangeExpired : CommitChangePasswordError()
    object BadNetwork : CommitChangePasswordError()
}