package com.mcs.emkn.ui.auth.viewmodels

sealed class SignInNavEvent {
    object ContinueSignIn : SignInNavEvent()
}