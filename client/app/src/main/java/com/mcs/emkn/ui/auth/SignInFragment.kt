package com.mcs.emkn.ui.auth

import android.app.AlertDialog
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.activity.addCallback
import androidx.core.view.isVisible
import androidx.core.widget.doAfterTextChanged
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import com.mcs.emkn.R
import com.mcs.emkn.core.Router
import com.mcs.emkn.databinding.FragmentSignInBinding
import com.mcs.emkn.ui.auth.viewmodels.SignInError
import com.mcs.emkn.ui.auth.viewmodels.SignInInteractor
import com.mcs.emkn.ui.auth.viewmodels.SignInNavEvent
import com.mcs.emkn.ui.auth.viewmodels.SignInViewModel
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class SignInFragment : Fragment() {
    private var _binding: FragmentSignInBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var router: Router

    private val signInInteractor: SignInInteractor by viewModels<SignInViewModel>()

    private var isLoadingStarted = false

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        _binding = FragmentSignInBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.forgotPasswordButton.setOnClickListener {
            router.goToForgotPasswordScreen()
        }
        binding.submitButton.setOnClickListener {
            clearErrorFields()
            val login = binding.loginEditText.text?.toString() ?: return@setOnClickListener
            val password = binding.passwordEditText.text?.toString() ?: return@setOnClickListener
            signInInteractor.onSignInClick(login, password)
            binding.submitButton.isEnabled = false
        }
        subscribeToFormFields()
        requireActivity().onBackPressedDispatcher.addCallback(this) {
            onBackButtonPressed()
            this.isEnabled = true
        }
        subscribeToLoadingStatus()
        subscribeToErrorsStatus()
        subscribeToNavStatus()

        insertLoginPassword()
    }

    private fun insertLoginPassword() {
        lifecycleScope.launch {
            signInInteractor.loadCredentialsAsync().await()?.let { credentials ->
                binding.loginEditText.setText(credentials.login)
                binding.passwordEditText.setText(credentials.password)
            }
        }
    }

    private fun decideSignInButtonEnabledState(login: String?, password: String?) {
        binding.submitButton.isEnabled =
            !(login.isNullOrBlank() || password.isNullOrBlank() || isLoadingStarted)
    }

    private fun subscribeToFormFields() {
        decideSignInButtonEnabledState(
            login = binding.loginEditText.text?.toString(),
            password = binding.passwordEditText.text?.toString()
        )
        binding.loginEditText.doAfterTextChanged { login ->
            decideSignInButtonEnabledState(
                login = login?.toString(),
                password = binding.passwordEditText.text?.toString()
            )
        }
        binding.passwordEditText.doAfterTextChanged { password ->
            decideSignInButtonEnabledState(
                login = binding.loginEditText.text?.toString(),
                password = password?.toString()
            )
        }
    }

    private fun onBackButtonPressed() {
        val login = binding.loginEditText.text?.toString()
        val password = binding.passwordEditText.text?.toString()
        if (login.isNullOrBlank() && password.isNullOrBlank()) {
            router.back()
            return
        }
        AlertDialog.Builder(requireContext())
            .setTitle(R.string.sign_in_back_alert_dialog_text)
            .setNegativeButton(R.string.sign_in_back_alert_dialog_cancel_button_text) { dialog, _ ->
                dialog?.dismiss()
            }
            .setPositiveButton(R.string.sign_in_back_alert_dialog_ok_button_text) { _, _ ->
                router.back()
            }
            .show()
    }

    private fun subscribeToLoadingStatus() {
        lifecycleScope.launch {
            lifecycle.repeatOnLifecycle(Lifecycle.State.STARTED) {
                signInInteractor.isLoadingFlow.collect { isLoading ->
                    binding.progressBar.isVisible = isLoading
                    isLoadingStarted = isLoading
                    decideSignInButtonEnabledState(
                        login = binding.loginEditText.text?.toString(),
                        password = binding.passwordEditText.text?.toString()
                    )
                }
            }
        }
    }

    private fun subscribeToErrorsStatus() {
        lifecycleScope.launch {
            lifecycle.repeatOnLifecycle(Lifecycle.State.STARTED) {
                signInInteractor.errorsFlow.collect { error ->
                    when (error) {
                        is SignInError.BadNetwork -> {
                            Toast
                                .makeText(
                                    requireContext(),
                                    resources.getString(R.string.bad_network_error),
                                    Toast.LENGTH_LONG
                                )
                                .show()
                        }
                        is SignInError.IncorrectLoginOrPassword -> {
                            binding.underPasswordTextView.text =
                                resources.getString(R.string.incorrect_login_or_password_error)
                        }
                    }
                }
            }
        }
    }

    private fun subscribeToNavStatus() {
        lifecycleScope.launch {
            lifecycle.repeatOnLifecycle(Lifecycle.State.STARTED) {
                signInInteractor.navEvents.collect { navEvent ->
                    when (navEvent) {
                        is SignInNavEvent.ContinueSignIn -> {
                            router.goToUserNavGraph()
                        }
                    }
                }
            }
        }
    }

    private fun clearErrorFields() {
        binding.underPasswordTextView.text = ""
    }
}