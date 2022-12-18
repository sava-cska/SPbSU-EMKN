package com.mcs.emkn.ui.auth.viewmodels

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.haroldadmin.cnradapter.NetworkResponse
import com.mcs.emkn.database.Database
import com.mcs.emkn.database.entities.SignUpAttempt
import com.mcs.emkn.network.Api
import com.mcs.emkn.network.dto.request.RegistrationRequestDto
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class SignUpViewModel @Inject constructor(
    private val api: Api,
    private val database: Database,
    private val observer: SignInSignUpQueriesObserver,
) : SignUpInteractor, ViewModel() {
    override val isLoadingFlow: Flow<Boolean>
        get() = _isLoadingFlow
    override val errorsFlow: Flow<SignUpError>
        get() = _errorsFlow
    override val navEvents: Flow<SignUpNavEvent>
        get() = _navEvents

    private val _isLoadingFlow = MutableStateFlow(false)
    private val _errorsFlow = MutableSharedFlow<SignUpError>()
    private val _navEvents = MutableSharedFlow<SignUpNavEvent>()
    private val isLoadingAtomic = observer.isLoading

    override fun onSignUpClick(email: String, login: String, password: String, name: String, surname: String) {
        if (isLoadingAtomic.get()) {
            return
        }
        viewModelScope.launch(Dispatchers.IO) {
            if (!isLoadingAtomic.compareAndSet(false, true)) {
                return@launch
            }
            try {
                _isLoadingFlow.emit(true)
                val response = api.accountsRegister(
                    RegistrationRequestDto(
                        login = login,
                        password = password,
                        email = email,
                        firstName = name,
                        lastName = surname
                    )
                )

                when (response) {
                    is NetworkResponse.Success -> {
                        database.runInTransaction {
                            database.accountsDao().deleteSignUpAttempts()
                            database.accountsDao().putSignUpAttempt(
                                SignUpAttempt(
                                    email = email,
                                    login = login,
                                    password = password,
                                    name = name,
                                    surName = surname,
                                    randomToken = response.body.tokenAndTimeDto.randomToken,
                                    expiresInSeconds = response.body.tokenAndTimeDto.expiresIn.toLong(),
                                    createdAt = System.currentTimeMillis(),
                                )
                            )
                        }
                        _navEvents.emit(SignUpNavEvent.ContinueSignUp)
                    }
                    is NetworkResponse.ServerError -> {
                        val errorsBody = response.body
                        if (errorsBody != null) {
                            if (errorsBody.errors.illegalLogin != null) {
                                _errorsFlow.emit(SignUpError.IllegalLogin)
                            }
                            if (errorsBody.errors.illegalPassword != null) {
                                _errorsFlow.emit(SignUpError.IllegalPassword)
                            }
                            if (errorsBody.errors.loginIsNotAvailable != null) {
                                _errorsFlow.emit(SignUpError.LoginIsNotAvailable)
                            }
                            if (errorsBody.errors.emailIsNotAvailable != null) {
                                _errorsFlow.emit(SignUpError.EmailIsNotAvailable)
                            }
                        }
                    }
                    is NetworkResponse.NetworkError -> _errorsFlow.emit(SignUpError.BadNetwork)
                    is NetworkResponse.UnknownError -> Unit
                }
            } finally {
                isLoadingAtomic.compareAndSet(true, false)
                _isLoadingFlow.emit(false)
            }
        }
    }
}