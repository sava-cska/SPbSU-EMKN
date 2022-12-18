package com.mcs.emkn.ui.changepassword.viewmodels

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.haroldadmin.cnradapter.NetworkResponse
import com.mcs.emkn.database.Database
import com.mcs.emkn.database.entities.ChangePasswordCommit
import com.mcs.emkn.network.Api
import com.mcs.emkn.network.dto.request.RevalidateCredentialsDto
import com.mcs.emkn.network.dto.request.ValidateChangePasswordRequestDto
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Deferred
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.async
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.launch
import java.util.concurrent.atomic.AtomicBoolean
import javax.inject.Inject

@HiltViewModel
class ChangePasswordViewModel @Inject constructor(
    private val api: Api,
    private val db: Database,
) : ViewModel(), ChangePasswordInteractor {
    override val errors: Flow<ChangePasswordError>
        get() = _errors
    override val navEvents: Flow<ChangePasswordNavEvent>
        get() = _navEvents
    override val timerFlow: Flow<Long>
        get() = _timerFlow

    private val _errors = MutableSharedFlow<ChangePasswordError>()
    private val _navEvents = MutableSharedFlow<ChangePasswordNavEvent>()
    private val _timerFlow = MutableSharedFlow<Long>()
    private val isSendingCodeAtomic = AtomicBoolean(false)
    private val isValidatingCodeAtomic = AtomicBoolean(false)

    override fun validateCode(code: String) {
        if (isValidatingCodeAtomic.get()) {
            return
        }
        viewModelScope.launch(Dispatchers.IO) {
            if (!isValidatingCodeAtomic.compareAndSet(false, true)) {
                return@launch
            }
            try {
                val attempt = db.accountsDao().getChangePasswordAttempts().firstOrNull() ?: return@launch
                val response =
                    api.accountsValidateChangePassword(ValidateChangePasswordRequestDto(attempt.randomToken, code))
                when (response) {
                    is NetworkResponse.Success -> {
                        db.runInTransaction {
                            db.accountsDao().deleteChangePasswordCommits()
                            db.accountsDao()
                                .putChangePasswordCommit(ChangePasswordCommit(response.body.token.changePasswordToken))
                        }
                        _navEvents.emit(ChangePasswordNavEvent.ContinueChangePassword)
                    }
                    is NetworkResponse.ServerError -> {
                        val errorsBody = response.body
                        if (errorsBody != null) {
                            if (errorsBody.errors.codeInvalid != null) {
                                _errors.emit(ChangePasswordError.InvalidCode)
                            }
                            if (errorsBody.errors.passwordChangeExpired != null) {
                                _errors.emit(ChangePasswordError.CodeExpired)
                            }
                        }
                    }
                    is NetworkResponse.NetworkError -> _errors.emit(ChangePasswordError.BadNetwork)
                    is NetworkResponse.UnknownError -> Unit
                }
            } finally {
                isValidatingCodeAtomic.compareAndSet(true, false)
            }
        }
    }

    override fun sendAnotherCode() {
        if (isSendingCodeAtomic.get()) {
            return
        }
        viewModelScope.launch(Dispatchers.IO) {
            if (!isSendingCodeAtomic.compareAndSet(false, true)) {
                return@launch
            }
            try {
                val attempt = db.accountsDao().getChangePasswordAttempts().firstOrNull() ?: return@launch
                when (val response =
                    api.accountsRevalidateChangePasswordCredentials(RevalidateCredentialsDto(attempt.randomToken))) {
                    is NetworkResponse.Success -> {
                        val newAttempt = attempt.copy(
                            createdAt = System.currentTimeMillis(),
                            expiresInSeconds = response.body.tokenAndTimeDto.expiresIn.toLong(),
                            randomToken = response.body.tokenAndTimeDto.randomToken,
                        )
                        db.runInTransaction {
                            db.accountsDao().deleteChangePasswordAttempts()
                            db.accountsDao().putChangePasswordAttempt(newAttempt)
                        }
                        _timerFlow.emit(newAttempt.expiresInSeconds * 1000)
                    }
                    is NetworkResponse.ServerError -> Unit
                    is NetworkResponse.NetworkError -> _errors.emit(ChangePasswordError.BadNetwork)
                    is NetworkResponse.UnknownError -> Unit
                }
            } finally {
                isSendingCodeAtomic.compareAndSet(true, false)
            }
        }
    }

    override fun loadTimerAsync(): Deferred<Long?> =
        viewModelScope.async(Dispatchers.IO) {
            val attempt = db.accountsDao().getChangePasswordAttempts().firstOrNull() ?: return@async null
            attempt.expiresInSeconds * 1000 - (System.currentTimeMillis() - attempt.createdAt)
        }
}