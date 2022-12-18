package com.mcs.emkn.ui.changepassword.viewmodels

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.haroldadmin.cnradapter.NetworkResponse
import com.mcs.emkn.database.Database
import com.mcs.emkn.database.entities.ChangePasswordAttempt
import com.mcs.emkn.network.Api
import com.mcs.emkn.network.dto.request.BeginChangePasswordRequestDto
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.launch
import java.util.concurrent.atomic.AtomicBoolean
import javax.inject.Inject

@HiltViewModel
class ForgotPasswordViewModel @Inject constructor(
    private val api: Api,
    private val database: Database,
) : ForgotPasswordInteractor, ViewModel() {
    override val errors: Flow<ForgotPasswordError>
        get() = _errors
    override val navEvents: Flow<ForgotPasswordNavEvent>
        get() = _navEvents

    private val _errors = MutableSharedFlow<ForgotPasswordError>()
    private val _navEvents = MutableSharedFlow<ForgotPasswordNavEvent>()

    private val submitAtomic = AtomicBoolean(false)

    override fun onSubmitClick(credentials: String) {
        if (submitAtomic.get()) {
            return
        }
        viewModelScope.launch(Dispatchers.IO) {
            if (!submitAtomic.compareAndSet(false, true)) {
                return@launch
            }
            try {
                val response = api.accountsBeginChangePassword(BeginChangePasswordRequestDto(credentials))
                when (response) {
                    is NetworkResponse.Success -> {
                        val attempt = ChangePasswordAttempt(
                            credentials = credentials,
                            randomToken = response.body.tokenAndTimeDto.randomToken,
                            expiresInSeconds = response.body.tokenAndTimeDto.expiresIn.toLong(),
                            createdAt = System.currentTimeMillis()
                        )
                        database.runInTransaction {
                            database.accountsDao().deleteChangePasswordAttempts()
                            database.accountsDao().putChangePasswordAttempt(attempt)
                        }
                        _navEvents.emit(ForgotPasswordNavEvent.ContinueForgotPassword)
                    }
                    is NetworkResponse.ServerError -> {
                        val errorsBody = response.body
                        if (errorsBody != null) {
                            if (errorsBody.errors.illegalEmail != null) {
                                _errors.emit(ForgotPasswordError.InvalidEmail)
                            }
                        }
                    }
                    is NetworkResponse.NetworkError -> _errors.emit(ForgotPasswordError.BadNetwork)
                    is NetworkResponse.UnknownError -> Unit
                }
            } finally {
                submitAtomic.compareAndSet(true, false)
            }
        }
    }
}