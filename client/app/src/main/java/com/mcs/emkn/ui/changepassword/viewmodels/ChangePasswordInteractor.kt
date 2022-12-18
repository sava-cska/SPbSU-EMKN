package com.mcs.emkn.ui.changepassword.viewmodels

import kotlinx.coroutines.Deferred
import kotlinx.coroutines.flow.Flow

interface ChangePasswordInteractor {
    val errors: Flow<ChangePasswordError>

    val navEvents: Flow<ChangePasswordNavEvent>

    /**
     * Отдаёт секунды таймера, как только он устанавливается
     */
    val timerFlow: Flow<Long>

    fun validateCode(code: String)

    fun sendAnotherCode()

    fun loadTimerAsync() : Deferred<Long?>
}