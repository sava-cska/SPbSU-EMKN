package com.mcs.emkn.ui.auth.viewmodels

import com.mcs.emkn.database.entities.Credentials
import kotlinx.coroutines.Deferred
import kotlinx.coroutines.flow.Flow

interface SignInInteractor {
    val isLoadingFlow: Flow<Boolean>

    val errorsFlow: Flow<SignInError>

    val navEvents: Flow<SignInNavEvent>

    fun loadCredentialsAsync() : Deferred<Credentials?>

    fun onSignInClick(login: String, password: String)
}