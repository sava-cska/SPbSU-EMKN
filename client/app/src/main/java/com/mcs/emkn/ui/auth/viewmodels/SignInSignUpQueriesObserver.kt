package com.mcs.emkn.ui.auth.viewmodels

import java.util.concurrent.atomic.AtomicBoolean
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class SignInSignUpQueriesObserver @Inject constructor() {
    val isLoading = AtomicBoolean(false)
}