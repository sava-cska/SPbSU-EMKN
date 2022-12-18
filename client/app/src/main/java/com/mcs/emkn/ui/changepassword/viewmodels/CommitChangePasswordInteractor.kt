package com.mcs.emkn.ui.changepassword.viewmodels

import kotlinx.coroutines.flow.Flow

interface CommitChangePasswordInteractor {
    val errors: Flow<CommitChangePasswordError>

    val navEvents: Flow<CommitChangePasswordNavEvent>

    fun submitNewPassword(password: String)
}