package com.mcs.emkn.ui.changepassword.viewmodels

sealed class CommitChangePasswordNavEvent {
    object ContinueChangePassword : CommitChangePasswordNavEvent()
}