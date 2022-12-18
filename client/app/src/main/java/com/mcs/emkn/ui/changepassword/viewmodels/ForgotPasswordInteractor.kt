package com.mcs.emkn.ui.changepassword.viewmodels

import kotlinx.coroutines.flow.Flow

interface ForgotPasswordInteractor {
    val errors: Flow<ForgotPasswordError>

    val navEvents: Flow<ForgotPasswordNavEvent>

    fun onSubmitClick(credentials: String)
}