package com.mcs.emkn.ui.emailconfirmation.viewmodels

sealed class EmailConfirmationError {
    object InvalidCode : EmailConfirmationError()
    object RegistrationExpired : EmailConfirmationError()
    object BadNetwork : EmailConfirmationError()
}