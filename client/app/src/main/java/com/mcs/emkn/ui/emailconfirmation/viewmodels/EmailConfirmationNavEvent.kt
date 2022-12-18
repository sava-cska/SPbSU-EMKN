package com.mcs.emkn.ui.emailconfirmation.viewmodels

sealed class EmailConfirmationNavEvent {
    object ContinueConfirmation : EmailConfirmationNavEvent()
}