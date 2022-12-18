package com.mcs.emkn.ui.changepassword.viewmodels

sealed class ChangePasswordError {
    object InvalidCode : ChangePasswordError()
    object BadNetwork : ChangePasswordError()
    object CodeExpired : ChangePasswordError()
}