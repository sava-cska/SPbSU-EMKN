package com.mcs.emkn.ui.auth

import android.app.AlertDialog
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.activity.addCallback
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import com.mcs.emkn.R
import com.mcs.emkn.core.Router
import com.mcs.emkn.databinding.FragmentSignUpBinding
import com.mcs.emkn.ui.auth.viewmodels.SignUpError
import com.mcs.emkn.ui.auth.viewmodels.SignUpInteractor
import com.mcs.emkn.ui.auth.viewmodels.SignUpNavEvent
import com.mcs.emkn.ui.auth.viewmodels.SignUpViewModel
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class SignUpFragment : Fragment() {
    private var _binding: FragmentSignUpBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var router: Router

    private val signUpInteractor: SignUpInteractor by viewModels<SignUpViewModel>()

    private var isLoadingStarted = false


    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        _binding = FragmentSignUpBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.signUpButton.setOnClickListener {
            clearErrorFields()
            binding.signUpButton.isEnabled = false

            signUpInteractor.onSignUpClick(
                email = binding.emailEditText.text?.toString() ?: return@setOnClickListener,
                login = binding.loginEditText.text?.toString() ?: return@setOnClickListener,
                password = binding.passwordEditText.text?.toString() ?: return@setOnClickListener,
                name = binding.firstnameEditText.text?.toString() ?: return@setOnClickListener,
                surname = binding.lastnameEditText.text?.toString() ?: return@setOnClickListener,
            )
        }
        subscribeToFormFields()
        requireActivity().onBackPressedDispatcher.addCallback(this) {
            onBackButtonPressed()
            this.isEnabled = true
        }

        subscribeToLoadingStatus()
        subscribeToErrorsStatus()
        subscribeToNavStatus()
    }

    private fun decideSignUpButtonEnabledState(
        firstName: String?,
        lastName: String?,
        login: String?,
        email: String?,
        password: String?
    ) {
        binding.signUpButton.isEnabled =
            !(firstName.isNullOrBlank() || lastName.isNullOrBlank() || login.isNullOrBlank() ||
                email.isNullOrBlank() || password.isNullOrBlank() || isLoadingStarted)
    }

    private fun subscribeToFormFields() {
        decideSignUpButtonEnabledState(
            firstName = binding.firstnameEditText.text?.toString(),
            lastName = binding.lastnameEditText.text?.toString(),
            login = binding.loginEditText.text?.toString(),
            email = binding.emailEditText.text?.toString(),
            password = binding.passwordEditText.text?.toString()
        )
        val watcher = object : TextWatcher {
            override fun beforeTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {}

            override fun onTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {}

            override fun afterTextChanged(p0: Editable?) {
                decideSignUpButtonEnabledState(
                    firstName = binding.firstnameEditText.text?.toString(),
                    lastName = binding.lastnameEditText.text?.toString(),
                    login = binding.loginEditText.text?.toString(),
                    email = binding.emailEditText.text?.toString(),
                    password = binding.passwordEditText.text?.toString(),
                )
            }
        }

        binding.firstnameEditText.addTextChangedListener(watcher)
        binding.lastnameEditText.addTextChangedListener(watcher)
        binding.loginEditText.addTextChangedListener(watcher)
        binding.emailEditText.addTextChangedListener(watcher)
        binding.passwordEditText.addTextChangedListener(watcher)
    }

    private fun onBackButtonPressed() {
        val firstname = binding.firstnameEditText.text?.toString()
        val lastname = binding.lastnameEditText.text?.toString()
        val nickname = binding.loginEditText.text?.toString()
        val email = binding.emailEditText.text?.toString()
        val password = binding.passwordEditText.text?.toString()
        if (firstname.isNullOrBlank()
            && lastname.isNullOrBlank()
            && nickname.isNullOrBlank()
            && email.isNullOrBlank()
            && password.isNullOrBlank()
        ) {
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
                signUpInteractor.isLoadingFlow.collect { isLoading ->
                    binding.progressBar.isVisible = isLoading
                    isLoadingStarted = isLoading
                    decideSignUpButtonEnabledState(
                        firstName = binding.firstnameEditText.text?.toString(),
                        lastName = binding.lastnameEditText.text?.toString(),
                        login = binding.loginEditText.text?.toString(),
                        email = binding.emailEditText.text?.toString(),
                        password = binding.passwordEditText.text?.toString()
                    )
                }
            }
        }
    }

    private fun subscribeToErrorsStatus() {
        lifecycleScope.launch {
            lifecycle.repeatOnLifecycle(Lifecycle.State.STARTED) {
                signUpInteractor.errorsFlow.collect { error ->
                    when (error) {
                        is SignUpError.BadNetwork -> {
                            Toast
                                .makeText(
                                    requireContext(),
                                    resources.getString(R.string.bad_network_error),
                                    Toast.LENGTH_LONG
                                )
                                .show()
                        }
                        is SignUpError.IllegalLogin -> {
                            binding.underLoginTextView.text =
                                resources.getString(R.string.incorrect_login_error)
                        }
                        is SignUpError.IllegalPassword -> {
                            binding.underPasswordTextView.text =
                                resources.getString(R.string.incorrect_password_error)
                        }
                        is SignUpError.LoginIsNotAvailable -> {
                            binding.underLoginTextView.text =
                                resources.getString(R.string.login_not_available_error)
                        }
                        is SignUpError.EmailIsNotAvailable -> {
                            binding.underEmailTextView.text =
                                resources.getString(R.string.email_not_available_error)
                        }
                    }
                }
            }
        }
    }

    private fun subscribeToNavStatus() {
        lifecycleScope.launch {
            lifecycle.repeatOnLifecycle(Lifecycle.State.STARTED) {
                signUpInteractor.navEvents.collect { navEvent ->
                    when (navEvent) {
                        is SignUpNavEvent.ContinueSignUp -> {
                            router.goToEmailConfirmationScreen()
                        }
                    }
                }
            }
        }
    }

    private fun clearErrorFields() {
        binding.underPasswordTextView.text = ""
        binding.underLoginTextView.text = ""
        binding.underEmailTextView.text = ""
    }
}